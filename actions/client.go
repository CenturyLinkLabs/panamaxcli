package actions

import (
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/client"
	"github.com/CenturyLinkLabs/panamaxcli/config"
)

type AgentClientFactory interface {
	New(config.Remote) client.Client
}

var DefaultAgentClientFactory AgentClientFactory

func init() {
	DefaultAgentClientFactory = &APIClientFactory{}
}

type APIClientFactory struct{}

func (f *APIClientFactory) New(r config.Remote) client.Client {
	return &client.APIClient{
		Endpoint:   r.Endpoint,
		Username:   r.Username,
		Password:   r.Password,
		PrivateKey: r.PrivateKey,
	}
}
