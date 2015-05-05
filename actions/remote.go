package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/CenturyLinkLabs/prettycli"
)

var format = regexp.MustCompile("^[a-zA-Z0-9]+$")

func AddRemote(config config.Config, name string, token []byte) (prettycli.Output, error) {
	if !format.MatchString(name) {
		return prettycli.PlainOutput{}, errors.New("Invalid name")
	}
	if _, err := config.Get(name); err == nil {
		return prettycli.PlainOutput{}, errors.New("Name already exists")
	}
	trimmedToken := strings.TrimSpace(string(token))
	if err := config.Save(name, trimmedToken); err != nil {
		return prettycli.PlainOutput{}, err
	}

	if len(config.Remotes()) == 1 {
		config.SetActive(name)
	}
	s := "Successfully added!"
	if config.Active() != nil {
		s += fmt.Sprintf(" '%s' is your active remote.", config.Active().Name)
	}
	return prettycli.PlainOutput{s}, nil
}

func AddRemoteByPath(config config.Config, name string, path string) (prettycli.Output, error) {
	token, err := ioutil.ReadFile(path)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	return AddRemote(config, name, token)
}

func RemoveRemote(config config.Config, name string) (prettycli.Output, error) {
	if err := config.Remove(name); err != nil {
		return prettycli.PlainOutput{}, err
	}
	out := fmt.Sprintf("Successfully removed remote '%s' from configuration!", name)
	return prettycli.PlainOutput{out}, nil
}

func ListRemotes(config config.Config) prettycli.Output {
	agents := config.Remotes()
	if len(agents) == 0 {
		return prettycli.PlainOutput{"No remotes"}
	}

	output := prettycli.ListOutput{Labels: []string{"Active", "Name", "Endpoint"}}
	for _, r := range config.Remotes() {
		activeMarker := ""
		if config.Active() != nil && *config.Active() == r {
			activeMarker = "*"
		}

		output.AddRow(map[string]string{
			"Active":   activeMarker,
			"Name":     r.Name,
			"Endpoint": r.Endpoint,
		})
	}
	return &output
}

func DescribeRemote(c config.Config, name string) (prettycli.Output, error) {
	r, err := c.Get(name)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	isActive := "false"
	if c.Active() != nil && c.Active().Name == r.Name {
		isActive = "true"
	}

	client := DefaultAgentClientFactory.New(r)
	metadata, err := client.GetMetadata()
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	adapterMetadataBytes, err := json.Marshal(metadata.Adapter)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	adapterMetadata := struct {
		Version   string
		Type      string
		IsHealthy bool
	}{}
	if err := json.Unmarshal(adapterMetadataBytes, &adapterMetadata); err != nil {
		return prettycli.PlainOutput{}, err
	}

	do := prettycli.DetailOutput{
		Details: map[string]string{
			"Name":               r.Name,
			"Active":             isActive,
			"Endpoint":           r.Endpoint,
			"Agent Version":      metadata.Agent.Version,
			"Adapter Version":    adapterMetadata.Version,
			"Adapter Type":       adapterMetadata.Type,
			"Adapter Is Healthy": strconv.FormatBool(adapterMetadata.IsHealthy),
		},
		Order: []string{"Name", "Active", "Endpoint"},
	}

	lo, err := ListDeployments(r)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}

	co := prettycli.CombinedOutput{}
	co.AddOutput("", do)
	co.AddOutput("Deployments", lo)

	return &co, nil
}

func SetActiveRemote(config config.Config, name string) (prettycli.Output, error) {
	if err := config.SetActive(name); err != nil {
		return prettycli.PlainOutput{}, err
	}
	return prettycli.PlainOutput{fmt.Sprintf("'%s' is now your active remote!", name)}, nil
}

func GetRemoteToken(c config.Config, name string) (prettycli.Output, error) {
	r, err := c.Get(name)
	if err != nil {
		return prettycli.PlainOutput{}, err
	}
	return prettycli.PlainOutput{r.Token}, nil
}
