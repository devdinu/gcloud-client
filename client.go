package main

import (
	"bytes"
	"encoding/json"
	"os"
)

type client struct {
	executor
}

func (c client) getInstances(cfg Config) ([]instance, error) {
	giCmd := GetInstancesCmd(cfg)
	out, err := c.Execute(giCmd)
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

//TODO: move this to instance
func (c client) getDescription(inst string, cfg Config) (Description, error) {
	out, err := c.Execute(DescribeCmd(inst, cfg))
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

func (c client) AddSSHKeys(inst string, cfg Config, keys []SSHKey) (string, error) {
	f, err := createTempFile(keys)
	if err != nil {
		return "", err
	}
	addCmd := AddSSHKeyCmd(inst, f.Name(), cfg)
	rdr, err := c.Execute(addCmd)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	buf := new(bytes.Buffer)
	buf.ReadFrom(rdr)
	return buf.String(), nil
}
