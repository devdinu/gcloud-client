package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type executor interface {
	Execute(c Command) (io.Reader, error)
}

// Should add mock and test it with sample.json
func (e commandExecutor) Execute(c Command) (io.Reader, error) {
	var out bytes.Buffer
	execCmd := exec.Command(c.Name(), c.Args()...)
	fmt.Println("Executing command: ", c.String())
	execCmd.Stdout = &out
	err := execCmd.Run()
	return &out, err
}

type commandExecutor struct{}
