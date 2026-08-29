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

	// Build a map of existing hosts for quick lookup
	existingHosts := make(map[string]*Host)
	for _, host := range c.Hosts {
		existingHosts[host.ID] = host
	}

	// Update existing hosts or add new ones based on the config
	for _, hostConfig := range config.Hosts {
		var host *Host
		for _, h := range c.Hosts {
			if h.ID == hostConfig.ID {
				host = h
				break
			}
		}
		if host == nil {
			host = NewHost(hostConfig.ID, c.pingServer)
			c.Hosts = append(c.Hosts, host)
		}
		host.UpdateConfig()
		delete(existingHosts, hostConfig.ID)
	}

	// Remove hosts that are no longer in the config
	for _, host := range existingHosts {
		c.removeHost(host.ID)
	}
}

// removeHost removes a host from the system by its ID. It stops the host and removes it from the Hosts slice.
func (c *System) removeHost(id string) {
	for i, h := range c.Hosts {
		if h.ID == id {
			h.Stop()
			c.Hosts = append(c.Hosts[:i], c.Hosts[i+1:]...)
			break
		}
	}
}
