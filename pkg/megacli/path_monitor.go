package megacli

import (
	"fmt"
	"log"
	"time"
	
	"github.com/manojsiriparthi/megaraid-exporter/config"
)

type PathMonitor struct {
	config     *config.Config
	isHealthy  bool
	lastCheck  time.Time
	checkInterval time.Duration
}

func NewPathMonitor(cfg *config.Config) *PathMonitor {
	return &PathMonitor{
		config:        cfg,
		checkInterval: 30 * time.Second,
	}
}

func (pm *PathMonitor) Start() {
	pm.checkPath()
	
	ticker := time.NewTicker(pm.checkInterval)
	go func() {
		for range ticker.C {
			pm.checkPath()
		}
	}()
}

func (pm *PathMonitor) checkPath() {
	path := pm.config.GetMegaCLIPath()
	pm.lastCheck = time.Now()
	
	if path == "" {
		pm.isHealthy = false
		log.Printf("WARNING: MegaCLI not found in any default locations")
		return
	}
	
	if config.IsValidMegaCLI(path) {
		if !pm.isHealthy {
			log.Printf("INFO: MegaCLI found and accessible at: %s", path)
		}
		pm.isHealthy = true
	} else {
		if pm.isHealthy {
			log.Printf("ERROR: MegaCLI path is no longer valid: %s", path)
		}
		pm.isHealthy = false
	}
}

func (pm *PathMonitor) IsHealthy() bool {
	return pm.isHealthy
}

func (pm *PathMonitor) GetPath() string {
	return pm.config.GetMegaCLIPath()
}

func (pm *PathMonitor) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"path":         pm.GetPath(),
		"healthy":      pm.isHealthy,
		"last_check":   pm.lastCheck.Format(time.RFC3339),
		"check_interval": pm.checkInterval.String(),
	}
}

func (pm *PathMonitor) ValidateOrError() error {
	if !pm.isHealthy {
		return fmt.Errorf("MegaCLI is not available at path: %s", pm.GetPath())
	}
	return nil
}
