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

	// The state of the last pings, kept up to date on every sample
	lastChange time.Time // when the ping result last changed (OK <-> failed) or pinging started
	failStreak int       // failed pings in a row at the end
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
		c.track(s)
	}

	// The program was not stopped properly: close the data before the downtime
	if n := len(c.samples); n > 0 && !c.samples[n-1].Gap {
		c.samples = append(c.samples, HistorySample{DT: c.samples[n-1].DT, Gap: true})
		c.track(c.samples[n])
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
	c.track(s)

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

// track updates lastChange and failStreak with the next sample
func (c *HostHistory) track(s HistorySample) {
	n := len(c.samples)
	var prev *HistorySample
	if n >= 2 {
		prev = &c.samples[n-2]
	}
	switch {
	case s.Gap:
		// Pinging stopped or started: what happened meanwhile is unknown
		c.failStreak = 0
		c.lastChange = time.Time{}
	case prev == nil || prev.Gap || prev.OK != s.OK:
		c.lastChange = s.DT
	}
	if !s.Gap {
		if s.OK {
			c.failStreak = 0
		} else {
			c.failStreak++
		}
	}
}

// LastChange returns when the ping result last changed between OK and failed
// (or pinging started); zero when nothing is known
func (c *HostHistory) LastChange() time.Time {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.lastChange
}

// FailStreak returns how many pings in a row have failed
func (c *HostHistory) FailStreak() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.failStreak
}

// Outage is a period when the host did not reply
type Outage struct {
	Start   time.Time
	End     time.Time // the first reply after it, or when pinging stopped
	Ongoing bool      // the host is still not replying; End is the last failed ping
}

// Outages returns the periods without replies since the time, the oldest first
func (c *HostHistory) Outages(since time.Time) []Outage {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	var res []Outage
	var cur *Outage
	i0 := sort.Search(len(c.samples), func(i int) bool { return !c.samples[i].DT.Before(since) })
	for _, s := range c.samples[i0:] {
		failed := !s.Gap && !s.OK
		switch {
		case failed && cur == nil:
			res = append(res, Outage{Start: s.DT, End: s.DT, Ongoing: true})
			cur = &res[len(res)-1]
		case failed:
			cur.End = s.DT
		case cur != nil:
			cur.End = s.DT
			cur.Ongoing = false
			cur = nil
		}
	}
	return res
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
	Min  time.Duration // of the successful pings
	Avg  time.Duration
	Max  time.Duration
	// Jitter is the average difference between the times of successive successful pings
	Jitter time.Duration
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
	var sum, jitterSum time.Duration
	jitterCount := 0
	var prev *HistorySample // the previous successful ping, nil after a failure or a gap
	i0 := sort.Search(len(c.samples), func(i int) bool { return !c.samples[i].DT.Before(since) })
	for i := range c.samples[i0:] {
		s := &c.samples[i0+i]
		if s.Gap {
			prev = nil
			continue
		}
		res.Sent++
		if !s.OK {
			res.Lost++
			prev = nil
			continue
		}
		if res.Sent-res.Lost == 1 || s.PingTime < res.Min {
			res.Min = s.PingTime
		}
		sum += s.PingTime
		res.Max = max(res.Max, s.PingTime)
		if prev != nil {
			d := s.PingTime - prev.PingTime
			jitterSum += max(d, -d)
			jitterCount++
		}
		prev = s
	}
	if ok := res.Sent - res.Lost; ok > 0 {
		res.Avg = sum / time.Duration(ok)
	}
	if jitterCount > 0 {
		res.Jitter = jitterSum / time.Duration(jitterCount)
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
