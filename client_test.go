package main_test

import (
	"io"
	"os"
	"testing"

	"github.com/devdinu/gcloud-client/command"
	"github.com/devdinu/gcloud-client/gcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	cfg := command.Config{Format: "json"}
	e := new(executorMock)
	client := gcloud.NewClient(e)
	f, err := os.Open("./testdata/reference.json")
	require.NoError(t, err)

	defer f.Close()
	e.On("Execute", command.GetInstancesCmd(cfg)).Return(f, nil)

	insts, err := client.GetInstances(cfg)

	require.NoError(t, err)
	expectedInstance := gcloud.Instance{
		Name: "some-instance-name",
		Zone: "https://www.googleapis.com/compute/v1/projects/some-cluster/zones/somezone-a",
		NetworkInterfaces: []gcloud.NetworkInterface{
			{
				NetworkIP: "10.11.12.13",
				AccessConfigs: []gcloud.AccessConfig{
					{NatIP: "12.34.56.78", Name: "external-nat"},
				},
			},
		},
		Status: "RUNNING",
	}
	e.AssertExpectations(t)
	assert.Equal(t, expectedInstance, insts[0])
}

type executorMock struct{ mock.Mock }

func (m *executorMock) Execute(c command.Command) (io.Reader, error) {
	args := m.Called(c)
	return args.Get(0).(io.Reader), args.Error(1)
}
