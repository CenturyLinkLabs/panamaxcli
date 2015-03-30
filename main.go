package main // import "github.com/CenturyLinkLabs/panamaxcli"

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/CenturyLinkLabs/panamaxcli/actions"
	"github.com/CenturyLinkLabs/panamaxcli/config"
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
					Action: remoteListAction,
				},
				{
					Name:        "add",
					Usage:       "Add a remote",
					Description: "Arguments are the name of the remote and the path to the token file.",
					Before:      actionRequiresArgument("remote name", "file path"),
					Action:      remoteAddAction,
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

func actionRequiresArgument(args ...string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if len(c.Args()) != len(args) {
			s := strings.Join(args, ", ")
			message := fmt.Sprintf("This command requires the following arguments: %s", s)
			log.Errorln(message)
			return errors.New(message)
		}

		return nil
	}
}

func noopAction(c *cli.Context) {
	fmt.Println("This command is not implemented.")
}

func remoteAddAction(c *cli.Context) {
	name := c.Args().First()
	path := c.Args().Get(1)

	usr, _ := user.Current()
	dir := usr.HomeDir
	fileConfig := config.FileConfig{Path: dir + "/.agents"}
	err := fileConfig.Load()
	if err != nil {
		log.Fatal(err)
	}
	config := config.Config(&fileConfig)

	output, err := actions.AddRemote(config, name, path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(output)
}

func remoteListAction(c *cli.Context) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	fileConfig := config.FileConfig{Path: dir + "/.agents"}
	err := fileConfig.Load()
	if err != nil {
		log.Fatal(err)
	}
	config := config.Config(&fileConfig)

	output := actions.ListRemotes(config)

	fmt.Printf(output)
}
