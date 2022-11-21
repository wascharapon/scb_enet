package config

import (
	"net"

	"github.com/acoshift/configfile"
)

type Config struct {
	Port   string
	BindIP string
}

// Address Join bind ip and port
func (c Config) Address() string {
	return net.JoinHostPort(c.BindIP, c.Port)
}

// r reader
var r = configfile.NewEnvReader()

// store Config in global variable
var c Config

// Init Config from os environment
// or .env
func Init() *Config {
	configfile.LoadDotEnv()

	c.Port = r.StringDefault("PORT", "8080")
	return &c
}

// Load global Config data
func Load() *Config {
	return &c
}
