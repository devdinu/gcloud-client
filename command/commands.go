package command

import (
	"fmt"
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

func ListProjects(cfg Config) gcloudCommand {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    "projects list",
		config: cfg,
	}
}

func DescribeCmd(inst string, cfg Config) gcloudCommand {
	return gcloudCommand{
		name:   "gcloud",
		cmd:    "compute instances describe " + inst,
		config: cfg,
	}
}

func Login(hosts []string, cmd string, cfg TmuxConfig) Command {
	if len(hosts) <= 0 {
		//TODO: introduce nop command
		return &tmux{}
	}
	return &tmux{
		hosts:      hosts,
		cmd:        cmd,
		TmuxConfig: cfg,
	}
}
