package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

type App struct {
	ID   int
	Name string
}

func main() {
	app := cli.NewApp()
	app.Name = "panamaxcli"
	app.Usage = "Panamax command-line utility."
	app.Version = "0.0.1"
	app.Authors = []cli.Author{{"CenturyLink Labs", "clt-labs-futuretech@centurylink.com"}}
	app.Commands = []cli.Command{
		{
			Name:        "run",
			Usage:       "Run an application template",
			Description: "Argument is an application template name.",
			Action:      noopAction,
		},
		{
			Name:   "status",
			Usage:  "Get status for a running Panamax instance",
			Action: noopAction,
		},
		{
			Name:  "app",
			Usage: "Work with running applications",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List applications",
					Action: appListAction,
				},
				{
					Name:        "describe",
					Usage:       "Get details for a specific application",
					Description: "Argument is an application ID.",
					Action:      noopActionRequiringApp,
				},
				{
					Name:        "logs",
					Usage:       "View an application's logs",
					Description: "Argument is an application ID.",
					Action:      noopActionRequiringApp,
				},
			},
		},
		{
			Name:  "service",
			Usage: "Work with running services",
			Subcommands: []cli.Command{
				{
					Name:        "describe",
					Usage:       "Get details for a specific service",
					Description: "Argument is a service ID.",
					Action:      noopActionRequiringService,
				},
				{
					Name:        "logs",
					Usage:       "View a service's logs",
					Description: "Argument is a service ID.",
					Action:      noopActionRequiringService,
				},
			},
		},
	}

	app.Run(os.Args)
}

func appListAction(c *cli.Context) {
	resp, err := http.Get("http://coreos:3001/apps.json")
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
	defer resp.Body.Close()

	apps := make([]App, 0)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apps); err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	for _, app := range apps {
		fmt.Printf("App: %d, %s\n", app.ID, app.Name)
	}
}

func noopAction(c *cli.Context) {
	fmt.Println("This command is unimplemented.")
}

func noopActionRequiringApp(c *cli.Context) {
	appID := c.Args().First()
	if appID == "" {
		log.Fatal("An app is required!")
	}

	fmt.Println("AppID:", appID)
	fmt.Println("This command is unimplemented.")
}

func noopActionRequiringService(c *cli.Context) {
	serviceID := c.Args().First()
	if serviceID == "" {
		log.Fatal("A service is required!")
	}

	fmt.Println("ServiceID:", serviceID)
	fmt.Println("This command is unimplemented.")
}
