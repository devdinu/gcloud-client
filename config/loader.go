package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/devdinu/gcloud-client/logger"
)

type InstanceCmdArgs struct {
	Prefix  string
	Regex   string
	Refresh bool
}

type Login struct {
	Session      string
	User         string
	TemplatesDir string
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
const SearchRegex CmdAction = "regex_search"

var args Args
var cmdAction CmdAction

func MustLoad() {
	var instanceArgs InstanceCmdArgs
	configDir := os.Getenv("HOME") + string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "gcloud-client"
	configFile := configDir + string(os.PathSeparator) + "config.json"

	sshCommand := flag.NewFlagSet("ssh_access", flag.ContinueOnError)
	instanceCommand := flag.NewFlagSet("instances", flag.ContinueOnError)
	flag.StringVar(&configFile, "config", configFile, "configuration to load for gcl")

	//TODO: override with env "GCL_CONFIG_DIR"
	flag.Parse()

	defaultCfg, err := loadDefaults(configFile)
	if err != nil {
		log.Fatalf("couldn't load default config, err: %s", err)
	}

	sshCommand.StringVar(&args.SSHFile, "ssh_key", defaultCfg.SSHFile, "new SSH Key file which have to be added to instances")
	sshCommand.StringVar(&args.Filter, "filter", "", "regexp to filter instances")
	sshCommand.StringVar(&args.InstanceName, "instance", "", "instance to add ssh key, take precedence over the regex filter, would require zone")
	sshCommand.BoolVar(&args.AddHosts, "add_hosts", false, "to add ip host mappings in /etc/hosts")
	sshCommand.StringVar(&args.Zone, "zone", "", "zone in which the given instance is present")

	// refresh should be subcommand and not as flag
	instanceCommand.BoolVar(&instanceArgs.Refresh, "refresh", true, "refresh instances list in store")
	instanceCommand.StringVar(&instanceArgs.Prefix, "prefix", "", "search instances by common prefix")
	instanceCommand.StringVar(&instanceArgs.Regex, "regex", "", "search instances by regex")
	instanceCommand.StringVar(&args.Login.Session, "session", "login-session", "login sesssion name")
	instanceCommand.StringVar(&args.Login.TemplatesDir, "templates", defaultCfg.TemplatesDir, "templates directory for tmuxinator")

	sshCommand.StringVar(&args.DBFile, "dbfile", defaultCfg.DBFile, "db file to store data")
	instanceCommand.StringVar(&args.DBFile, "dbfile", defaultCfg.DBFile, "db file to store data")
	instanceCommand.StringVar(&args.LogLevel, "level", defaultCfg.LogLevel, "log level [info/debug/all]")
	sshCommand.StringVar(&args.LogLevel, "level", defaultCfg.LogLevel, "log level [info/debug/all]")
	instanceCommand.IntVar(&args.Limit, "limit", 0, "limit number of instances to search")
	sshCommand.IntVar(&args.Limit, "limit", 0, "limit number of instancddefaultDBFilees to add")
	sshCommand.StringVar(&args.User, "user", defaultCfg.User, "username to add ssh key, if empty $USER will be taken")
	instanceCommand.StringVar(&args.Login.User, "user", defaultCfg.User, "username for ssh")
	sshCommand.StringVar(&args.Projects, "projects", "", "projects to search for as comma seperated values: proj1,proj2,lastproject (project id)")
	instanceCommand.StringVar(&args.Projects, "projects", "", "projects to search for as comma seperated values: proj1,proj2,lastproject (project id)")

	//sshCommand.SetOutput(ioutil.Discard)
	//instanceCommand.SetOutput(ioutil.Discard)

	if len(os.Args) >= 2 {
		if os.Args[1] == "ssh_access" || os.Args[1] == "" {
			if err := sshCommand.Parse(os.Args[2:]); err != nil {
				sshCommand.Usage()
				log.Fatalf("[Config] Error defining ssh access command %v", err)
			}
			if args.Filter == "" && args.InstanceName == "" {
				log.Fatalf("[Config] mention instances search filter for access")
			}
			cmdAction = SshAccess
		} else if os.Args[1] == "instances" {
			if len(os.Args) < 3 {
				instanceCommand.Usage()
				log.Fatalf("[Config] Error defining instances command, no action mentioned")
			}
			switch os.Args[2] {
			case "search":
				if err := instanceCommand.Parse(os.Args[3:]); err != nil {
					log.Fatalf("[Config] Error defining instances command %v", err)
				}
				if instanceCommand.Parsed() {
					if instanceArgs.Regex != "" {
						cmdAction = SearchRegex
					} else if instanceArgs.Prefix != "" {
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
				instanceCommand.Usage()
			}
		}
	}
	if cmdAction == "" {
		flag.Usage()
		sshCommand.Usage()
		instanceCommand.Usage()
		log.Fatalf("Try again with proper command")
	}

	args.InstanceCmdArgs = instanceArgs
	args.Format = "json"
	logger.Debugf("action %s args: %+v \ncmd args: %v", cmdAction, args, os.Args)
}

func GetInstanceCmdArgs() InstanceCmdArgs { return args.InstanceCmdArgs }
func GetArgs() Args                       { return args }
func GetActionName() CmdAction            { return cmdAction }
func GetDBFileName() string               { return args.DBFile }

func LogLevel() string {
	if args.LogLevel == "" {
		return "debug"
	}
	return strings.ToLower(args.LogLevel)
}
func Projects() []string {
	if args.Projects != "" {
		return strings.Split(args.Projects, ",")
	}
	return []string{}
}
