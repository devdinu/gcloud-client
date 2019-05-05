package gcloud

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Instance struct {
	Name              string             `json:"name"`
	Zone              string             `json:"zone"`
	NetworkInterfaces []NetworkInterface `json:"networkInterfaces"`
	Status            string             `json:"status"`
	Project           string             `json:"projectID"`
}

func (ins Instance) String() string {
	res := fmt.Sprintf("%-30s : %-50s %-10s", ins.Project, ins.Name, ins.IP())
	if ins.ExternalIP() != "" {
		return fmt.Sprintf("%s External: %s", res, ins.ExternalIP())
	}
	return res
}

func (i Instance) IP() string {
	if len(i.NetworkInterfaces) == 0 {
		return ""
	}
	return i.NetworkInterfaces[0].NetworkIP
}

func (i Instance) ExternalIP() string {
	if len(i.NetworkInterfaces) > 0 &&
		len(i.NetworkInterfaces[0].AccessConfigs) > 0 {
		return i.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}
	return ""
}

type NetworkInterface struct {
	NetworkIP     string `json:"networkIP"`
	AccessConfigs []AccessConfig
}

type AccessConfig struct {
	NatIP string `natIP`
	Name  string `name`
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
