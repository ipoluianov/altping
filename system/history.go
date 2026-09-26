package system

import (
	"sort"
	"sync"
	"time"
)

// historyDepth is how long ping results are kept in memory
const historyDepth = 24 * time.Hour

type HistorySample struct {
	DT       time.Time
	PingTime time.Duration
	OK       bool
	// Gap marks where pinging was started or stopped: a break in the data, not a failure
	Gap bool
}

// HostHistory is the in-memory ping history of one host.
// It is written by the host goroutine and read by the UI.
type HostHistory struct {
	mtx     sync.RWMutex
	samples []HistorySample
}

func NewHostHistory() *HostHistory {
	return &HostHistory{}
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
	}
}

// Range returns the samples within [from, to] plus the nearest sample on each side
func (c *HostHistory) Range(from, to time.Time) []HistorySample {
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
	res := make([]HistorySample, i1-i0)
	copy(res, c.samples[i0:i1])
	return res
}
