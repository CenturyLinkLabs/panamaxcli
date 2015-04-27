package main // import "github.com/CenturyLinkLabs/panamaxcli"

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/client"
	"github.com/CenturyLinkLabs/panamaxcli/actions"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	Config   config.Config
	Commands []cli.Command
)

const verificationWarning = `There was a problem verifying the Panamax Agent's SSL certificate! Please check the README if you are unsure why this might occur:

https://github.com/CenturyLinkLabs/panamaxcli

If you're positive that this is not an issue, you can rerun your command with the --insecure flag. The error is:
%s`

func init() {
	client.DefaultHTTPTimeout = 10

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
					Name:        "describe",
					Aliases:     []string{"d"},
					Usage:       "Describe a remote",
					Description: "Arguments is optionally the name of the remote. When omitted, the active remote will be used.",
					Before:      actionRequiresArgument("optional:remote name"),
					Action:      remoteDescribeAction,
				},
				{
					Name:        "add",
					Usage:       "Add a remote",
					Description: "Arguments are the name of the remote and the path to the token file.",
					Before:      actionRequiresArgument("remote name", "token path"),
					Action:      remoteAddAction,
				},
				{
					Name:        "active",
					Usage:       "Set the active remote",
					Description: "Argument is a remote name.",
					Before:      actionRequiresArgument("remote name"),
					Action:      setActiveRemoteAction,
				},
				{
					Name:        "remove",
					Usage:       "Remove a remote",
					Description: "Argument is a remote name.",
					Before:      actionRequiresArgument("remote name"),
					Action:      removeRemoteAction,
				},
				{
					Name:        "token",
					Usage:       "Show the remote's token",
					Description: "Argument is a remote name.",
					Before:      actionRequiresArgument("optional:remote name"),
					Action:      getTokenAction,
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
					Name:        "create",
					Usage:       "Deploy a template",
					Description: "Argument is the path to a Panamax template.",
					Before:      actionRequiresArgument("template path"),
					Action:      createDeploymentAction,
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
	app.Before = initializeApp
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable verbose logging",
		},
		cli.BoolFlag{
			Name:  "insecure",
			Usage: "Skip SSL certificate verification",
		},
	}

	app.Run(os.Args)
}

func initializeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		// Remote Agent Client will write to logrus with helpful info!
		log.SetLevel(log.DebugLevel)
	}

	if c.GlobalBool("insecure") {
		// Remote Agent Client will not verify SSL cert. This is probably bad, but
		// is useful for old certs that don't have the proper SAN IP settings.
		client.SkipSSLVerify = true
	}

	// Surprise! CLI wants an error from this method but, only uses it to abort
	// execution, not for display anywhere.
	if err := loadConfig(c); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func loadConfig(c *cli.Context) error {
	path, err := makeConfigPath()
	if err != nil {
		return err
	}

	fileConfig := config.FileConfig{Path: path}
	if err := fileConfig.Load(); err != nil {
		return err
	}
	Config = config.Config(&fileConfig)

	return nil
}

func actionRequiresArgument(args ...string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		requiredCount := len(args)
		for i, arg := range args {
			if strings.HasPrefix(arg, "optional:") {
				requiredCount = i
			}
		}

		if len(c.Args()) < requiredCount || len(c.Args()) > len(args) {
			s := strings.Join(args, ", ")
			message := fmt.Sprintf("This command requires the following arguments: %s", s)
			log.Errorln(message)
			return errors.New(message)
		}

		return nil
	}
}

func actionRequiresActiveRemote(c *cli.Context) error {
	arg := c.Args().First()
	isHelp := (arg == "help" || arg == "h")
	if !isHelp && Config.Active() == nil {
		message := "an active remote is required for this command"
		log.Errorln(message)
		return errors.New(message)
	}

	return nil
}

func remoteAddAction(c *cli.Context) {
	name := c.Args().First()
	path := c.Args().Get(1)

	output, err := actions.AddRemoteByPath(Config, name, path)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func removeRemoteAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.RemoveRemote(Config, name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func remoteListAction(c *cli.Context) {
	output := actions.ListRemotes(Config)
	fmt.Println(output.ToPrettyOutput())
}

func remoteDescribeAction(c *cli.Context) {
	name, err := explicitOrActiveRemoteName(c)
	if err != nil {
		log.Error(err)
		return
	}

	output, err := actions.DescribeRemote(Config, name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func setActiveRemoteAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.SetActiveRemote(Config, name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func deploymentsListAction(c *cli.Context) {
	output, err := actions.ListDeployments(*Config.Active())

	if err != nil {
		fatalError(err)
	}
	fmt.Println(output.ToPrettyOutput())
}

func createDeploymentAction(c *cli.Context) {
	path := c.Args().First()
	output, err := actions.CreateDeployment(*Config.Active(), path)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func describeDeploymentAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.DescribeDeployment(*Config.Active(), name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func redeployDeploymentAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.RedeployDeployment(*Config.Active(), name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func deleteDeploymentAction(c *cli.Context) {
	name := c.Args().First()
	output, err := actions.DeleteDeployment(*Config.Active(), name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func getTokenAction(c *cli.Context) {
	name, err := explicitOrActiveRemoteName(c)
	if err != nil {
		log.Error(err)
		return
	}

	output, err := actions.GetRemoteToken(Config, name)
	if err != nil {
		fatalError(err)
	}

	fmt.Println(output.ToPrettyOutput())
}

func fatalError(err error) {
	if uErr, ok := err.(*url.Error); ok {
		if hErr, ok := uErr.Err.(x509.HostnameError); ok {
			err = fmt.Errorf(verificationWarning, hErr.Error())
		}
	}

	log.Fatal(err)
}

func explicitOrActiveRemoteName(c *cli.Context) (string, error) {
	if c.Args().First() == "" && Config.Active() == nil {
		return "", errors.New("you must provide a remote name or set an active remote!")
	}

	name := c.Args().First()
	if name == "" {
		name = Config.Active().Name
	}

	return name, nil
}

func makeConfigPath() (string, error) {
	// Stolen from: https://github.com/awslabs/aws-sdk-go/pull/136 Originally
	// cleaner with os/user.Current(), but that failed under cross-compilation on
	// non-linux platforms.
	dir := os.Getenv("HOME") // *nix
	if dir == "" {           // Windows
		dir = os.Getenv("USERPROFILE")
	}
	if dir == "" {
		return "", errors.New("Couldn't determine your home directory!")
	}

	panamaxDir := filepath.Join(dir, ".panamax")
	if _, err := os.Stat(panamaxDir); os.IsNotExist(err) {
		if err := os.Mkdir(panamaxDir, 0600); err != nil {
			return "", err
		}
	}

	return filepath.Join(panamaxDir, "remotes"), nil
}
