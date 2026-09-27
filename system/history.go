package system

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ipoluianov/altping/config"
)

// historyDepth is how long ping results are kept in memory and on disk
const historyDepth = 24 * time.Hour

// On-disk record: DT (unix nanoseconds), PingTime (nanoseconds), flags.
// Little endian, fixed size, so a record torn by a crash is just dropped.
const (
	historyFileExt    = ".hist"
	historyRecordSize = 8 + 8 + 1

	historyFlagOK  = 1
	historyFlagGap = 2
)

type HistorySample struct {
	DT       time.Time
	PingTime time.Duration
	OK       bool
	// Gap marks where pinging was started or stopped: a break in the data, not a failure
	Gap bool
}

// HostHistory is the ping history of one host.
// It is written by the host goroutine and read by the UI.
// Every sample is also appended to a file, so the history survives restarts.
type HostHistory struct {
	mtx      sync.RWMutex
	samples  []HistorySample
	filePath string
	file     *os.File
}

// historyDirectory returns the directory with the history files, next to the configs
func historyDirectory() string {
	return filepath.Join(config.ConfigDirectory(), "history")
}

func historyFilePath(hostID string) string {
	return filepath.Join(historyDirectory(), filepath.Base(hostID)+historyFileExt)
}

// NewHostHistory loads the saved history of the host and opens its file for appending
func NewHostHistory(hostID string) *HostHistory {
	c := &HostHistory{filePath: historyFilePath(hostID)}
	c.load()
	return c
}

func (c *HostHistory) load() {
	bs, err := os.ReadFile(c.filePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("History read error:", err)
	}

	limit := time.Now().Add(-historyDepth)
	for len(bs) >= historyRecordSize {
		s := decodeHistorySample(bs[:historyRecordSize])
		bs = bs[historyRecordSize:]
		if s.DT.Before(limit) {
			continue
		}
		// Keep the samples sorted even if the system clock went back
		if n := len(c.samples); n > 0 && s.DT.Before(c.samples[n-1].DT) {
			continue
		}
		c.samples = append(c.samples, s)
	}

	// The program was not stopped properly: close the data before the downtime
	if n := len(c.samples); n > 0 && !c.samples[n-1].Gap {
		c.samples = append(c.samples, HistorySample{DT: c.samples[n-1].DT, Gap: true})
	}

	c.rewriteFile()
}

// rewriteFile replaces the file with the in-memory samples and reopens it for appending
func (c *HostHistory) rewriteFile() {
	c.closeFile()

	err := os.MkdirAll(filepath.Dir(c.filePath), 0755)
	if err != nil {
		fmt.Println("History write error:", err)
		return
	}

	var buf bytes.Buffer
	buf.Grow(len(c.samples) * historyRecordSize)
	for _, s := range c.samples {
		buf.Write(encodeHistorySample(s))
	}
	tmpPath := c.filePath + ".tmp"
	err = os.WriteFile(tmpPath, buf.Bytes(), 0644)
	if err == nil {
		err = os.Rename(tmpPath, c.filePath)
	}
	if err != nil {
		fmt.Println("History write error:", err)
		return
	}

	c.file, err = os.OpenFile(c.filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("History write error:", err)
		c.file = nil
	}
}

func (c *HostHistory) closeFile() {
	if c.file != nil {
		c.file.Close()
		c.file = nil
	}
}

// Close closes the history file; the in-memory samples stay readable
func (c *HostHistory) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.closeFile()
}

func (c *HostHistory) Add(s HistorySample) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.samples = append(c.samples, s)

	// Drop samples older than historyDepth.
	// Compact only when a noticeable part has expired to avoid copying on every sample.
	limit := s.DT.Add(-historyDepth)
	expired := sort.Search(len(c.samples), func(i int) bool { return !c.samples[i].DT.Before(limit) })
	if expired > 0 && expired >= len(c.samples)/8 {
		c.samples = append(c.samples[:0:0], c.samples[expired:]...)
		c.rewriteFile()
		return
	}

	if c.file != nil {
		_, err := c.file.Write(encodeHistorySample(s))
		if err != nil {
			fmt.Println("History write error:", err)
		}
	}
}

// Visit calls fn for the samples within [from, to] plus the nearest sample on each side,
// in time order, without copying them. fn must not call the history.
func (c *HostHistory) Visit(from, to time.Time, fn func(s HistorySample)) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	i0 := sort.Search(len(c.samples), func(i int) bool { return !c.samples[i].DT.Before(from) })
	i1 := sort.Search(len(c.samples), func(i int) bool { return c.samples[i].DT.After(to) })
	if i0 > 0 {
		i0--
	}
	if i1 < len(c.samples) {
		i1++
	}
	for _, s := range c.samples[i0:i1] {
		fn(s)
	}
}

// HistoryStats summarizes the pings of a period
type HistoryStats struct {
	Sent int
	Lost int
	Avg  time.Duration // of the successful pings
	Max  time.Duration
}

// LossPercent returns the share of lost pings, 0 when nothing was sent
func (s HistoryStats) LossPercent() float64 {
	if s.Sent == 0 {
		return 0
	}
	return float64(s.Lost) * 100 / float64(s.Sent)
}

// Stats summarizes the pings since the time
func (c *HostHistory) Stats(since time.Time) HistoryStats {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	var res HistoryStats
	var sum time.Duration
	i0 := sort.Search(len(c.samples), func(i int) bool { return !c.samples[i].DT.Before(since) })
	for _, s := range c.samples[i0:] {
		if s.Gap {
			continue
		}
		res.Sent++
		if !s.OK {
			res.Lost++
			continue
		}
		sum += s.PingTime
		res.Max = max(res.Max, s.PingTime)
	}
	if ok := res.Sent - res.Lost; ok > 0 {
		res.Avg = sum / time.Duration(ok)
	}
	return res
}

func encodeHistorySample(s HistorySample) []byte {
	bs := make([]byte, historyRecordSize)
	binary.LittleEndian.PutUint64(bs[0:], uint64(s.DT.UnixNano()))
	binary.LittleEndian.PutUint64(bs[8:], uint64(s.PingTime))
	var flags byte
	if s.OK {
		flags |= historyFlagOK
	}
	if s.Gap {
		flags |= historyFlagGap
	}
	bs[16] = flags
	return bs
}

func decodeHistorySample(bs []byte) HistorySample {
	flags := bs[16]
	return HistorySample{
		DT:       time.Unix(0, int64(binary.LittleEndian.Uint64(bs[0:]))),
		PingTime: time.Duration(binary.LittleEndian.Uint64(bs[8:])),
		OK:       flags&historyFlagOK != 0,
		Gap:      flags&historyFlagGap != 0,
	}
}

// removeStaleHistoryFiles deletes the files of hosts that have not been pinged
// for longer than historyDepth (removed hosts, other configurations)
func removeStaleHistoryFiles() {
	entries, err := os.ReadDir(historyDirectory())
	if err != nil {
		return
	}
	limit := time.Now().Add(-historyDepth)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), historyFileExt) {
			continue
		}
		info, err := e.Info()
		if err != nil || !info.ModTime().Before(limit) {
			continue
		}
		os.Remove(filepath.Join(historyDirectory(), e.Name()))
	}
}
