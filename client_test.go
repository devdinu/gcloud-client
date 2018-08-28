package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	cfg := Config{}
	e := new(executorMock)
	c := client{e}
	f, err := os.Open("./testdata/reference.json")
	require.NoError(t, err)

	defer f.Close()
	e.On("Execute", GetInstancesCmd(cfg)).Return(f, nil)

	insts, err := c.getInstances(cfg)

	require.NoError(t, err)
	expectedInstance := instance{
		Name:              "some-instance-name",
		Zone:              "https://www.googleapis.com/compute/v1/projects/some-cluster/zones/somezone-a",
		NetworkInterfaces: []NetworkInterface{{NetworkIP: "10.11.12.13"}},
		Status:            "RUNNING",
	}
	assert.Equal(t, expectedInstance, insts[0])
}

type executorMock struct{ mock.Mock }

func (m *executorMock) Execute(c Command) (io.Reader, error) {
	args := m.Called(c)
	return args.Get(0).(io.Reader), args.Error(1)
}
