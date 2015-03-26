package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type App struct {
	ID   int
	Name string
}

type PanamaxAPI struct {
}

func (p PanamaxAPI) GetApps() []App {
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

	return apps
}

type Panamax interface {
	GetApps() []App
}

func ListApps(p Panamax) (string, error) {
	apps := p.GetApps()
	out := "Running Apps\n"
	for _, app := range apps {
		out += fmt.Sprintf("App: %d, %s\n", app.ID, app.Name)
	}
	return out, nil
}
