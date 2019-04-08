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

type Refresher struct {
	ctx   context.Context
	store dbStore
	lister
}

func (r Refresher) RefreshInstances(c gcloud.Client, cfg config.Args) error {
	instCfg := config.GetInstanceCmdArgs()
	if !instCfg.Refresh {
		return nil
	}
	err := r.refreshProjects(r.ctx, c)
	logger.Debugf("[Refresh] waiting for all goroutines to complete")
	return err
}

func (r Refresher) refreshProjects(ctx context.Context, c gcloud.Client) error {
	var lwg sync.WaitGroup
	//TODO: replace the chunk with listprojects
	args := config.GetArgs()
	cmdCfg := command.Config{Zone: args.Zone, Limit: args.Limit, Format: args.Format}
	projs, err := r.lister.Projects(ctx, c)
	if err != nil {
		return err
	}
	lwg.Add(len(projs))
	for _, pr := range projs {
		if err := r.store.CreateBuckets([]string{pr.ProjectID}); err != nil {
			logger.Errorf("[Instances] couldn't create bucket %v", pr)
			continue
		}
		logger.Debugf("[Instances] created bucket %s\n", pr.ProjectID)
		go r.getInstancesForProject(ctx, c, pr, cmdCfg, &lwg)
	}
	lwg.Wait()
	return nil
}

func (r Refresher) getInstancesForProject(ctx context.Context, c gcloud.Client, pr gcloud.Project, cmdCfg command.Config, wg *sync.WaitGroup) error {
	defer wg.Done()

	var lwg sync.WaitGroup
	instances := make(chan gcloud.Instance, 10000)
	cmdCfg.Project = pr.ProjectID
	insts, err := c.GetInstances(cmdCfg)
	if err != nil {
		return fmt.Errorf("[Instances] list instances failed with error %v", err)
	}
	totalStoreConcurrency := 10
	lwg.Add(totalStoreConcurrency)
	for i := 0; i < totalStoreConcurrency; i++ {
		go func(id int) {
			if err := r.store.Save(ctx, instances, &lwg); err != nil {
				logger.Errorf("[Refresh] error storing instance: %v", err)
			}
			logger.Infof("[Refresh] goroutine :%d completed", id)
		}(i)
	}
	for _, i := range insts {
		i.Project = pr.ProjectID
		instances <- i
		if err != nil {
			return errors.New("Error storing instance: " + err.Error())
		}
	}
	logger.Debugf("[Refresher] Got instances size: %d, writing with concurrency: %d chan:%d", len(insts), totalStoreConcurrency, len(instances))
	close(instances)
	lwg.Wait()
	return nil
}
