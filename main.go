package main // import "github.com/CenturyLinkLabs/panamaxcli"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

type App struct {
	ID   int
	Name string
}

func main() {
	app := cli.NewApp()
	app.Name = "panamaxcli"
	app.Version = "0.0.1"
	app.Usage = "Panamax command-line utility."
	app.Authors = []cli.Author{{"CenturyLink Labs", "clt-labs-futuretech@centurylink.com"}}
	app.Commands = []cli.Command{
		{
			Name:   "app",
			Action: handleApp,
		},
	}

	app.Run(os.Args)
}

func handleApp(c *cli.Context) {
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
