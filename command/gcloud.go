package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
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

type Executor struct{}

func (e Executor) Execute(c Command) (io.Reader, error) {
	var out bytes.Buffer
	execCmd := exec.Command(c.Name(), c.Args()...)
	fmt.Println("Executing command: ", c.String())
	execCmd.Stdout = &out
	err := execCmd.Run()
	return &out, err
}
