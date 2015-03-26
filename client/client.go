package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PanamaxClient interface {
	GetApps() ([]App, error)
	GetApp(id int) (App, error)
}
type App struct {
	ID   int
	Name string
}

type PanamaxAPI struct {
	URL string
}

func (p PanamaxAPI) GetApps() ([]App, error) {
	url := fmt.Sprintf("%s/apps", p.URL)
	body, err := doGet(url)
	if err != nil {
		return nil, err
	}

	var apps []App
	if err := json.Unmarshal(body, &apps); err != nil {
		return nil, errors.New("error decoding JSON")
	}
	return apps, nil
}

func (p PanamaxAPI) GetApp(id int) (App, error) {
	url := fmt.Sprintf("%s/apps/%d", p.URL, id)
	body, err := doGet(url)
	if err != nil {
		return App{}, err
	}

	var app App
	if err := json.Unmarshal(body, &app); err != nil {
		return App{}, err
	}
	return app, nil
}

func doGet(url string) ([]byte, error) {
	hr, err := http.NewRequest("GET", url, nil)
	hr.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(hr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected status '%v'", resp.Status)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
