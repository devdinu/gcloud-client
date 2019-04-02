package action

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
)

type dbStore interface {
	CreateBuckets(names []string) error
	Save(context.Context, <-chan gcloud.Instance, *sync.WaitGroup) error
}

func RefreshInstances(ctx context.Context, s dbStore) Action {
	return func(c gcloud.Client, cfg config.Args) error {
		instCfg := config.GetInstanceCmdArgs()
		if !instCfg.Refresh {
			return nil
		}
		var wg sync.WaitGroup
		wg.Add(1)
		insts := refreshProjects(ctx, c, s)
		go func() {
			if err := s.Save(ctx, insts, &wg); err != nil {
				logger.Errorf("[Refresh] error storing instance: %v", err)
			}
			logger.Infof("[Refresh] stored all instances")
		}()
		logger.Debugf("[Refresh] waiting for all goroutines to complete")
		wg.Wait()
		return nil
	}
}

func listProjects(ctx context.Context, c gcloud.Client) (gcloud.Projects, error) {
	args := config.GetArgs()
	cmdCfg := command.Config{Zone: args.Zone, Limit: args.Limit, Format: args.Format}
	projs, err := c.ListProjects(cmdCfg)
	if err != nil {
		logger.Infof("[Instances] list projects failed with error %v", err)
		return nil, err
	}
	return gcloud.Projects(projs), nil
}

func refreshProjects(ctx context.Context, c gcloud.Client, s dbStore) <-chan gcloud.Instance {
	instances := make(chan gcloud.Instance, 10)
	go func() {
		var lwg sync.WaitGroup
		//TODO: replace the chunk with listprojects
		args := config.GetArgs()
		cmdCfg := command.Config{Zone: args.Zone, Limit: args.Limit, Format: args.Format}
		projs, err := c.ListProjects(cmdCfg)
		if err != nil {
			logger.Errorf("[Instances] list projects failed with error %v", err)
			return
		}
		lwg.Add(len(projs))
		for _, pr := range projs {
			if err := s.CreateBuckets([]string{pr.ProjectID}); err != nil {
				logger.Errorf("[Instances] couldn't create bucket %v", pr)
				continue
			}
			logger.Debugf("[Instances] created bucket %s\n", pr.ProjectID)
			go getInstancesForProject(ctx, instances, c, pr, cmdCfg, &lwg)
		}
		lwg.Wait()
		close(instances)
	}()
	return instances
}

func getInstancesForProject(ctx context.Context, instances chan<- gcloud.Instance, c gcloud.Client, pr gcloud.Project, cmdCfg command.Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	cmdCfg.Project = pr.ProjectID
	insts, err := c.GetInstances(cmdCfg)
	if err != nil {
		return fmt.Errorf("[Instances] list instances failed with error %v", err)
	}
	for _, i := range insts {
		i.Project = pr.ProjectID
		instances <- i
		if err != nil {
			return errors.New("Error storing instance: " + err.Error())
		}
	}
	return nil
}
