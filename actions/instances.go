package action

import (
	"fmt"
	"sync"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
)

func RefreshInstances(c gcloud.Client, cfg config.Args) error {
	instCfg := config.GetInstanceCmdArgs()
	if !instCfg.Refresh {
		return nil
	}
	args := config.GetArgs()
	cmdCfg := command.Config{Zone: args.Zone, Limit: args.Limit, Format: args.Format}
	projs, err := c.ListProjects(cmdCfg)
	if err != nil {
		return fmt.Errorf("[Instances] list projects failed with error %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(len(projs))
	for _, pr := range projs {
		go getInstancesForProject(c, pr, cmdCfg, &wg)
	}
	wg.Wait()
	return nil
}

func getInstancesForProject(c gcloud.Client, pr gcloud.Project, cmdCfg command.Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	cmdCfg.Project = pr.ProjectID
	insts, err := c.GetInstances(cmdCfg)
	if err != nil {
		return fmt.Errorf("[Instances] list instances failed with error %v", err)
	}
	for _, i := range insts {
		fmt.Println(i.Name)
	}
	return nil
}
