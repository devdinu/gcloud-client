package gcloud

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/devdinu/gcloud-client/command"
)

type executor interface {
	Execute(c command.Command) (io.Reader, error)
}

type Client struct {
	executor
}

func (c Client) GetInstances(cfg command.Config) ([]Instance, error) {
	giCmd := command.GetInstancesCmd(cfg)
	out, err := c.Execute(giCmd)
	if err != nil {
		return nil, err
	}
	var insts []Instance
	err = json.NewDecoder(out).Decode(&insts)
	if err != nil {
		return nil, err
	}
	return insts, nil
}

//TODO: move this to instance
func (c Client) GetDescription(inst string, cfg command.Config) (Description, error) {
	out, err := c.Execute(command.DescribeCmd(inst, cfg))
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

func (c Client) AddSSHKeys(inst string, cfg command.Config, keys []SSHKey) (string, error) {
	f, err := createTempFile(keys)
	if err != nil {
		return "", err
	}
	addCmd := command.AddSSHKeyCmd(inst, f.Name(), cfg)
	rdr, err := c.Execute(addCmd)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	buf := new(bytes.Buffer)
	buf.ReadFrom(rdr)
	return buf.String(), nil
}
func NewClient(e executor) Client {
	return Client{e}
}
