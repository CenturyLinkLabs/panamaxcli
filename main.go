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

var (
	Config   config.Config
	Commands []cli.Command
)

func init() {
	Commands = []cli.Command{
		{
			Name:    "remote",
			Aliases: []string{"re"},
			Usage:   "Manage remotes",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "List remotes",
					Action:  remoteListAction,
				},
				{
					Name:    "describe",
					Aliases: []string{"d"},
					Usage:   "Describe a remote",
					Before:  actionRequiresArgument("remote name"),
					Action:  remoteDescribeAction,
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
					Usage:       "Set the active remote",
					Description: "Passing a remote name as an argument makes it the active remote.",
					Before:      actionRequiresArgument("remote name"),
					Action:      setActiveRemoteAction,
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
			Name:    "deployment",
			Aliases: []string{"de"},
			Usage:   "Manage deployments",
			Before:  actionRequiresActiveRemote,
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "List deployments",
					Action:  deploymentsListAction,
				},
				{
					Name:        "describe",
					Aliases:     []string{"d"},
					Usage:       "Describe a deployment",
					Description: "Argument is a deployment ID.",
					Before:      actionRequiresArgument("deployment ID"),
					Action:      describeDeploymentAction,
				},
				{
					Name:        "redeploy",
					Usage:       "Redeploy a deployment",
					Description: "Argument is a deployment ID.",
					Before:      actionRequiresArgument("deployment ID"),
					Action:      redeployDeploymentAction,
				},
				{
					Name:        "delete",
					Usage:       "Delete a deployment",
					Description: "Argument is a deployment ID.",
					Before:      actionRequiresArgument("deployment ID"),
					Action:      deleteDeploymentAction,
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
	app.Before = loadConfig

	app.Run(os.Args)
}

func loadConfig(c *cli.Context) error {
	usr, _ := user.Current()
	dir := usr.HomeDir
	fileConfig := config.FileConfig{Path: dir + "/.agents"}
	err := fileConfig.Load()
	if err != nil {
		log.Error(err)
		return err
	}
	Config = config.Config(&fileConfig)

	return nil
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

func actionRequiresActiveRemote(c *cli.Context) error {
	if Config.Active() == nil {
		message := "an active remote is required for this command"
		log.Errorln(message)
		return errors.New(message)
	}

	return nil
}

func noopAction(c *cli.Context) {
	fmt.Println("This command is not implemented.")
}

func remoteAddAction(c *cli.Context) {
	name := c.Args().First()
	path := c.Args().Get(1)

	output, err := actions.AddRemote(Config, name, path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func remoteListAction(c *cli.Context) {
	output := actions.ListRemotes(Config)
	fmt.Println(output.ToPrettyOutput())
}

func remoteDescribeAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.DescribeRemote(Config, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func setActiveRemoteAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.SetActiveRemote(Config, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func deploymentsListAction(c *cli.Context) {
	output, err := actions.ListDeployments(*Config.Active())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output.ToPrettyOutput())
}

func describeDeploymentAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.DescribeDeployment(*Config.Active(), name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func redeployDeploymentAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.RedeployDeployment(*Config.Active(), name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func deleteDeploymentAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.DeleteDeployment(*Config.Active(), name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output.ToPrettyOutput())
}
