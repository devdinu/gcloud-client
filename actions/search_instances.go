package action

import (
	"context"
	"fmt"

	"github.com/devdinu/gcloud-client/config"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
	"github.com/devdinu/gcloud-client/store"
)

type finder interface {
	Search(ctx context.Context, projs []string, matcher store.Predicate) ([]gcloud.Instance, error)
}

type searcher struct {
	lister
	finder
	ctx context.Context
}

func (s searcher) SearchInstancesPrefix(c gcloud.Client, args config.Args) error {
	pattern := args.InstanceCmdArgs.Prefix
	insts, err := s.searchInstances(c, args, store.PrefixMatcher(pattern))
	if err != nil {
		return fmt.Errorf("[SearchPrefix] couldn't search instances with prefix %s err: %v", pattern, err)
	}
	logger.Infof("Search By Prefix Result: ")
	s.formatInstances(insts, args)
	return nil
}

func (s searcher) SearchInstancesRegex(c gcloud.Client, args config.Args) error {
	regex := args.InstanceCmdArgs.Regex
	matcher, err := store.RegexMatcher(regex)
	if err != nil {
		return fmt.Errorf("[SearchRegex] couldn't create regex matcher with pattern :%s, error :%v", regex, err)
	}
	insts, err := s.searchInstances(c, args, matcher)
	if err != nil {
		return fmt.Errorf("[SearchRegex] couldn't search instances with regex %s err: %v", regex, err)
	}
	logger.Infof("Search By Regex Result: ")
	s.formatInstances(insts, args)
	return nil
}

func (s searcher) searchInstances(c gcloud.Client, args config.Args, matcher store.Predicate) ([]gcloud.Instance, error) {
	projs, err := s.lister.Projects(s.ctx, c)
	if err != nil {
		return nil, err
	}
	return s.finder.Search(s.ctx, projs.Names(), matcher)
}

func (s searcher) formatInstances(insts []gcloud.Instance, args config.Args) error {
	hostMapping := args.InstanceCmdArgs.HostMapping
	if hostMapping {
		for _, ins := range insts {
			fmt.Printf("%-16s  %s\n", ins.IP(), ins.Name)
		}
		return nil
	}
	for _, ins := range insts {
		fmt.Printf("%s\n", ins)
	}
	return nil
}
