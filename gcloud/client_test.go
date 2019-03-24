package gcloud

import (
	"io"
	"os"
	"testing"

	"github.com/devdinu/gcloud-client/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListProjectsSuccess(t *testing.T) {
	exec := new(executorMock)
	c := Client{exec}
	cfg := command.Config{Format: "json"}
	f, err := os.Open("../testdata/projects.json")
	require.NoError(t, err)
	defer f.Close()
	exec.On("Execute", command.ListProjects(cfg)).Return(f, nil)
	expectedProjs := []Project{
		Project{Name: "project-1", ProjectID: "gcloud-client-project", State: "ACTIVE"},
		Project{Name: "project-2", ProjectID: "gcloud-client-project", State: "ACTIVE"},
	}

	projs, err := c.ListProjects(cfg)

	require.NoError(t, err)
	assert.Equal(t, 2, len(projs))
	for i, proj := range projs {
		assert.Equal(t, expectedProjs[i], proj)
	}
}

type executorMock struct{ mock.Mock }

func (m *executorMock) Execute(c command.Command) (io.Reader, error) {
	args := m.Called(c)
	return args.Get(0).(io.Reader), args.Error(1)
}
