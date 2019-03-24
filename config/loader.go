package config

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Zone, Format, InstanceName, User, Filter string
	AddHosts                                 bool
	Limit                                    int
	SSHFile                                  string
}

var args Args

func Load() {
	var sshFile, instanceName, zone, filter string
	var addHosts bool

	flag.StringVar(&sshFile, "ssh_key", "", "new SSH Key file which have to be added to instances")
	flag.StringVar(&filter, "filter", "", "regexp to filter instances")
	flag.StringVar(&args.User, "user", "", "username to add ssh key, if empty $USER will be taken")
	flag.IntVar(&args.Limit, "limit", 1, "limit number of instances to add")
	flag.StringVar(&instanceName, "instance", "", "instance to add ssh key, take precedence over the regex filter, would require zone")
	flag.StringVar(&zone, "zone", "", "zone in which the given instance is present")
	flag.BoolVar(&addHosts, "add_hosts", false, "to add ip host mappings in /etc/hosts")
	flag.Parse()

	if sshFile == "" {
		sshFile = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
		fmt.Println("Using default $HOME/.ssh/id_rsa.pub as ssh_key")
	}

	args = Args{
		Zone:         zone,
		Format:       "json",
		Filter:       filter,
		AddHosts:     addHosts,
		InstanceName: instanceName,
		SSHFile:      sshFile,
	}
}

func GetArgs() Args { return args }
