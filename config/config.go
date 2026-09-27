package config

import (
	"encoding/json"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

type ConfigHost struct {
	ID          string
	DisplayName string
	Hostname    string
	Port        int `json:",omitempty"` // TCP port to connect to instead of ping; 0 - ping
	IntervalMs  int `json:",omitempty"` // between pings; 0 - DefaultIntervalMs
	TimeoutMs   int `json:",omitempty"` // to wait for a reply; 0 - DefaultTimeoutMs
	SlowMs      int `json:",omitempty"` // average above it is shown as slow; 0 - off
	// Notify: beep and mark the window when the host goes down or comes back
	Notify bool `json:",omitempty"`
}

// Target returns what to check: the host and the TCP port ("" - ping).
// An old "example.com:443" in Hostname still means the port.
func (c ConfigHost) Target() (host string, port string) {
	if c.Port > 0 {
		return c.Hostname, strconv.Itoa(c.Port)
	}
	if h, p, err := net.SplitHostPort(c.Hostname); err == nil {
		return h, p
	}
	return c.Hostname, ""
}

// Address is the host with the port, if one is checked: "example.com:443", "[2001:db8::1]:443"
func (c ConfigHost) Address() string {
	host, port := c.Target()
	if port == "" {
		return host
	}
	return net.JoinHostPort(host, port)
}

// Interval returns the time between pings of the host
func (c ConfigHost) Interval() time.Duration {
	ms := DefaultIntervalMs
	if c.IntervalMs > 0 {
		ms = max(MinIntervalMs, min(c.IntervalMs, MaxIntervalMs))
	}
	return time.Duration(ms) * time.Millisecond
}

// Timeout returns how long to wait for a reply from the host
func (c ConfigHost) Timeout() time.Duration {
	ms := DefaultTimeoutMs
	if c.TimeoutMs > 0 {
		ms = max(MinTimeoutMs, min(c.TimeoutMs, MaxTimeoutMs))
	}
	return time.Duration(ms) * time.Millisecond
}

type Config struct {
	ID    string
	Name  string
	Hosts []*ConfigHost
}

var currentConfig *Config

func NewConfig() *Config {
	var c Config
	return &c
}

func (c *Config) AddHost(host ConfigHost) {
	c.Hosts = append(c.Hosts, &host)
}

func (c *Config) RemoveHost(id string) {
	for i, host := range c.Hosts {
		if host.ID == id {
			c.Hosts = append(c.Hosts[:i], c.Hosts[i+1:]...)
			return
		}
	}
}

func (c *Config) GetHost(id string) ConfigHost {
	for _, host := range c.Hosts {
		if host.ID == id {
			return *host
		}
	}
	var empty ConfigHost
	return empty
}

func (c *Config) Save() error {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	fullPath := path.Join(ConfigDirectory(), c.ID+".ws")
	mkdirallErr := os.MkdirAll(ConfigDirectory(), 0755)
	if mkdirallErr != nil {
		return mkdirallErr
	}
	return os.WriteFile(fullPath, bs, 0644)
}

func (c *Config) Load(id string) error {
	fullPath := path.Join(ConfigDirectory(), id+".ws")
	bs, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, c)
	if err != nil {
		return err
	}
	return nil
}
