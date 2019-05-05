package store

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"
)

type DB struct {
	*bolt.DB
}

func (db DB) CreateBuckets(names []string) error {
	return db.Update(func(tx *bolt.Tx) error {
		for _, n := range names {
			_, err := tx.CreateBucketIfNotExists([]byte(n))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db DB) Save(ctx context.Context, instances <-chan gcloud.Instance, wg *sync.WaitGroup) error {
	defer wg.Done()

	var batchInsts []gcloud.Instance
	for inst := range instances {
		batchInsts = append(batchInsts, inst) //TODO: make batch size configurable
		if len(batchInsts) >= 100 {
			err := db.updateBatch(batchInsts)
			if err != nil {
				return err
			}
			logger.Debugf("[DB] Written batch of size:%d\n", len(batchInsts))
			batchInsts = []gcloud.Instance{}
		}
	}
	return db.updateBatch(batchInsts)
}

func (db DB) updateBatch(instances []gcloud.Instance) error {
	if len(instances) <= 0 {
		return nil
	}
	for _, inst := range instances {
		err := db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(inst.Project))
			if b == nil {
				return fmt.Errorf("[DB] save instances failed for bucket %s inst: %+v", inst.Project, inst)
			}
			var data bytes.Buffer
			if err := gob.NewEncoder(&data).Encode(inst); err != nil {
				return err
			}
			logger.Debugf("[DB] storing instance %s into bucket: %s", inst.Name, inst.Project)
			return b.Put([]byte(inst.Name), data.Bytes())
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type KeyVals map[string][]byte

// Bucket must exist before write
func (db DB) Write(ctx context.Context, bucket string, data KeyVals) error {
	for k, v := range data {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				return fmt.Errorf("[DB] Write failed since bucket %s not found", bucket)
			}
			return b.Put([]byte(k), v)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (db DB) Search(ctx context.Context, projs []string, matches Predicate) ([]gcloud.Instance, error) {
	var insts []gcloud.Instance

	for _, proj := range projs {
		logger.Debugf("[Search] searching project: %s", proj)
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(proj))
			if b == nil {
				return fmt.Errorf("Bucket not found for project %s", proj)
			}
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if !matches(k) {
					continue
				}
				var data bytes.Buffer
				var found gcloud.Instance
				err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&found)
				if err != nil {
					return err
				}
				logger.Debugf("[Search] found: %s %v", string(k), data.String())
				insts = append(insts, found)
			}
			return nil
		})
		if err != nil {
			logger.Warnf("Searching in bucket: %s err: %v\n", proj, err)
		}
	}

	return insts, nil
}

func NewDB(file string) (DB, error) {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return DB{}, fmt.Errorf("Init db failed %v, for: %s", err, file)
	}
	return DB{db}, nil
}
