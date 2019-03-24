package gcloud

import (
	"fmt"
	"strings"
)

type Description struct {
	Name     string `json:"name"`
	Zone     string `json:"zone"`
	Metadata `json:"metadata"`
}

func (d Description) SshKeys() []SSHKey {
	for _, i := range d.Items {
		if i.Key == "ssh-keys" {
			return parseSSHKeys(i.Value)
		}
	}
	return nil
}

type SSHKey struct {
	Username string
	Key      string
	ID       string
}

func (sk SSHKey) String() string {
	return fmt.Sprintf("%s %s %s", sk.Username, sk.Key, sk.ID)
}

func parseSSHKeys(sshKeys string) []SSHKey {
	var keys []SSHKey
	for _, k := range strings.Split(sshKeys, "\n") {
		fields := strings.Fields(k)
		if len(fields) >= 3 {
			keys = append(keys, SSHKey{Username: fields[0], Key: fields[1], ID: fields[2]})
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
