package command

import (
	"fmt"
	"strings"
)

type gcloudCommand struct {
	name   string
	cmd    string
	config Config
	args   []string
}

func (cmd gcloudCommand) Name() string {
	return cmd.name
}

func (cmd gcloudCommand) Args() []string {
	cmdslice := strings.Fields(cmd.cmd)
	cmdslice = append(cmdslice, cmd.config.Flags()...)
	return cmdslice
}

func (cmd gcloudCommand) String() string {
	return fmt.Sprintf("%s %s %s", cmd.name, cmd.cmd, strings.Join(cmd.config.Flags(), " "))
}

type Command interface {
	Name() string
	Args() []string
	String() string
}
