package main

import (
	"fmt"
	"strings"
)

type Command interface {
	Name() string
	Args() []string
	String() string
}

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

func DescribeCmd(inst string, cfg Config) Command {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    "compute instances describe " + inst,
		config: cfg,
	}
}

func GetInstancesCmd(cfg Config) Command {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    "compute instances list",
		config: cfg}
}

func AddSSHKeyCmd(inst, ssh_key_path string, cfg Config) Command {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    fmt.Sprintf("compute instances add-metadata %s --metadata-from-file ssh-keys=%s", inst, ssh_key_path),
		config: cfg,
	}
}
