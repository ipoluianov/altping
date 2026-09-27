package system

import (
	"sync"

	"github.com/ipoluianov/altping/config"
)

type System struct {
	pingServer *PingServer
	Hosts      []*Host
	running    bool

	// Ping history by host ID. Hosts are recreated on every restart,
	// so the history is kept here.
	historyMtx sync.Mutex
	history    map[string]*HostHistory

	historyCleanup sync.Once
}

var systemInstance *System

func init() {
	systemInstance = newSystem()
}

func newSystem() *System {
	var c System
	c.pingServer = NewPingServer()
	c.history = make(map[string]*HostHistory)
	return &c
}

func Get() *System {
	return systemInstance
}

func (c *System) Start() {
	if c.running {
		return
	}
	c.running = true
	c.pingServer.Start()
	c.SyncHosts()
}

func (c *System) Stop() {
	if !c.running {
		return
	}
	c.running = false
	for _, host := range c.Hosts {
		host.signalStop()
	}
	c.pingServer.Stop()
	for _, host := range c.Hosts {
		host.waitStopped()
	}
	c.Hosts = nil
}

func (c *System) IsRunning() bool {
	return c.running
}

func (c *System) PingServerMode() string {
	return c.pingServer.Mode()
}

func (c *System) GetHostLastState(id string) HostState {
	for _, host := range c.Hosts {
		if host.ID == id {
			return host.GetState()
		}
	}
	return HostState{}
}

// GetHostHistory returns the ping history of the host (nil for an unknown host)
func (c *System) GetHostHistory(id string) *HostHistory {
	c.historyMtx.Lock()
	defer c.historyMtx.Unlock()
	return c.history[id]
}

// SyncHosts applies the current config to the running hosts: added hosts are started,
// removed ones are stopped, and a host whose address changed is restarted.
// The other hosts keep running with their state and statistics.
func (c *System) SyncHosts() {
	cfg := config.Get()
	addresses := make(map[string]string, len(cfg.Hosts))
	for _, hostConfig := range cfg.Hosts {
		addresses[hostConfig.ID] = hostConfig.Hostname
	}

	kept := make(map[string]*Host)
	var stopping []*Host
	for _, host := range c.Hosts {
		if address, ok := addresses[host.ID]; ok && address == host.configHost.Hostname {
			kept[host.ID] = host
		} else {
			stopping = append(stopping, host)
		}
	}
	for _, host := range stopping {
		host.signalStop()
	}
	for _, host := range stopping {
		host.waitStopped()
	}

	c.syncHistory(cfg)

	c.Hosts = nil
	if !c.running {
		return
	}
	for _, hostConfig := range cfg.Hosts {
		host, ok := kept[hostConfig.ID]
		if !ok {
			host = NewHost(hostConfig.ID, c.pingServer, c.GetHostHistory(hostConfig.ID))
			host.UpdateConfig()
			host.Start()
		}
		c.Hosts = append(c.Hosts, host)
	}
}

// syncHistory keeps the history of the hosts in the config, loading it for new hosts
func (c *System) syncHistory(cfg *config.Config) {
	c.historyCleanup.Do(removeStaleHistoryFiles)

	c.historyMtx.Lock()
	defer c.historyMtx.Unlock()
	history := make(map[string]*HostHistory)
	for _, hostConfig := range cfg.Hosts {
		h, ok := c.history[hostConfig.ID]
		if !ok {
			h = NewHostHistory(hostConfig.ID)
		}
		history[hostConfig.ID] = h
	}
	for id, h := range c.history {
		if _, ok := history[id]; !ok {
			h.Close()
		}
	}
	c.history = history
}
