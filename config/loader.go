package config

import (
	"flag"
	"os"
)

type Args struct {
	Zone, Format, InstanceName, User, Filter string
	AddHosts                                 bool
	Limit                                    int
	SSHFile                                  string
	Action                                   string
	InstanceCmdArgs
}

type InstanceCmdArgs struct {
	Refresh bool
}

var args Args

func Load() {
	var instanceArgs InstanceCmdArgs
	var sshFile, instanceName, zone, filter string
	var addHosts bool

	sshCommand := flag.NewFlagSet("ssh_access", flag.ExitOnError)
	instanceCommand := flag.NewFlagSet("instances", flag.ExitOnError)

	defaultSSHFile := os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	sshCommand.StringVar(&sshFile, "ssh_key", defaultSSHFile, "new SSH Key file which have to be added to instances")
	sshCommand.StringVar(&filter, "filter", "", "regexp to filter instances")
	sshCommand.StringVar(&args.User, "user", "", "username to add ssh key, if empty $USER will be taken")
	sshCommand.StringVar(&args.InstanceName, "instance", "", "instance to add ssh key, take precedence over the regex filter, would require zone")
	sshCommand.BoolVar(&addHosts, "add_hosts", false, "to add ip host mappings in /etc/hosts")

	flag.StringVar(&zone, "zone", "", "zone in which the given instance is present")
	flag.IntVar(&args.Limit, "limit", 1, "limit number of instances to add")

	instanceCommand.BoolVar(&instanceArgs.Refresh, "refresh", true, "refresh instances list in store")

	flag.Parse()

	if len(os.Args) >= 2 {
		if os.Args[1] == "ssh_access" || os.Args[1] == "" {
			sshCommand.Parse(os.Args[2:])
			args.Action = "AddSSHKeys"
		} else if os.Args[1] == "instances" {
			instanceCommand.Parse(os.Args[2:])
			args.Action = "RefreshInstances"
		}
	}

	args = Args{
		Zone:            zone,
		Format:          "json",
		Filter:          filter,
		AddHosts:        addHosts,
		InstanceName:    instanceName,
		SSHFile:         sshFile,
		InstanceCmdArgs: instanceArgs,
	}
}

func GetInstanceCmdArgs() InstanceCmdArgs { return args.InstanceCmdArgs }
func GetArgs() Args                       { return args }
