package main

import "fmt"

type Config struct {
	zone   string
	format string
	limit  int
}

func (c *Config) Zone() string {
	return c.zone
}

func (c *Config) Format() string {
	return c.format
}

func (c *Config) Limit() string {
	return fmt.Sprintf("%d", c.limit)
}

func (c *Config) Flags() []string {
	var flags []string
	if c.zone != "" {
		flags = append(flags, "--zone", c.zone)
	}
	if c.format != "" {
		flags = append(flags, "--format", c.format)
	}
	if c.limit != 0 {
		flags = append(flags, "--limit", fmt.Sprintf("%d", c.limit))
	}
	return flags
}
