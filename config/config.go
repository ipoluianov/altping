package config

type ConfigHost struct {
	ID        string
	Hostname  string
	TimeoutMs int
	DataSize  int
}

type Config struct {
	Hosts []*ConfigHost
}

func Get() *Config {
	var c Config
	c.Hosts = make([]*ConfigHost, 0)
	/*{
		var host ConfigHost
		host.ID = "001"
		host.Hostname = "localhost"
		host.TimeoutMs = 1000
		host.DataSize = 64
		c.Hosts = append(c.Hosts, &host)
	}*/
	{
		var host ConfigHost
		host.ID = "002"
		host.Hostname = "example.com"
		host.TimeoutMs = 2000
		host.DataSize = 128
		c.Hosts = append(c.Hosts, &host)
	}
	/*{
		var host ConfigHost
		host.ID = "003"
		host.Hostname = "192.168.254.12"
		host.TimeoutMs = 1000
		host.DataSize = 64
		c.Hosts = append(c.Hosts, &host)
	}*/
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
