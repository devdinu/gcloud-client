package action

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
)

func AddSSHKeys(c gcloud.Client, args config.Args) error {
	var insts []gcloud.Instance
	if args.InstanceName == "" || args.Zone == "" {
		var err error
		cfg := command.Config{Zone: args.Zone, Format: args.Format, Filter: args.Filter}
		if args.Projects != "" {
			cfg.Project = config.Projects()[0]
		}
		insts, err = c.GetInstances(cfg)
		if err != nil {
			fmt.Println(fmt.Errorf("get instances errored %v", err))
			return err
		}
	} else {
		insts = []gcloud.Instance{{Name: args.InstanceName, Zone: args.Zone}}
	}
	for _, inst := range insts {
		conf := command.Config{Format: args.Format, Zone: inst.Zone}
		if args.Projects != "" {
			conf.Project = config.Projects()[0]
		}
		desc, err := c.GetDescription(inst.Name, conf)
		if err != nil {
			fmt.Println(fmt.Errorf("describe instance errored %v", err))
			return err
		}
		keys := desc.SSHKeys()
		newKey, err := readKey(args.User, args.SSHFile)
		if err != nil {
			fmt.Printf("Error adding key to instance %s err: %v\n", inst.Name, err)
			return err
		}
		keys = append(keys, newKey)
		out, err := c.AddSSHKeys(inst.Name, command.Config{Zone: inst.Zone}, keys)
		if err != nil {
			fmt.Printf("Error adding key to instance %s err: %v\n", inst.Name, err)
			return err
		}
		fmt.Printf("Added key to instance: %s ip: %s status: %s %s\n", inst.Name, inst.IP(), inst.Status, out)
	}
	return nil
}

func readKey(user, filename string) (gcloud.SSHKey, error) {
	f, err := os.Open(filename)
	if err != nil {
		return gcloud.SSHKey{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return gcloud.SSHKey{}, err
	}
	key := string(b)
	fields := strings.Fields(key)
	if len(fields) != 3 {
		return gcloud.SSHKey{}, errors.New("Invalid SSH Key Format")
	}
	if user == "" {
		user = os.Getenv("USER")
	}
	return gcloud.SSHKey{Username: user + ":" + fields[0], Key: fields[1], ID: fields[2]}, nil
}
