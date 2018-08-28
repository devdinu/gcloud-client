package main

import (
	"io/ioutil"
	"os"
)

type instance struct {
	Name              string             `json:"name"`
	Zone              string             `json:"zone"`
	networkInterfaces []NetworkInterface `json:"networkInterfaces"`
}

func (i instance) IP() string {
	if len(i.networkInterfaces) == 0 {
		return ""
	}
	return i.networkInterfaces[0].NetworkIP
}

type NetworkInterface struct {
	NetworkIP string `json:"networkIP"`
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
