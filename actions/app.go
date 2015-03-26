package actions

import (
	"fmt"
	"strconv"

	"github.com/CenturyLinkLabs/panamaxcli/client"
)

func ListApps(p client.PanamaxClient) (string, error) {
	apps, err := p.GetApps()
	if err != nil {
		return "", err
	}
	out := "Running Apps\n"
	for _, app := range apps {
		out += fmt.Sprintf("App: %d, %s\n", app.ID, app.Name)
	}
	return out, nil
}

func DescribeApp(p client.PanamaxClient, ID int) (string, error) {
	app, err := p.GetApp(ID)
	if err != nil {
		return "", err
	}
	out := "App Details"
	out += "ID: " + strconv.Itoa(app.ID)
	out += "Name: " + app.Name
	return out, nil
}
