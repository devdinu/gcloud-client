package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type instance struct {
	Name string `json:"name"`
	Zone string `json:"zone"`
}

func (in instance) AddSSHKeys(cfg Config, keys []SSHKey) error {
	fmt.Println("adding keys", keys)
	f, err := createTempFile(keys)
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	fmt.Println("Wrote ssh-keys to file: ", f.Name())
	return nil
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
