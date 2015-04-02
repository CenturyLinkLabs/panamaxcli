package actions

import (
	"strconv"

	"github.com/CenturyLinkLabs/panamaxcli/config"
)

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
