package action

import (
	"context"
	"fmt"

	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/store"
)

type Action func(gcloud.Client, config.Args) error

var actions map[config.CmdAction]Action

func MapActions(ctx context.Context, db store.DB) {
	actions = make(map[config.CmdAction]Action)

	lister := lister{store: db}
	login := InstanceLogin{ctx: ctx, f: db, lister: lister}
	srch := searcher{ctx: ctx, lister: lister, finder: db}

	actions[config.SSHAccess] = AddSSHKeys
	actions[config.RefreshInstances] = Refresher{ctx: ctx, store: db}.RefreshInstances
	actions[config.SearchPrefix] = srch.SearchInstancesPrefix
	actions[config.SearchRegex] = srch.SearchInstancesRegex
	actions[config.LoginInstances] = login.Login

}

func GetAction(ca config.CmdAction) (Action, error) {
	val, ok := actions[ca]
	if !ok {
		return nil, fmt.Errorf("Action not found for cmd: %s", ca)
	}
	return val, nil
}
