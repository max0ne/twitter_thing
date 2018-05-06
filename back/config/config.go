package config

import (
	"fmt"
	"regexp"
)

// Config api/db shared config object
type Config struct {
	// 'api' or 'db'
	Role string

	// address of db
	DBAddr string

	// port of db
	DBPort string

	// port of vr service
	VRPort string

	// list of vr peer urls
	VRPeerURLs []string

	// list of peer urls of dbs
	DBPeerURLs []string
}

// DBURL convenience getter
func (config Config) DBURL() string {
	return fmt.Sprintf("%s:%s", config.DBAddr, config.DBPort)
}

// VRURL convenience getter
func (config Config) VRURL() string {
	return fmt.Sprintf("%s:%s", config.DBAddr, config.VRPort)
}

// Validate - -
func (config Config) Validate() error {
	if config.Role == "api" {
		if !regexp.MustCompile("^.*?:.*?$").Match([]byte(config.DBURL())) {
			return fmt.Errorf("misconfigured DBAddr or DBPort")
		}
		return nil
	}
	if config.Role == "db" {
		if !regexp.MustCompile("^.*?:.*?$").Match([]byte(config.VRURL())) {
			return fmt.Errorf("misconfigured DBAddr or VRPort")
		}
		if config.VRMe() == -1 {
			return fmt.Errorf("VRURL %s must be in VRPeerURLs", config.VRURL())
		}
		return nil
	}
	return fmt.Errorf("unknown role %s", config.Role)
}

// VRMe index of local machine in vr peer array
func (config Config) VRMe() int {
	for idx, url := range config.VRPeerURLs {
		if url == config.VRURL() {
			return idx
		}
	}
	return -1
}
