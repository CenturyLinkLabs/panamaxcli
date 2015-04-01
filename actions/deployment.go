package actions

import (
	"strconv"

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

func ListDeployments(remote config.Remote) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	deps, _ := c.ListDeployments()

	o := ListOutput{Labels: []string{"ID", "Name"}}
	for _, d := range deps {
		o.AddRow(map[string]string{
			"ID":   strconv.Itoa(d.ID),
			"Name": d.Name,
		})
	}

	return &o, nil
}
