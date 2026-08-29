package config

type ConfigHost struct {
	ID          string
	DisplayName string
	Hostname    string
	TimeoutMs   int
	DataSize    int
}

type Config struct {
	Hosts []*ConfigHost
}

func Get() *Config {
	var c Config
	c.Hosts = make([]*ConfigHost, 0)
	{
		var host ConfigHost
		host.ID = "001"
		host.DisplayName = "LocalHost"
		host.Hostname = "localhost"
		host.TimeoutMs = 1000
		host.DataSize = 64
		c.Hosts = append(c.Hosts, &host)
	}
	{
		var host ConfigHost
		host.ID = "002"
		host.DisplayName = "Example"
		host.Hostname = "example.com"
		host.TimeoutMs = 2000
		host.DataSize = 128
		c.Hosts = append(c.Hosts, &host)
	}
	{
		var host ConfigHost
		host.ID = "003"
		host.DisplayName = ""
		host.Hostname = "192.0.2.1"
		host.TimeoutMs = 1000
		host.DataSize = 64
		c.Hosts = append(c.Hosts, &host)
	}
	{
		var host ConfigHost
		host.ID = "004"
		host.DisplayName = ""
		host.Hostname = "altbins.com"
		host.TimeoutMs = 2000
		host.DataSize = 64
		c.Hosts = append(c.Hosts, &host)
	}
	{
		var host ConfigHost
		host.ID = "005"
		host.DisplayName = ""
		host.Hostname = "altbins.pro"
		host.TimeoutMs = 2000
		host.DataSize = 64
		c.Hosts = append(c.Hosts, &host)
	}
	return &c
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
