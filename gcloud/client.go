package gcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/logger"
)

type executor interface {
	Execute(c command.Command) (io.Reader, error)
}

type Client struct {
	executor
}

//TODO: get project name as arg
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

func (c Client) ListProjects(cfg command.Config) ([]Project, error) {
	projs, err := c.Execute(command.ListProjects(cfg))
	if err != nil {
		return nil, err
	}
	var projects []Project
	err = json.NewDecoder(projs).Decode(&projects)
	if err != nil {
		return nil, err
	}
	return projects, err
}

//TODO: move to separate as it doesn't deal with gcloud
func (c Client) Login(ctx context.Context, insts []Instance, cmd string, cfg command.TmuxConfig) (string, error) {
	var hosts []string
	for _, inst := range insts {
		if inst.Status != "RUNNING" {
			logger.Debugf("host: %s not running, state: %s", inst.Name, inst.Status)
		}
		hosts = append(hosts, inst.IP())
	}
	output, err := c.Execute(command.Login(hosts, cmd, cfg))
	res, _ := ioutil.ReadAll(output)
	return string(res), err
}

func NewClient(e executor) Client {
	return Client{e}
}
