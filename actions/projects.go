package action

import (
	"context"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
)

type lister struct {
	store dbStore
}

func (l lister) Projects(ctx context.Context, c gcloud.Client) (gcloud.Projects, error) {
	args := config.GetArgs()
	cmdCfg := command.Config{Zone: args.Zone, Format: args.Format}
	customProjects := config.Projects()
	allProjs, err := c.ListProjects(cmdCfg)
	if err != nil {
		logger.Infof("[Instances] list projects failed with error %v", err)
		return nil, err
	}
	var projs gcloud.Projects
	if len(customProjects) > 0 {
		logger.Debugf("[Lister] using custom configuration: %s", customProjects)
		for _, p := range allProjs {
			if contains(customProjects, p.ProjectID) {
				projs = append(projs, p)
			}
		}
	} else {
		projs = gcloud.Projects(allProjs)
	}
	logger.Debugf("[Lister] final projects list: %v", projs)
	return projs, nil
}

func contains(list []string, elem string) bool {
	for _, x := range list {
		if x == elem {
			return true
		}
	}
	return false
}
