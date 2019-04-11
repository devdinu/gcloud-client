package main

import (
	"context"
	"log"

	action "github.com/devdinu/gcloud-client/actions"
	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
	"github.com/devdinu/gcloud-client/store"
)

func main() {
	config.MustLoad()
	logger.SetLevel(config.LogLevel())
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
	logger.Debugf("[gcloud-client] Executing action: %s", config.GetActionName())
	if err := action(c, args); err != nil {
		log.Fatal(err)
	}
}
