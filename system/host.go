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

	// chanStop is closed to stop the host goroutine, chanDone is closed when it has finished
	chanStop chan struct{}
	chanDone chan struct{}

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

	return &c
}

func (c *Host) Start() {
	c.mtx.Lock()
	if c.started {
		c.mtx.Unlock()
		return
	}
	c.started = true
	c.stopping = false
	c.chanStop = make(chan struct{})
	c.chanDone = make(chan struct{})
	c.mtx.Unlock()
	c.resetStat()
	go c.thWork()

	fmt.Println("Host Start", c.configHost.ID)
}

// Stop stops the host and waits for its goroutine to finish
func (c *Host) Stop() {
	c.signalStop()
	c.waitStopped()
}

// signalStop asks the host to stop without waiting,
// so that many hosts can be stopped at once
func (c *Host) signalStop() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if !c.started || c.stopping {
		return
	}
	c.stopping = true
	close(c.chanStop)
}

func (c *Host) waitStopped() {
	c.mtx.Lock()
	chanDone := c.chanDone
	c.mtx.Unlock()
	if chanDone != nil {
		<-chanDone
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
	c.addGapToHistory()
	defer func() {
		c.addGapToHistory()
		c.mtx.Lock()
		c.started = false
		c.mtx.Unlock()
		close(c.chanDone)
	}()

	// The first ping goes right away, so a new host shows its state at once
	timeout := time.Duration(0)

	for {
		select {
		case <-time.After(timeout):
		case <-c.chanStop:
			return
		}
		timeout = time.Duration(1000) * time.Millisecond

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
}
