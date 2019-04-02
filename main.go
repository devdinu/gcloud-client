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
	ctx := context.Background()
	db, err := store.NewDB(config.GetDBFileName())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	action.MapActions(ctx, db)

	action, err := action.GetAction(config.GetActionName())
	if err != nil {
		log.Fatal(err)
	}
	if err := action(c, args); err != nil {
		log.Fatal(err)
	}
}
