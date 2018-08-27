package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var sshFile, instanceName, zone, filter, user string
	var limit int
	flag.StringVar(&sshFile, "ssh_key", "", "new SSH Key file which have to be added to instances")
	flag.StringVar(&filter, "filter", "", "regexp to filter instances")
	flag.StringVar(&user, "user", "", "username to add ssh key, if empty $USER will be taken")
	flag.IntVar(&limit, "limit", 1, "limit number of instances to add")
	flag.StringVar(&instanceName, "instance", "", "instance to add ssh key, take precedence over the regex filter, would require zone")
	flag.StringVar(&zone, "zone", "", "zone in which the given instance is present")
	flag.Parse()

	if sshFile == "" {
		sshFile = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
		fmt.Println("Using default $HOME/.ssh/id_rsa.pub as ssh_key")
	}
	var insts []instance
	cfg := Config{format: "json", limit: limit, filter: filter}
	if instanceName == "" || zone == "" {
		insts, _ = getInstances(cfg)
	} else {
		insts = []instance{{Name: instanceName, Zone: zone}}
	}
	for _, inst := range insts {
		conf := Config{format: "json", zone: inst.Zone}
		desc, err := getDescription(inst.Name, conf)
		if err != nil {
			fmt.Println(fmt.Errorf("describe instance errored %v", err))
			return
		}
		keys := desc.sshKeys()
		newKey, err := readKey(user, sshFile)
		if err != nil {
			fmt.Println("Error adding key to instance %s err: %v\n", inst.Name, err)
			return
		}
		keys = append(keys, newKey)
		out, err := inst.AddSSHKeys(Config{zone: inst.Zone}, keys)
		if err != nil {
			fmt.Printf("Error adding key to instance %s err: %v\n", inst.Name, err)
			return
		}
		fmt.Printf("Added key to instance: %s %s\n", inst.Name, out)
	}
}

func readKey(user, filename string) (SSHKey, error) {
	f, err := os.Open(filename)
	if err != nil {
		return SSHKey{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return SSHKey{}, err
	}
	key := string(b)
	fields := strings.Fields(key)
	if len(fields) != 3 {
		return SSHKey{}, errors.New("Invalid SSH Key Format")
	}
	if user == "" {
		user = os.Getenv("USER")
	}
	return SSHKey{username: user + ":" + fields[0], key: fields[1], id: fields[2]}, nil
}

type Description struct {
	Name     string `json:"name"`
	Zone     string `json:"zone"`
	Metadata `json:"metadata"`
}

func (d Description) sshKeys() []SSHKey {
	for _, i := range d.Items {
		if i.Key == "ssh-keys" {
			return parseSSHKeys(i.Value)
		}
	}
	return nil
}

type SSHKey struct {
	username string
	key      string
	id       string
}

func (sk SSHKey) String() string {
	return fmt.Sprintf("%s %s %s", sk.username, sk.key, sk.id)
}

func parseSSHKeys(sshKeys string) []SSHKey {
	var keys []SSHKey
	for _, k := range strings.Split(sshKeys, "\n") {
		fields := strings.Fields(k)
		if len(fields) >= 3 {
			keys = append(keys, SSHKey{username: fields[0], key: fields[1], id: fields[2]})
		}
	}
	return keys
}

type Metadata struct {
	Items []Item `json:"items"`
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func execute(c Command) (io.Reader, error) {
	var out bytes.Buffer
	execCmd := exec.Command(c.Name(), c.Args()...)
	fmt.Println("Executing command: ", c.String())
	execCmd.Stdout = &out
	err := execCmd.Run()
	return &out, err
}

func getInstances(cfg Config) ([]instance, error) {
	giCmd := GetInstancesCmd(cfg)
	out, err := execute(giCmd)
	if err != nil {
		return nil, err
	}
	var insts []instance
	err = json.NewDecoder(out).Decode(&insts)
	if err != nil {
		return nil, err
	}
	return insts, nil
}

func getDescription(inst string, cfg Config) (Description, error) {
	out, err := execute(DescribeCmd(inst, cfg))
	if err != nil {
		return Description{}, err
	}
	var desc Description
	err = json.NewDecoder(out).Decode(&desc)
	if err != nil {
		return Description{}, err
	}
	return desc, nil
}
