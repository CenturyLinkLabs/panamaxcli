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
)

var format = regexp.MustCompile("^[a-zA-Z0-9]+$")

func AddRemote(config config.Config, name string, path string) (Output, error) {
	if !format.MatchString(name) {
		return PlainOutput{}, errors.New("Invalid name")
	}
	if config.Exists(name) {
		return PlainOutput{}, errors.New("Name already exists")
	}
	token, err := ioutil.ReadFile(path)
	if err != nil {
		return PlainOutput{}, err
	}
	trimmedToken := strings.TrimSpace(string(token))
	if err = config.Save(name, trimmedToken); err != nil {
		return PlainOutput{}, err
	}

	if len(config.Remotes()) == 1 {
		config.SetActive(name)
		s := fmt.Sprintf("Successfully added! '%s' is now your active remote.", name)
		return PlainOutput{s}, nil
	}
	return PlainOutput{"Successfully added!"}, nil
}

func ListRemotes(config config.Config) Output {
	agents := config.Remotes()
	if len(agents) == 0 {
		return PlainOutput{"No remotes"}
	}

	output := ListOutput{Labels: []string{"Active", "Name", "Endpoint"}}
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

func DescribeRemote(c config.Config, name string) (Output, error) {
	var remote config.Remote
	for _, r := range c.Remotes() {
		if r.Name == name {
			remote = r
			break
		}
	}
	if remote.Name == "" {
		return PlainOutput{}, fmt.Errorf("the remote '%s' does not exist", name)
	}

	isActive := "false"
	if c.Active() != nil && c.Active().Name == remote.Name {
		isActive = "true"
	}

	client := DefaultAgentClientFactory.New(remote)
	metadata, err := client.GetMetadata()
	if err != nil {
		return PlainOutput{}, err
	}

	adapterMetadataBytes, err := json.Marshal(metadata.Adapter)
	if err != nil {
		return PlainOutput{}, err
	}

	adapterMetadata := struct {
		Version   string
		Type      string
		IsHealthy bool
	}{}
	if err := json.Unmarshal(adapterMetadataBytes, &adapterMetadata); err != nil {
		return PlainOutput{}, err
	}

	do := DetailOutput{
		Details: map[string]string{
			"Name":               remote.Name,
			"Active":             isActive,
			"Endpoint":           remote.Endpoint,
			"Agent Version":      metadata.Agent.Version,
			"Adapter Version":    adapterMetadata.Version,
			"Adapter Type":       adapterMetadata.Type,
			"Adapter Is Healthy": strconv.FormatBool(adapterMetadata.IsHealthy),
		},
	}

	lo, err := ListDeployments(remote)
	if err != nil {
		return PlainOutput{}, err
	}

	co := CombinedOutput{}
	co.AddOutput("", do)
	co.AddOutput("Deployments", lo)

	return &co, nil
}

func SetActiveRemote(config config.Config, name string) (Output, error) {
	if err := config.SetActive(name); err != nil {
		return PlainOutput{}, err
	}
	return PlainOutput{"Success!"}, nil
}
