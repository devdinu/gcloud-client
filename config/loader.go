package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type InstanceCmdArgs struct {
	Prefix  string
	Refresh bool
}

type Args struct {
	Zone, Format, InstanceName, User, Filter string
	AddHosts                                 bool
	Limit                                    int
	DBFile, SSHFile                          string
	InstanceCmdArgs
}

type CmdAction string

const SshAccess CmdAction = "ssh_access"
const RefreshInstances CmdAction = "refresh"
const SearchPrefix CmdAction = "prefix_search"

var args Args
var cmdAction CmdAction

func Load() {
	var instanceArgs InstanceCmdArgs
	var sshFile, instanceName, zone, filter string
	var addHosts bool
	var dbfile string

	sshCommand := flag.NewFlagSet("ssh_access", flag.ContinueOnError)
	instanceCommand := flag.NewFlagSet("instances", flag.ContinueOnError)

	defaultSSHFile := os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	sshCommand.StringVar(&sshFile, "ssh_key", defaultSSHFile, "new SSH Key file which have to be added to instances")
	sshCommand.StringVar(&filter, "filter", "", "regexp to filter instances")
	sshCommand.StringVar(&args.User, "user", "", "username to add ssh key, if empty $USER will be taken")
	sshCommand.StringVar(&args.InstanceName, "instance", "", "instance to add ssh key, take precedence over the regex filter, would require zone")
	sshCommand.BoolVar(&addHosts, "add_hosts", false, "to add ip host mappings in /etc/hosts")

	sshCommand.StringVar(&zone, "zone", "", "zone in which the given instance is present")
	sshCommand.StringVar(&dbfile, "dbfile", "hosts.db", "db file to store data")
	sshCommand.IntVar(&args.Limit, "limit", 1, "limit number of instances to add")

	// refresh should be subcommand and not as flag
	instanceCommand.BoolVar(&instanceArgs.Refresh, "refresh", true, "refresh instances list in store")
	instanceCommand.StringVar(&instanceArgs.Prefix, "prefix", "", "search instances by common prefix")
	instanceCommand.StringVar(&dbfile, "dbfile", "hosts.db", "db file to store data")

	flag.Parse()
	//sshCommand.SetOutput(ioutil.Discard)
	//instanceCommand.SetOutput(ioutil.Discard)

	fmt.Printf("parse success: %v val: %s \nflagargs: %v \nosArgs:%v %d\n", flag.Parsed(), dbfile, flag.Args(), os.Args, len(os.Args))

	if len(os.Args) >= 2 {
		if os.Args[1] == "ssh_access" || os.Args[1] == "" {
			if err := sshCommand.Parse(os.Args[2:]); err != nil {
				log.Fatalf("[Config] Error defining ssh access command %v", err)
			}
			cmdAction = SshAccess
		} else if os.Args[1] == "instances" {
			if len(os.Args) < 3 {
				log.Fatalf("[Config] Error defining ssh access command no action mentioned")
			}
			switch os.Args[2] {
			case "search":
				if err := instanceCommand.Parse(os.Args[3:]); err != nil {
					log.Fatalf("[Config] Error defining instances command %v", err)
				}
				if instanceCommand.Parsed() {
					if instanceArgs.Prefix != "" {
						cmdAction = SearchPrefix
					}
				}
			case "refresh":
				cmdAction = RefreshInstances
			}
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
		DBFile:          dbfile,
	}
	fmt.Printf("action %s args: %+v \ncmd args: %v\n", cmdAction, args, os.Args)
}

func GetInstanceCmdArgs() InstanceCmdArgs { return args.InstanceCmdArgs }
func GetArgs() Args                       { return args }
func GetActionName() CmdAction            { return cmdAction }
func GetDBFileName() string               { return args.DBFile }
