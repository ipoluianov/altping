package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
)

func getLastWorkspaceID() string {
	configPath := ConfigDirectory()
	fullPath := path.Join(configPath, "last_workspace.txt")
	bs, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}
	return string(bs)
}

func SetLastWorkspaceID(id string) error {
	configPath := ConfigDirectory()
	fullPath := path.Join(configPath, "last_workspace.txt")
	mkdirallErr := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if mkdirallErr != nil {
		return mkdirallErr
	}
	return os.WriteFile(fullPath, []byte(id), 0644)
}

func LoadLastWorkspace() (*Config, error) {
	id := getLastWorkspaceID()
	if id == "" {
		return nil, nil
	}
	var c Config
	err := c.Load(id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GenerateRandomID() string {
	rndbytes := make([]byte, 8)
	_, err := rand.Read(rndbytes)
	if err != nil {
		return ""
	}
	id := hex.EncodeToString(rndbytes)
	return id
}

func ConfigDirectory() string {
	localExePath := ""
	localExePath, err := os.Executable()
	if err != nil {
		localExePath = "."
	}
	configPath := filepath.Join(filepath.Dir(localExePath), ".altping")
	return configPath
}

func loadDefaultConfig() {
	configDir := ConfigDirectory()

	configIdLast := ""
	bsLastConfigId, err := os.ReadFile(path.Join(configDir, "last_config_id.txt"))
	if err == nil {
		configIdLast = string(bsLastConfigId)
	}

	cfg := NewConfig()
	err = cfg.Load(configIdLast)
	if err != nil {
		cfg, _ = CreateNewConfig("Default")
		if cfg == nil {
			cfg = NewConfig()
		}
	}
	currentConfig = cfg
}

func Init() {
	loadDefaultConfig()
}

func Get() *Config {
	return currentConfig
}

func SaveLastConfigId() {
	if currentConfig == nil {
		return
	}
	configPath := ConfigDirectory()
	fullPath := path.Join(configPath, "last_config_id.txt")
	_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
	_ = os.WriteFile(fullPath, []byte(currentConfig.ID), 0644)
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

func CreateNewConfig(name string) (*Config, error) {
	var err error
	id := generateId()
	config := NewConfig()
	config.ID = id
	config.Name = name
	config.Hosts = make([]*ConfigHost, 0)
	err = config.Save()
	if err != nil {
		return config, err
	}
	return config, nil
}

func RemoveConfig(id string) error {
	configPath := ConfigDirectory()
	fullPath := path.Join(configPath, id+".ws")
	err := os.Remove(fullPath)
	if err != nil {
		return err
	}
	return nil
}

func Configs() []*Config {
	configPath := ConfigDirectory()
	files, err := os.ReadDir(configPath)
	if err != nil {
		return []*Config{}
	}
	var configs []*Config
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".ws" {
			ws := &Config{}
			ws.Load(file.Name()[:len(file.Name())-3])
			configs = append(configs, ws)
		}
	}
	return configs
}

func generateId() string {
	rndbytes := make([]byte, 8)
	_, err := rand.Read(rndbytes)
	if err != nil {
		return "default"
	}
	id := hex.EncodeToString(rndbytes)
	return id
}
