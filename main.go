package main // import "github.com/CenturyLinkLabs/panamaxcli"

import (
	"fmt"
	"os"
	"strconv"

	"github.com/CenturyLinkLabs/panamaxcli/actions"
	"github.com/CenturyLinkLabs/panamaxcli/client"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "panamaxcli"
	app.Version = "0.0.1"
	app.Usage = "Panamax command-line utility."
	app.Authors = []cli.Author{{"CenturyLink Labs", "clt-labs-futuretech@centurylink.com"}}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-url",
			Value:  "http://localhost:3001",
			Usage:  "The URL of the Panamax API server",
			EnvVar: "PANAMAX_API_URL",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "status",
			Usage:  "Get status for a running Panamax instance",
			Action: noopAction,
		},
		{
			Name:        "search",
			Usage:       "Search for templates and images",
			Description: "Argument is an search string.",
			Action:      searchAction,
		},
		{
			Name:        "run",
			Usage:       "Run an application template",
			Description: "Argument is an application template or image.",
			Action:      noopAction,
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
					Before:      appRequired,
					Action:      appDescribeAction,
				},
				{
					Name:        "logs",
					Usage:       "View an application's logs",
					Description: "Argument is an application ID.",
					Before:      appRequired,
					Action:      noopAction,
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
					Before:      serviceRequired,
					Action:      noopAction,
				},
				{
					Name:        "logs",
					Usage:       "View a service's logs",
					Description: "Argument is a service ID.",
					Before:      serviceRequired,
					Action:      noopAction,
				},
			},
		},
	}

	app.Run(os.Args)
}

func searchAction(c *cli.Context) {
	p := client.PanamaxAPI{
		URL: c.GlobalString("api-url"),
	}

	terms := c.Args()
	output, err := actions.Search(p, terms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(output)
}

func appListAction(c *cli.Context) {
	p := client.PanamaxAPI{
		URL: c.GlobalString("api-url"),
	}
	output, err := actions.ListApps(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(output)
}

func appDescribeAction(c *cli.Context) {
	p := client.PanamaxAPI{
		URL: c.GlobalString("api-url"),
	}
	appID, err := strconv.Atoi(c.Args().First())
	if err != nil {
		log.Fatal(err)
	}
	output, err := actions.DescribeApp(p, appID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(output)
}

func noopAction(c *cli.Context) {
	fmt.Println("This command is unimplemented.")
}

func appRequired(c *cli.Context) error {
	appID := c.Args().First()
	if appID == "" {
		log.Fatal("An app is required!")
	}

	return nil
}

func serviceRequired(c *cli.Context) error {
	serviceID := c.Args().First()
	if serviceID == "" {
		log.Fatal("A service is required!")
	}

	return nil
}
