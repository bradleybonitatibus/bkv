package config

import "fmt"

// Config contains the BKV config properties
type Config struct {
	GRPCPort int
	GRPCHost string
}

// NewConfig initializes and returns a pointer to a config struct
func NewConfig(port int, host string) *Config {
	return &Config{
		GRPCPort: port,
		GRPCHost: host,
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("%v:%v", c.GRPCHost, c.GRPCPort)
}
