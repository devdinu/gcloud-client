package action

import (
	"context"
	"fmt"

	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
)

type finder interface {
	Search(ctx context.Context, projs []string, pattern string) ([]gcloud.Instance, error)
}

func SearchInstancesPrefix(ctx context.Context, f finder) Action {
	return func(c gcloud.Client, args config.Args) error {
		projs, err := listProjects(ctx, c)
		if err != nil {
			return err
		}
		pattern := args.InstanceCmdArgs.Prefix
		//TODO: override with commandline projects to reduce search space
		insts, err := f.Search(ctx, projs.Names(), pattern)
		if err != nil {
			fmt.Printf("[Search] couldn't search instances with prefix %s err: %v", pattern, err)
			return err
		}
		logger.Infof("Search By Prefix Result: ")
		for _, ins := range insts {
			logger.Infof("%s: name: %s\tip: %s\t", ins.Project, ins.Name, ins.IP())
		}
		return nil
	}
}
