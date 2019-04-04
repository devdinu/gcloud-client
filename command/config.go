package command

import (
	"fmt"
)

type Config struct {
	Zone    string
	Format  string
	Limit   int
	Filter  string
	Project string
}

func (c *Config) Flags() []string {
	var flags []string
	if c.Zone != "" {
		flags = append(flags, "--zone", c.Zone)
	}
	if c.Format != "" {
		flags = append(flags, "--format", c.Format)
	}
	if c.Limit != 0 {
		flags = append(flags, "--limit", fmt.Sprintf("%d", c.Limit))
	}
	if c.Filter != "" {
		flags = append(flags, "--filter", fmt.Sprintf("name~'%s'", c.Filter))
	}
	if c.Project != "" {
		flags = append(flags, "--project", c.Project)
	}
	return flags
}
