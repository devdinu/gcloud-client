package action

import (
	"context"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
	"github.com/devdinu/gcloud-client/store"
)

type InstanceLogin struct {
	ctx context.Context
	f   finder
	lister
}

//TODO: login could be done as search instances and login, otherwise display
func (il InstanceLogin) Login(c gcloud.Client, args config.Args) error {
	projs, err := il.lister.Projects(il.ctx, c)
	if err != nil {
		return err
	}
	pattern := args.InstanceCmdArgs.Prefix
	//TODO: override with commandline projects to reduce search space
	// use regex match if available
	insts, err := il.f.Search(il.ctx, projs.Names(), store.PrefixMatcher(pattern))
	if err != nil {
		logger.Errorf("[Search] couldn't search instances with prefix %s err: %v", pattern, err)
		return err
	}
	logger.Infof("Search By Prefix Result: ")
	for _, ins := range insts {
		logger.Infof("%s: name: %s\tip: %s external: %s\t", ins.Project, ins.Name, ins.IP())
	}
	tmuxCfg := command.TmuxConfig{
		Project: "ssh_instances_pane_cmd",
		Session: args.Login.Session,
	}
	tmuxCfg.AddArg("user", args.Login.User)
	_, err = c.Login(il.ctx, insts, "hostname", tmuxCfg)
	return err
}
