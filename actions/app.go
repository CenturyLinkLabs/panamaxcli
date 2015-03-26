package actions

import (
	"fmt"

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
