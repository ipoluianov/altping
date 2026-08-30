package config

import (
	"encoding/json"
	"os"
	"path"
)

type ConfigHost struct {
	ID          string
	DisplayName string
	Hostname    string
}

type Config struct {
	ID    string
	Name  string
	Hosts []*ConfigHost
}

var currentConfig *Config

func init() {
	var c Config
	c.ID = "default"
	c.Name = "Default Config"
	c.Hosts = make([]*ConfigHost, 0)
	c.Load("default")
	currentConfig = &c
}

func Get() *Config {
	return currentConfig
}

func LoadConfig(id string) error {
	var c Config
	c.ID = id
	err := c.Load(id)
	if err != nil {
		return err
	}
	currentConfig = &c
	return nil
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
