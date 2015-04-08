package actions

import (
	"fmt"
	"strconv"

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

	o := ListOutput{Labels: []string{"ID", "Name", "Services"}}
	for _, d := range deps {
		o.AddRow(map[string]string{
			"ID":       strconv.Itoa(d.ID),
			"Name":     d.Name,
			"Services": strconv.Itoa(len(d.ServiceIDs)),
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

	do := DetailOutput{
		Details: map[string]string{
			"Name":         desc.Name,
			"ID":           strconv.Itoa(desc.ID),
			"Redeployable": strconv.FormatBool(desc.Redeployable),
		},
	}

	lo := ListOutput{Labels: []string{"ID", "State"}}
	for _, s := range desc.Status.Services {
		lo.AddRow(map[string]string{
			"ID":    s.ID,
			"State": s.ActualState,
		})
	}

	co := CombinedOutput{}
	co.AddOutput("", do)
	co.AddOutput("Services", lo)
	return &co, nil
}

func RedeployDeployment(remote config.Remote, id string) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	desc, err := c.RedeployDeployment(id)
	if err != nil {
		return PlainOutput{}, err
	}

	o := PlainOutput{fmt.Sprintf("Redeployed '%s' as Deployment ID %d", desc.Name, desc.ID)}
	return &o, nil
}

func DeleteDeployment(remote config.Remote, id string) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	if err := c.DeleteDeployment(id); err != nil {
		return PlainOutput{}, err
	}

	o := PlainOutput{fmt.Sprintf("Successfully deleted deployment '%s'", id)}
	return &o, nil
}
