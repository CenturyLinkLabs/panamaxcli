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
	if _, err := config.Get(name); err == nil {
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

func RemoveRemote(config config.Config, name string) (Output, error) {
	if err := config.Remove(name); err != nil {
		return PlainOutput{}, err
	}
	out := fmt.Sprintf("Successfully removed remote '%s' from configuration!", name)
	return PlainOutput{out}, nil
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
	r, err := c.Get(name)
	if err != nil {
		return PlainOutput{}, err
	}

	isActive := "false"
	if c.Active() != nil && c.Active().Name == r.Name {
		isActive = "true"
	}

	client := DefaultAgentClientFactory.New(r)
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
			"Name":               r.Name,
			"Active":             isActive,
			"Endpoint":           r.Endpoint,
			"Agent Version":      metadata.Agent.Version,
			"Adapter Version":    adapterMetadata.Version,
			"Adapter Type":       adapterMetadata.Type,
			"Adapter Is Healthy": strconv.FormatBool(adapterMetadata.IsHealthy),
		},
	}

	lo, err := ListDeployments(r)
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
	return PlainOutput{fmt.Sprintf("'%s' is now your active remote!", name)}, nil
}

func GetRemoteToken(c config.Config, name string) (Output, error) {
	r, err := c.Get(name)
	if err != nil {
		return PlainOutput{}, err
	}
	return PlainOutput{r.Token}, nil
}
