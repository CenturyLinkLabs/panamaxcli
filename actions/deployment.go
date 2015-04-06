package actions

import (
	"strconv"
	"strings"

	"github.com/CenturyLinkLabs/panamaxcli/config"
)

func ListDeployments(remote config.Remote) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	deps, err := c.ListDeployments()
	if err != nil {
		return PlainOutput{}, err
	}

	if len(deps) == 0 {
		return PlainOutput{"No Deployments"}, nil
	}

	o := ListOutput{Labels: []string{"ID", "Name"}}
	for _, d := range deps {
		o.AddRow(map[string]string{
			"ID":   strconv.Itoa(d.ID),
			"Name": d.Name,
		})
	}

	return &o, nil
}

func DescribeDeployment(remote config.Remote, id string) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	desc, err := c.DescribeDeployment(id)
	if err != nil {
		return PlainOutput{}, err
	}

	statuses := make([]string, len(desc.Status.Services))
	for i, s := range desc.Status.Services {
		statuses[i] = s.ActualState
	}

	o := DetailOutput{
		Details: map[string]string{
			"Name":             desc.Name,
			"ID":               strconv.Itoa(desc.ID),
			"Redeployable":     strconv.FormatBool(desc.Redeployable),
			"Service Statuses": strings.Join(statuses, ", "),
		},
	}

	return &o, nil
}
