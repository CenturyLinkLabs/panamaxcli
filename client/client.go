package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type App struct {
	ID   int
	Name string
}

type PanamaxAPI struct {
	URL string
}

func (p PanamaxAPI) GetApps() ([]App, error) {
	hr, err := http.NewRequest("GET", p.URL+"/apps", nil)
	hr.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(hr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		s, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(s))
		return nil, fmt.Errorf("unexpected status '%v'", resp.Status)
	}
	defer resp.Body.Close()

	var apps []App
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apps); err != nil {
		return nil, errors.New("error decoding JSON")
	}

	return apps, nil
}

type PanamaxClient interface {
	GetApps() ([]App, error)
}
