package config

import (
	"fmt"
)

// Config api/db shared config object
type Config struct {
	// 'api' or 'db'
	Role string

	// address of db
	DBAddr string

	// port of db
	DBPort string
}

// DBURL convenience getter
func (config Config) DBURL() string {
	return fmt.Sprintf("%s:%s", config.DBAddr, config.DBPort)
}
