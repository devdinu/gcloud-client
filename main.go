package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	cfg := Config{format: "json", limit: 2}
	insts, err := getInstances(cfg)
	for _, inst := range insts {
		ssh, err := getSSHMeta(inst, Config{format: "json"})
		fmt.Println("meta", inst.Name, ssh, err)
	}
	fmt.Println(insts, err)
}

type SSHKey struct {
	Name     string `json:"name"`
	Zone     string `json:"zone"`
	Metadata `json:"metadata"`
}

type Metadata struct {
	Items []Item `json:"items"`
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type instance struct {
	Name string `json:"name"`
	Zone string `json:"zone"`
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

func getSSHMeta(inst instance, cfg Config) ([]SSHKey, error) {
	out, err := execute(DescribeCmd(inst.Name, cfg))
	if err != nil {
		return nil, err
	}
	var keys []SSHKey
	err = json.NewDecoder(out).Decode(&keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
