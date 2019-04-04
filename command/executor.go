package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/devdinu/gcloud-client/logger"
)

type Executor struct{}

func (e Executor) Execute(c Command) (io.Reader, error) {
	var out bytes.Buffer
	execCmd := exec.Command(c.Name(), c.Args()...)
	logger.Debugf("[Executor] Executing command: %s", c.String())
	execCmd.Stdout = &out
	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr
	err := execCmd.Run()
	return &out, err
}
