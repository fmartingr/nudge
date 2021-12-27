package server

import "fmt"

type Config struct {
	Port     int
	LogLevel string

	IPs      []string
	Interval int // seconds
}

func (c *Config) Print() {
	fmt.Printf("port = %d\n", c.Port)
	fmt.Printf("loglevel = %s\n", c.LogLevel)
	fmt.Printf("ips = %s\n", c.IPs)
	fmt.Printf("interval = %d\n", c.Interval)
}

func NewDefaultConfig() *Config {
	return &Config{
		Port:     3000,
		LogLevel: "warn",
		IPs:      []string{"9.9.9.9", "1.1.1.1"},
		Interval: 60,
	}
}
