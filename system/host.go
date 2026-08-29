package system

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ipoluianov/altping/config"
)

type Host struct {
	mtx      sync.Mutex
	ID       string
	started  bool
	stopping bool

	config     *config.Config
	configHost config.ConfigHost

	lastState HostState

	defaultData    map[int][]byte
	sequenceNumber uint16
	counter        int

	statOK  int
	statERR int
	statIP  string
}

type HostState struct {
	ConfigHost config.ConfigHost
	Started    bool
	Stopping   bool
	LastError  error
	LastCheck  time.Time
	PingTime   time.Duration
	StatOK     int
	StatERR    int
	StatIP     string
}

func NewHost(id string) *Host {
	var c Host
	c.ID = id

	c.defaultData = make(map[int][]byte)
	for s := 0; s < 1500; s++ {
		d := make([]byte, s)
		for i := 0; i < s; i++ {
			d[i] = byte(i%26) + 0x41
		}
		c.defaultData[s] = d
	}

	return &c
}

func (c *Host) Start() {
	c.mtx.Lock()
	if c.started {
		c.mtx.Unlock()
		return
	}
	c.stopping = false
	c.mtx.Unlock()
	c.resetStat()
	go c.thWork()
}

func (c *Host) Stop() {
	c.mtx.Lock()
	c.stopping = true
	c.mtx.Unlock()
	for c.started {
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Host) GetState() HostState {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.lastState
}

func (c *Host) SetState(state HostState) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.lastState = state

	//fmt.Println("State", state.ConfigHost.ID, "Started", state.Started, "Stopping", state.Stopping, "LastError", state.LastError, "LastCheck", state.LastCheck, "PingTime", state.PingTime)
}

func (c *Host) resetStat() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.statOK = 0
	c.statERR = 0
}

func (c *Host) UpdateConfig() {
	config := config.Get()
	c.config = config
	c.configHost = config.GetHost(c.ID)
	c.resetStat()
}

func (c *Host) thWork() {
	c.mtx.Lock()
	c.started = true
	c.mtx.Unlock()
	for {
		c.mtx.Lock()
		if !c.started || c.stopping {
			c.mtx.Unlock()
			break
		}
		c.mtx.Unlock()
		time.Sleep(1000 * time.Millisecond)
		result, peer, err := c.ping(c.configHost.Hostname, 64, 1000, false)
		_ = peer
		if err != nil {
			c.statERR++
			c.statIP = ""
		} else {
			c.statOK++
			c.statIP = "123"
		}
		c.mtx.Lock()
		var state HostState
		state.ConfigHost = c.configHost
		state.Started = c.started
		state.Stopping = c.stopping
		state.LastError = err
		state.LastCheck = time.Now()
		state.StatOK = c.statOK
		state.StatERR = c.statERR
		state.StatIP = c.statIP
		fmt.Println("WORK", c.configHost.ID, "StatOK", c.statOK, "StatERR", c.statERR, "StatIP", c.statIP)
		if err == nil {
			state.PingTime = time.Duration(result) * time.Millisecond
		}
		c.mtx.Unlock()
		c.SetState(state)
	}
	c.mtx.Lock()
	c.started = false
	c.mtx.Unlock()
}

func (c *Host) ping(addr string, dataSize int, timeoutMs int, useUdpSocket bool) (result int, peer net.Addr, err error) {
	var data []byte
	data, _ = c.defaultData[dataSize]
	if data == nil {
		data = make([]byte, dataSize)
	}

	c.mtx.Lock()
	seqIndex := c.sequenceNumber
	srcIndex := uint16(os.Getpid() & 0xFFFF)
	c.sequenceNumber++
	if c.sequenceNumber > 65534 {
		c.sequenceNumber = 1
	}
	c.mtx.Unlock()

	return GetServer().PingHost(addr, data, timeoutMs, srcIndex, seqIndex, useUdpSocket)
}
