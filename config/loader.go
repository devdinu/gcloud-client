package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type InstanceCmdArgs struct {
	Prefix  string
	Refresh bool
}

type Login struct {
	Session string
	User    string
}

type Args struct {
	Zone, Format, InstanceName, User, Filter string
	AddHosts                                 bool
	Limit                                    int
	DBFile, SSHFile                          string
	InstanceCmdArgs
	Login
	LogLevel string
	Projects string
}

type CmdAction string

const SshAccess CmdAction = "ssh_access"
const RefreshInstances CmdAction = "refresh"
const LoginInstances CmdAction = "login"
const SearchPrefix CmdAction = "prefix_search"

var args Args
var cmdAction CmdAction

func Load() {
	var instanceArgs InstanceCmdArgs

	sshCommand := flag.NewFlagSet("ssh_access", flag.ContinueOnError)
	instanceCommand := flag.NewFlagSet("instances", flag.ContinueOnError)

	defaultSSHFile := os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	defaultUser := os.Getenv("USER")
	defaultDBFile := os.Getenv("HOME") + "/hosts.db"

	sshCommand.StringVar(&args.SSHFile, "ssh_key", defaultSSHFile, "new SSH Key file which have to be added to instances")
	sshCommand.StringVar(&args.Filter, "filter", "", "regexp to filter instances")
	sshCommand.StringVar(&args.InstanceName, "instance", "", "instance to add ssh key, take precedence over the regex filter, would require zone")
	sshCommand.BoolVar(&args.AddHosts, "add_hosts", false, "to add ip host mappings in /etc/hosts")
	sshCommand.StringVar(&args.Zone, "zone", "", "zone in which the given instance is present")

	// refresh should be subcommand and not as flag
	instanceCommand.BoolVar(&instanceArgs.Refresh, "refresh", true, "refresh instances list in store")
	instanceCommand.StringVar(&instanceArgs.Prefix, "prefix", "", "search instances by common prefix")
	instanceCommand.StringVar(&instanceArgs.Prefix, "regex", "", "search instances by regex")
	instanceCommand.StringVar(&args.Login.Session, "session", "login-session", "login sesssion name")

	sshCommand.StringVar(&args.DBFile, "dbfile", defaultDBFile, "db file to store data")
	instanceCommand.StringVar(&args.DBFile, "dbfile", defaultDBFile, "db file to store data")
	instanceCommand.StringVar(&args.LogLevel, "level", "info", "log level [info/debug/all]")
	sshCommand.StringVar(&args.LogLevel, "level", "info", "log level [info/debug/all]")
	instanceCommand.IntVar(&args.Limit, "limit", 0, "limit number of instances to search")
	sshCommand.IntVar(&args.Limit, "limit", 0, "limit number of instances to add")
	sshCommand.StringVar(&args.User, "user", defaultUser, "username to add ssh key, if empty $USER will be taken")
	instanceCommand.StringVar(&args.Login.User, "user", defaultUser, "username for ssh")
	sshCommand.StringVar(&args.Projects, "projects", "", "projects to search for as comma seperated values: proj1,proj2,lastproject (project id)")
	instanceCommand.StringVar(&args.Projects, "projects", "", "projects to search for as comma seperated values: proj1,proj2,lastproject (project id)")

	flag.Parse()
	//sshCommand.SetOutput(ioutil.Discard)
	//instanceCommand.SetOutput(ioutil.Discard)

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
				if err := instanceCommand.Parse(os.Args[3:]); err != nil {
					log.Fatalf("[Config] Error defining instances command %v", err)
				}
				cmdAction = RefreshInstances
			case "login":
				if err := instanceCommand.Parse(os.Args[3:]); err != nil {
					log.Fatalf("[Config] Error defining instances command %v", err)
				}
				cmdAction = LoginInstances
			default:
				fmt.Println("[Config] no matching commands mentioned")
				flag.Usage()
			}
		}
	}

	args.InstanceCmdArgs = instanceArgs
	args.Format = "json"
	fmt.Printf("action %s args: %+v \ncmd args: %v\n", cmdAction, args, os.Args)
}

func GetInstanceCmdArgs() InstanceCmdArgs { return args.InstanceCmdArgs }
func GetArgs() Args                       { return args }
func GetActionName() CmdAction            { return cmdAction }
func GetDBFileName() string               { return args.DBFile }
func LogLevel() string                    { return args.LogLevel }
func Projects() []string {
	if args.Projects != "" {
		return strings.Split(args.Projects, ",")
	}
	return []string{}
}
