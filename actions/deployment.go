package actions

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/CenturyLinkLabs/prettycli"
	"github.com/ghodss/yaml"
)

func ListDeployments(remote config.Remote) (prettycli.Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	deps, err := c.ListDeployments()
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	if len(deps) == 0 {
		return prettycli.PlainOutput{"No Deployments"}, nil
	}

	o := prettycli.ListOutput{Labels: []string{"ID", "Name", "Services"}}
	for _, d := range deps {
		o.AddRow(map[string]string{
			"ID":       strconv.Itoa(d.ID),
			"Name":     d.Name,
			"Services": strconv.Itoa(len(d.ServiceIDs)),
		})
	}

	return &o, nil
}

func DescribeDeployment(remote config.Remote, id string) (prettycli.Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	desc, err := c.DescribeDeployment(id)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	do := prettycli.DetailOutput{
		Details: map[string]string{
			"Name":         desc.Name,
			"ID":           strconv.Itoa(desc.ID),
			"Redeployable": strconv.FormatBool(desc.Redeployable),
		},
		Order: []string{"ID", "Name", "Redeployable"},
	}

	lo := prettycli.ListOutput{Labels: []string{"ID", "State"}}
	for _, s := range desc.Status.Services {
		lo.AddRow(map[string]string{
			"ID":    s.ID,
			"State": s.ActualState,
		})
	}

	co := prettycli.CombinedOutput{}
	co.AddOutput("", do)
	co.AddOutput("Services", lo)
	return &co, nil
}

func CreateDeployment(remote config.Remote, path string) (prettycli.Output, error) {
	templateBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	bp := agent.DeploymentBlueprint{}
	if err := yaml.Unmarshal(templateBytes, &bp.Template); err != nil {
		return prettycli.PlainOutput{}, err
	}

	dr, err := DefaultAgentClientFactory.New(remote).CreateDeployment(bp)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	return prettycli.PlainOutput{fmt.Sprintf("Template successfully deployed as '%d'", dr.ID)}, nil
}

func RedeployDeployment(remote config.Remote, id string) (prettycli.Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	desc, err := c.RedeployDeployment(id)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	o := prettycli.PlainOutput{fmt.Sprintf("Redeployed '%s' as Deployment ID %d", desc.Name, desc.ID)}
	return &o, nil
}

func DeleteDeployment(remote config.Remote, id string) (prettycli.Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	if err := c.DeleteDeployment(id); err != nil {
		return prettycli.PlainOutput{}, err
	}

	o := prettycli.PlainOutput{fmt.Sprintf("Successfully deleted deployment '%s'", id)}
	return &o, nil
}
