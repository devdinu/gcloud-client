package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func AddSSHKeyCmd(inst, ssh_key_path string, cfg Config) gcloudCommand {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    fmt.Sprintf("compute instances add-metadata %s --metadata-from-file ssh-keys=%s", inst, ssh_key_path),
		config: cfg,
	}
}

func GetInstancesCmd(cfg Config) gcloudCommand {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    "compute instances list",
		config: cfg}
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
