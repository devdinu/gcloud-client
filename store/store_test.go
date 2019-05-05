package store_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/devdinu/gcloud-client/logger"

	"github.com/devdinu/gcloud-client/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSearchInstances(t *testing.T) {
	suite.Run(t, new(SearchSuite))
}

type SearchSuite struct {
	suite.Suite
	db        store.DB
	bucket    string
	instances []gcloud.Instance
	ctx       context.Context
}

func (s *SearchSuite) TestSetup() {
	logger.SetLevel("error")
	s.ctx, _ = context.WithTimeout(context.Background(), time.Second*3)
	t := s.T()
	var err error
	s.db, err = store.NewDB("./testdata/search.db")
	s.bucket = "search-instances"
	s.instances = []gcloud.Instance{
		{Name: "integration-01"},
		{Name: "integration-02"},
		{Name: "prod-01-vm"},
		{Name: "prod-02-vm"},
		{Name: "prod-a-01"},
		{Name: "integration-03"},
		{Name: "prod-db-03"},
	}
	require.NoError(t, err)

	require.NoError(t, s.db.CreateBuckets([]string{s.bucket}))
	for _, inst := range s.instances {
		err := s.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(s.bucket))
			require.NotNil(t, b)
			var data bytes.Buffer
			require.NoError(t, gob.NewEncoder(&data).Encode(inst))
			return b.Put([]byte(inst.Name), data.Bytes())
		})
		require.NoError(t, err)
	}
}

func (s *SearchSuite) TestShouldMatchIntancesWithPrefix() {
	t := s.T()
	projs := []string{s.bucket}

	insts, err := s.db.Search(s.ctx, projs, store.PrefixMatcher("integration"))

	require.NoError(t, err)
	require.Equal(t, 3, len(insts))
	assert.Equal(t, insts[0].Name, "integration-01")
	assert.Equal(t, insts[1].Name, "integration-02")
	assert.Equal(t, insts[2].Name, "integration-03")
}

func (s *SearchSuite) TestShouldMatchIntancesWithRegex() {
	t := s.T()
	projs := []string{s.bucket}
	rmatch, err := store.RegexMatcher(`prod-\d+-vm`)
	require.NoError(t, err)

	insts, err := s.db.Search(s.ctx, projs, rmatch)

	require.NoError(t, err)
	require.Equal(t, 2, len(insts))
	assert.Equal(t, insts[0].Name, "prod-01-vm")
	assert.Equal(t, insts[1].Name, "prod-02-vm")
}
