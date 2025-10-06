package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	MegaCLIPath string
	Port        string
	LogLevel    string
}

// Common MegaCLI installation paths
var DefaultMegaCLIPaths = []string{
	"/opt/MegaRAID/MegaCli/MegaCli64",
	"/opt/MegaCli/MegaCli64",
	"/usr/sbin/MegaCli64",
	"/usr/local/bin/MegaCli64",
	"/opt/lsi/MegaCLI/MegaCli64",
}

func NewConfig() *Config {
	return &Config{
		Port:     "8080",
		LogLevel: "info",
	}
}

func (c *Config) SetMegaCLIPath(path string) {
	c.MegaCLIPath = path
}

func (c *Config) GetMegaCLIPath() string {
	if c.MegaCLIPath != "" {
		return c.MegaCLIPath
	}
	
	// Try to discover MegaCLI automatically
	if path := DiscoverMegaCLI(); path != "" {
		c.MegaCLIPath = path
		return path
	}
	
	return ""
}

func DiscoverMegaCLI() string {
	for _, path := range DefaultMegaCLIPaths {
		if IsValidMegaCLI(path) {
			return path
		}
	}
	return ""
}

func IsValidMegaCLI(path string) bool {
	if path == "" {
		return false
	}
	
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	
	// Check if file is executable
	if err := checkExecutable(path); err != nil {
		return false
	}
	
	return true
}

func checkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	
	mode := info.Mode()
	if mode&0111 == 0 {
		return os.ErrPermission
	}
	
	return nil
}
