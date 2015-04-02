package actions

import (
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/client"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/stretchr/testify/assert"
)

type FakeFactory struct {
	NewedRemotes []config.Remote
}

func (f *FakeFactory) New(r config.Remote) client.Client {
	f.NewedRemotes = append(f.NewedRemotes, r)
	return &fakeClient
}

type FakeClient struct {
	Deployments      []agent.DeploymentResponseLite
	Metadata         agent.Metadata
	ErrorForMetadata error
}

func (c FakeClient) ListDeployments() ([]agent.DeploymentResponseLite, error) {
	return c.Deployments, nil
}

func (c FakeClient) GetMetadata() (agent.Metadata, error) {
	return c.Metadata, c.ErrorForMetadata
}

var (
	fakeFactory = FakeFactory{}
	fakeClient  = FakeClient{}
)

func init() {
	DefaultAgentClientFactory = &fakeFactory
}

func setupFactory() {
	fakeFactory = FakeFactory{}
	fakeClient = FakeClient{}
}

func TestAPIClientFactoryNew(t *testing.T) {
	r := config.Remote{Endpoint: "http://example.com"}
	f := APIClientFactory{}
	c := f.New(r)
	ac, ok := c.(*client.APIClient)
	if assert.True(t, ok) {
		assert.Equal(t, "http://example.com", ac.Endpoint)
	}
}
