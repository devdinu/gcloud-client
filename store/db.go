package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/devdinu/gcloud-client/gcloud"
)

type DB struct{}

func (db DB) Save(ctx context.Context, instances <-chan gcloud.Instance, wg *sync.WaitGroup) error {
	defer wg.Done()

	for inst := range instances {
		fmt.Printf("storing instance %v \n", inst)
	}
	fmt.Println("stored all instances")
	return nil
}

func NewDB() DB {
	return DB{}
}
