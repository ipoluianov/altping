package system

import (
	"sync"

	"github.com/ipoluianov/altping/config"
)

type System struct {
	pingServer  *PingServer
	Hosts       []*Host
	chanStopped chan struct{}

	// Ping history by host ID. Hosts are recreated on every restart,
	// so the history is kept here.
	historyMtx sync.Mutex
	history    map[string]*HostHistory
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
	c.UpdateConfig()
	c.chanStopped = make(chan struct{})
	c.pingServer.Start()
	for _, host := range c.Hosts {
		host.Start(c.chanStopped)
	}
}

func (c *System) Stop() {
	if c.chanStopped == nil {
		return
	}

	c.pingServer.Stop()
	if c.chanStopped != nil {
		close(c.chanStopped)
		c.chanStopped = nil
	}
	// Signal all hosts to stop
	for _, host := range c.Hosts {
		go host.Stop()
	}
	allStopped := false
	for !allStopped {
		allStopped = true
		for _, host := range c.Hosts {
			if host.IsRunning() {
				allStopped = false
				break
			}
		}
	}
	c.Hosts = nil
}

func (c *System) IsRunning() bool {
	if c.chanStopped == nil {
		return false
	}
	return true
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

func (c *System) UpdateConfig() {
	config := config.Get()
	c.Hosts = nil

	c.historyMtx.Lock()
	history := make(map[string]*HostHistory)
	for _, hostConfig := range config.Hosts {
		h, ok := c.history[hostConfig.ID]
		if !ok {
			h = NewHostHistory()
		}
		history[hostConfig.ID] = h
	}
	c.history = history
	c.historyMtx.Unlock()

	for _, hostConfig := range config.Hosts {
		var host *Host
		host = NewHost(hostConfig.ID, c.pingServer, history[hostConfig.ID])
		c.Hosts = append(c.Hosts, host)
		host.UpdateConfig()
	}
}
