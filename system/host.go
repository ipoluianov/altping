package system

import (
	"errors"
	"fmt"
	"net"
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

	defaultData map[int][]byte
	counter     int

	statOK             int
	statERR            int
	IP                 string
	resultErr          error
	resultLastLiveIP   string
	resultLastPingTime time.Duration
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

func (c *Host) checkIP() bool {
	if len(c.IP) == 0 {
		ips, err := net.LookupIP(c.configHost.Hostname)
		if err != nil {
			c.mtx.Lock()
			c.IP = ""
			c.mtx.Unlock()
			c.resultErr = errors.New("cannot resolve hostname")
			return false
		}

		// Find IP v4 address
		var ip4 net.IP
		for _, ip := range ips {
			if ip.To4() != nil {
				ip4 = ip
				break
			}
		}
		if ip4 == nil {
			c.mtx.Lock()
			c.IP = ""
			c.mtx.Unlock()
			c.resultErr = fmt.Errorf("no IPv4 address found for %s", c.configHost.Hostname)
			return false
		}
		c.mtx.Lock()
		c.IP = ip4.String()
		c.mtx.Unlock()
	}
	return true
}

func (c *Host) updateState() {
	c.mtx.Lock()
	var state HostState
	state.ConfigHost = c.configHost
	state.Started = c.started
	state.Stopping = c.stopping
	state.LastError = c.resultErr
	state.LastCheck = time.Now()
	state.StatOK = c.statOK
	state.StatERR = c.statERR
	state.StatIP = c.resultLastLiveIP
	state.PingTime = c.resultLastPingTime
	c.mtx.Unlock()
	c.SetState(state)
}

func (c *Host) thWork() {
	c.mtx.Lock()
	c.started = true
	c.mtx.Unlock()
	for {
		time.Sleep(1000 * time.Millisecond)

		c.mtx.Lock()
		if !c.started || c.stopping {
			c.mtx.Unlock()
			break
		}
		c.mtx.Unlock()

		if c.checkIP() {
			result, peer, err := c.ping(c.IP, 64, 1000)
			c.resultLastPingTime = time.Duration(result) * time.Millisecond

			liveIP := ""

			if err != nil {
				c.statERR++
				c.resultErr = err
			} else {
				c.statOK++
				c.resultErr = nil
				ipWithoutPort, _, _ := net.SplitHostPort(peer.String())
				liveIP = ipWithoutPort
			}
			c.resultLastLiveIP = liveIP
		} else {
			c.mtx.Lock()
			c.statERR++
			c.mtx.Unlock()
		}

		c.updateState()
	}
	c.mtx.Lock()
	c.started = false
	c.mtx.Unlock()
}

func (c *Host) ping(addr string, dataSize int, timeoutMs int) (result int, peer net.Addr, err error) {
	return GetServer().PingHost(addr, dataSize, timeoutMs)
}
