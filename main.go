package main

import (
	"context"
	"flag"
	"log"
	"os"

	action "github.com/devdinu/gcloud-client/actions"
	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/store"
)

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		return
	}
	config.Load()
	args := config.GetArgs()

	c := gcloud.NewClient(command.Executor{})
	var cmdAction action.Action
	if os.Args[1] == "ssh_access" || os.Args[1] == "" {
		cmdAction = action.AddSSHKeys
	} else if os.Args[1] == "instances" {
		cmdAction = action.RefreshInstances(context.Background(), store.NewDB())
	} else {
		flag.Usage()
		return
	}
	if err := cmdAction(c, args); err != nil {
		log.Fatal(err)
	}
}
