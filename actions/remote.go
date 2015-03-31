package actions

import (
	"errors"
	"io/ioutil"
	"regexp"
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
	return PlainOutput{"Success!"}, nil
}

func ListRemotes(config config.Config) Output {
	agents := config.Remotes()
	if len(agents) == 0 {
		return PlainOutput{"No remotes"}
	}

	output := ListOutput{Labels: []string{"Active", "Name"}}
	for _, r := range config.Remotes() {
		activeMarker := ""
		if *config.Active() == r {
			activeMarker = "*"
		}

		output.AddRow(map[string]string{
			"Active": activeMarker,
			"Name":   r.Name,
		})
	}
	return &output
}

func SetActiveRemote(config config.Config, name string) (Output, error) {
	if err := config.SetActive(name); err != nil {
		return PlainOutput{}, err
	}
	return PlainOutput{"Success!"}, nil
}
