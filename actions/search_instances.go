package action

import (
	"context"
	"fmt"

	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
)

type finder interface {
	Search(ctx context.Context, pattern string) ([]gcloud.Instance, error)
}

func SearchInstancesPrefix(ctx context.Context, f finder) Action {
	return func(c gcloud.Client, args config.Args) error {
		pattern := args.InstanceCmdArgs.Prefix
		fmt.Printf("searching instances with prefix %s\n", pattern)
		insts, err := f.Search(ctx, pattern)
		if err != nil {
			fmt.Printf("[Search] couldn't search instances with prefix %s err: %v", pattern, err)
		}
		return nil
		fmt.Println("Search By Prefix Result: ")
		for _, ins := range insts {
			fmt.Println("instance %s", ins)
		}
		return nil
	}
}
