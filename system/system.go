package system

import "github.com/ipoluianov/altping/config"

type System struct {
	pingServer  *PingServer
	Hosts       []*Host
	chanStopped chan struct{}
}

var systemInstance *System

func init() {
	systemInstance = newSystem()
}

func newSystem() *System {
	var c System
	c.pingServer = NewPingServer()
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

func (c *System) UpdateConfig() {
	config := config.Get()
	c.Hosts = nil
	for _, hostConfig := range config.Hosts {
		var host *Host
		host = NewHost(hostConfig.ID, c.pingServer)
		c.Hosts = append(c.Hosts, host)
		host.UpdateConfig()
	}
}
