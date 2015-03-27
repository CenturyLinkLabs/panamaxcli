package main // import "github.com/CenturyLinkLabs/panamaxcli"

import (
	"errors"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var Commands []cli.Command

func init() {
	Commands = []cli.Command{
		{
			Name:  "remote",
			Usage: "Manage remotes",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List remotes",
					Action: noopAction,
				},
				{
					Name:        "add",
					Usage:       "Add a remote",
					Description: "Argument is the path to a plain text file containing a token.",
					Before:      actionRequiresArgument("file path"),
					Action:      noopAction,
				},
				{
					Name:        "active",
					Usage:       "Set or get the active remote",
					Description: "Passing a remote name as an argument makes it the active remote.",
					Before:      actionRequiresArgument("remote name"),
					Action:      noopAction,
				},
				{
					Name:        "remove",
					Usage:       "Remove a remote",
					Description: "Argument is a remote name.",
					Before:      actionRequiresArgument("remote name"),
					Action:      noopAction,
				},
			},
		},
		{
			Name:  "deployment",
			Usage: "Manage deployments",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List deployments",
					Action: noopAction,
				},
				{
					Name:        "describe",
					Usage:       "Describe a deployment",
					Description: "Argument is a deployment ID.",
					Before:      actionRequiresArgument("deployment ID"),
					Action:      noopAction,
				},
				{
					Name:        "redeploy",
					Usage:       "Redeploy a deployment",
					Description: "Argument is a deployment ID.",
					Before:      actionRequiresArgument("deployment ID"),
					Action:      noopAction,
				},
				{
					Name:        "delete",
					Usage:       "Delete a deployment",
					Description: "Argument is a deployment ID.",
					Before:      actionRequiresArgument("deployment ID"),
					Action:      noopAction,
				},
			},
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "panamaxcli"
	app.Version = "0.0.1"
	app.Usage = "Panamax command-line utility."
	app.Authors = []cli.Author{{"CenturyLink Labs", "clt-labs-futuretech@centurylink.com"}}
	app.Commands = Commands

	app.Run(os.Args)
}

func actionRequiresArgument(name string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.Args().First() == "" {
			message := fmt.Sprintf("This command requires a %s as an argument.", name)
			log.Errorln(message)
			return errors.New(message)
		}

		return nil
	}
}

func noopAction(c *cli.Context) {
	fmt.Println("This command is not implemented.")
}
