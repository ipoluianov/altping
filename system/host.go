package system

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
	"sync"
	"syscall"
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

	statOK  int
	statERR int
	IP      string
	// When the name was last looked up, and the failed pings in a row since then (see checkIP)
	resolvedAt  time.Time
	failedPings int

	// The u00.io API key to send the results to, "" - not shared. Set from the UI,
	// so it can be turned on and off without restarting the host.
	shareKey  string
	shareName string // the name shown on the page

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
	c.SetShare(c.configHost)
	c.resetStat()
}

// SetShare starts or stops sending the results to u00.io as the host config says
func (c *Host) SetShare(hostConfig config.ConfigHost) {
	key := ""
	if hostConfig.Share {
		key = hostConfig.ShareKey
	}
	// The name as the table shows it
	name := hostConfig.DisplayName
	if name == "" {
		name = hostConfig.Address()
	}
	c.mtx.Lock()
	c.shareKey = key
	c.shareName = name
	c.mtx.Unlock()
}

// share sends the result of the last ping to u00.io if the host is shared
func (c *Host) share() {
	// A ping interrupted by Stop is not a real result
	select {
	case <-c.chanStop:
		return
	default:
	}
	c.mtx.Lock()
	key, name := c.shareKey, c.shareName
	c.mtx.Unlock()
	if key != "" {
		shareSender.push(key, shareItem{value: shareValue(c.resultLastPingTime, c.resultErr), name: name})
	}
}

// target returns the host to resolve and the TCP port to connect to ("" - ping)
func (c *Host) target() (host string, port string) {
	return c.configHost.Target()
}

// untilStop returns a context that is also cancelled when the host is stopped,
// so Stop does not wait for a slow lookup or connection
func (c *Host) untilStop(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-c.chanStop:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

// The name of a host is looked up again after this time, and after this many failed pings in a row
const (
	reresolveInterval      = 5 * time.Minute
	reresolveAfterFailures = 3
)

// checkIP finds the address to ping. The name is looked up again from time to
// time and when the host stops replying, so a changed address is followed.
// The address in use is kept while the DNS still returns it: a name with
// several addresses (round robin) does not jump between them. When the
// lookup fails, the known address is kept: a DNS failure is not a host failure.
func (c *Host) checkIP() bool {
	if c.IP != "" && c.failedPings < reresolveAfterFailures && time.Since(c.resolvedAt) < reresolveInterval {
		return true
	}
	host, _ := c.target()
	ctx, cancel := c.untilStop(context.Background())
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	cancel()
	c.resolvedAt = time.Now()
	c.failedPings = 0
	ip := pickIP(ips)
	if err != nil || ip == nil {
		if c.IP != "" {
			return true
		}
		c.resultErr = errors.New("cannot resolve hostname")
		return false
	}
	if c.IP != "" && slices.ContainsFunc(ips, net.ParseIP(c.IP).Equal) {
		return true
	}
	c.mtx.Lock()
	c.IP = ip.String()
	c.mtx.Unlock()
	return true
}

// connectTCP measures how long it takes to open a connection to the port
func (c *Host) connectTCP(port string, timeout time.Duration) (time.Duration, net.Addr, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), timeout)
	defer cancelTimeout()
	ctx, cancel := c.untilStop(timeoutCtx)
	defer cancel()
	start := time.Now()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(c.IP, port))
	elapsed := time.Since(start)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
			return elapsed, nil, errors.New("timeout")
		}
		// Windows reports its own error code, so the text is checked too
		if errors.Is(err, syscall.ECONNREFUSED) || strings.Contains(err.Error(), "refused") {
			return elapsed, nil, errors.New("port closed")
		}
		return elapsed, nil, err
	}
	peer := conn.RemoteAddr()
	conn.Close()
	return elapsed, peer, nil
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
		timeout = c.configHost.Interval()

		if c.checkIP() {
			var result time.Duration
			var peer net.Addr
			var err error
			if _, port := c.target(); port != "" {
				result, peer, err = c.connectTCP(port, c.configHost.Timeout())
			} else {
				result, peer, err = c.pingServer.PingHost(c.IP, 64, int(c.configHost.Timeout().Milliseconds()), c.chanStop)
			}
			c.resultLastPingTime = result

			liveIP := ""

			if err != nil {
				c.statERR++
				c.failedPings++
				c.resultErr = err
			} else {
				fmt.Println("OK", c.configHost.ID)
				c.statOK++
				c.failedPings = 0
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
				if tcpAddr, ok := peer.(*net.TCPAddr); ok {
					liveIP = tcpAddr.IP.String()
				}
			}
			c.resultLastLiveIP = liveIP
		} else {
			c.mtx.Lock()
			c.statERR++
			c.mtx.Unlock()
		}

		c.addToHistory()
		c.share()
		c.updateState()
	}
}
