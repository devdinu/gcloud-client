package action

import (
	"context"
	"fmt"
	"sync"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
)

type store interface {
	Save(context.Context, <-chan gcloud.Instance, *sync.WaitGroup) error
}

func RefreshInstances(ctx context.Context, s store) Action {
	return func(c gcloud.Client, cfg config.Args) error {
		instCfg := config.GetInstanceCmdArgs()
		if !instCfg.Refresh {
			return nil
		}
		var wg sync.WaitGroup
		wg.Add(1)
		insts := refreshProjects(ctx, c)
		go s.Save(ctx, insts, &wg)
		fmt.Println("waiting for all goroutines to complete")
		wg.Wait()
		return nil
	}
}

func refreshProjects(ctx context.Context, c gcloud.Client) <-chan gcloud.Instance {
	instances := make(chan gcloud.Instance, 10)
	go func() {
		var lwg sync.WaitGroup
		args := config.GetArgs()
		cmdCfg := command.Config{Zone: args.Zone, Limit: args.Limit, Format: args.Format}
		projs, err := c.ListProjects(cmdCfg)
		if err != nil {
			fmt.Errorf("[Instances] list projects failed with error %v", err)
			return
		}
		lwg.Add(len(projs))
		for _, pr := range projs {
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
		instances <- i
		if err != nil {
			return fmt.Errorf("Error storing instance: %v", err)
		}
	}
	return nil
}
