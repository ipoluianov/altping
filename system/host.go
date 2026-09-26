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

	pingServer *PingServer
	history    *HostHistory

	config     *config.Config
	configHost config.ConfigHost

	lastState HostState

	defaultData map[int][]byte
	counter     int

	chanStop chan struct{}

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

func NewHost(id string, pingServer *PingServer, history *HostHistory) *Host {
	var c Host
	c.ID = id
	c.pingServer = pingServer
	c.history = history

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

func (c *Host) Start(chanStopped chan struct{}) {
	c.mtx.Lock()
	if c.started {
		c.mtx.Unlock()
		return
	}
	c.stopping = false
	c.mtx.Unlock()
	c.resetStat()
	c.chanStop = chanStopped
	go c.thWork()

	fmt.Println("Host Start", c.configHost.ID)
}

func (c *Host) Stop() {
	c.mtx.Lock()
	c.stopping = true
	c.mtx.Unlock()
	for c.started {
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *Host) IsRunning() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.started
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

func (c *Host) addToHistory() {
	// A ping interrupted by Stop is not a real result
	select {
	case <-c.chanStop:
		return
	default:
	}

	c.history.Add(HistorySample{
		DT:       time.Now(),
		PingTime: c.resultLastPingTime,
		OK:       c.resultErr == nil,
	})
}

func (c *Host) addGapToHistory() {
	c.history.Add(HistorySample{DT: time.Now(), Gap: true})
}

func (c *Host) thWork() {
	c.mtx.Lock()
	c.started = true
	c.mtx.Unlock()
	c.addGapToHistory()

	timeout := time.Duration(1000) * time.Millisecond

	for {
		// select with timeout
		select {
		case <-time.After(timeout):
		case <-c.chanStop:
			c.addGapToHistory()
			c.mtx.Lock()
			c.started = false
			c.mtx.Unlock()
			return
		}

		c.mtx.Lock()
		if !c.started || c.stopping {
			c.mtx.Unlock()
			break
		}
		c.mtx.Unlock()

		if c.checkIP() {
			result, peer, err := c.pingServer.PingHost(c.IP, 64, 1000, c.chanStop)
			c.resultLastPingTime = result

			liveIP := ""

			if err != nil {
				c.statERR++
				c.resultErr = err
			} else {
				fmt.Println("OK", c.configHost.ID)
				c.statOK++
				c.resultErr = nil

				var ipAddr *net.IPAddr
				ipAddr, ok := peer.(*net.IPAddr)
				if ok {
					liveIP = ipAddr.IP.String()
				}
				udpAddr, ok := peer.(*net.UDPAddr)
				if ok {
					liveIP = udpAddr.IP.String()
				}
			}
			c.resultLastLiveIP = liveIP
		} else {
			c.mtx.Lock()
			c.statERR++
			c.mtx.Unlock()
		}

		c.addToHistory()
		c.updateState()
	}
	c.addGapToHistory()
	c.mtx.Lock()
	c.started = false
	c.mtx.Unlock()
}
