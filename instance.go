package main

import (
	"bytes"
	"io/ioutil"
	"os"
)

type instance struct {
	Name string `json:"name"`
	Zone string `json:"zone"`
}

func (in instance) AddSSHKeys(cfg Config, keys []SSHKey) (string, error) {
	f, err := createTempFile(keys)
	if err != nil {
		return "", err
	}
	addCmd := AddSSHKeyCmd(in.Name, f.Name(), cfg)
	rdr, err := execute(addCmd)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	buf := new(bytes.Buffer)
	buf.ReadFrom(rdr)
	return buf.String(), nil
}

func createTempFile(keys []SSHKey) (*os.File, error) {
	f, err := ioutil.TempFile("", "ssh-key")
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		f.WriteString(k.String() + "\n")
	}
	f.Sync()
	return f, nil
}
