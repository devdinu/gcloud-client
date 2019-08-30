package action

import (
	"context"
	"fmt"

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

func (il InstanceLogin) Login(c gcloud.Client, args config.Args) error {
	projs, err := il.lister.Projects(il.ctx, c)
	if err != nil {
		return err
	}
	pattern := args.InstanceCmdArgs.Prefix
	//TODO: override with commandline projects to reduce search space

	matcher, err := getMatcher(args.InstanceCmdArgs)
	if err != nil {
		logger.Errorf("[Login] couldn't create matcher regexp: %s prefix: %s",
			args.InstanceCmdArgs.Regex, args.InstanceCmdArgs.Prefix)
	}
	insts, err := il.f.Search(il.ctx, projs.Names(), matcher)
	if err != nil {
		logger.Errorf("[Login] couldn't search instances with prefix %s err: %v", pattern, err)
		return err
	}
	logger.Infof("Search By Prefix Result: ")
	for _, ins := range insts {
		fmt.Println(ins.String())
	}
	tmuxCfg := command.TmuxConfig{
		TemplatesDir: args.TemplatesDir,
		Project:      "ssh_instances_pane_cmd",
		Session:      args.Login.Session,
	}
	tmuxCfg.AddArg("user", args.Login.User)
	_, err = c.Login(il.ctx, insts, "hostname", tmuxCfg)
	return err
}

func getMatcher(iargs config.InstanceCmdArgs) (store.Predicate, error) {
	if iargs.Regex != "" {
		return store.RegexMatcher(iargs.Regex)
	}
	return store.PrefixMatcher(iargs.Prefix), nil
}
