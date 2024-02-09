package config

import (
	"os"

	logger "github.com/invest-scraping/logg"
	"gopkg.in/yaml.v3"
)

const (
	defaultMonitorsDir    = "./config/env/monitors/"
	stagingMonitorsDir    = "./config/env/monitors/"
	productionMonitorsDir = "./config/env/monitors/"
)

var log = logger.NewDefaultLog()

func RegisterMonitors(cfg *Config) {
	monitorsDir := getMonitorDir()
	files, err := os.ReadDir(monitorsDir)
	if err != nil {
		log.Panic("error reading monitor file names. Reason: ", err.Error())
	}

	var monitors []*Monitor
	for _, f := range files {
		b, err := os.ReadFile(monitorsDir + f.Name())
		if err != nil {
			log.Panic("error reading monitor configurations. Reason: ", err.Error())
		}

		m := &Monitor{}
		if err := yaml.Unmarshal(b, m); err != nil {
			log.Panic("error unmarshalling monitor configuration. Reason: ", err.Error())

		}
		monitors = append(monitors, m)
	}

	cfg.Monitors = monitors
}

func getMonitorDir() string {
	env := os.Getenv(envProfileKey)
	dir := ""

	switch {
	case env == "staging":
		dir = stagingMonitorsDir
	case env == "production":
		dir = productionMonitorsDir
	default:
		dir = defaultMonitorsDir
	}

	if _, err := os.Stat(dir); err != nil {
		return defaultMonitorsDir
	}

	return dir
}

func findMonitorByKey(key string, monitors []*Monitor) *Monitor {
	for _, monitor := range monitors {
		if monitor.Key == key {
			return monitor
		}
	}
	return nil
}
