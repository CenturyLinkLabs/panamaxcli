package actions

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/CenturyLinkLabs/panamaxcli/config"
)

var format = regexp.MustCompile("^[a-zA-Z0-9]+$")

func AddRemote(config config.Config, name string, path string) (string, error) {
	if !format.MatchString(name) {
		return "", errors.New("Invalid name")
	}
	if config.Exists(name) {
		return "", errors.New("Name already exists")
	}
	token, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	trimmedToken := strings.TrimSpace(string(token))
	if err = config.Save(name, trimmedToken); err != nil {
		return "", err
	}
	return "Success!", nil
}
