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
	var sshFile string
	flag.StringVar(&sshFile, "ssh_key", "", "new SSH Key file which have to be added to instances")
	flag.Parse()

	if sshFile == "" {
		panic("SSH File is mandatory")
	}
	cfg := Config{format: "json", limit: 1}
	insts, _ := getInstances(cfg)
	for _, inst := range insts {
		desc, err := getDescription(inst.Name, Config{format: "json", zone: inst.Zone})
		fmt.Println("meta", inst.Name, desc.sshKeys(), err)
		keys := desc.sshKeys()
		newKey, err := readKey(sshFile)
		keys = append(keys, newKey)
		inst.AddSSHKeys(cfg, keys)
	}
}

func readKey(filename string) (SSHKey, error) {
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
	user := os.Getenv("USER")
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
		fmt.Println(len(fields))
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
	fmt.Println(giCmd)
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
