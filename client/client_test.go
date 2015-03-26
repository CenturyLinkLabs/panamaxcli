package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetApps(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.Handle("/apps", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		json := `[
			{
				"categories": [
				{
					"id": 1,
					"name": "Web Tier",
					"position": null
				},
				{
					"id": 2,
					"name": "DB Tier",
					"position": null
				}
				],
				"documentation": "Wordpress with MySQL\nFoo",
				"errors": null,
				"from": "Template: Wordpress with MySQL",
				"id": 1,
				"name": "Wordpress with MySQL"
			}
		]`
		fmt.Fprintf(w, json)
	}))

	p := PanamaxAPI{URL: s.URL}
	apps, err := p.GetApps()
	assert.NoError(t, err)
	if assert.Len(t, apps, 1) {
		assert.Equal(t, 1, apps[0].ID)
		assert.Equal(t, "Wordpress with MySQL", apps[0].Name)
	}
}

func TestErroredJSONDecodeGetApps(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.Handle("/apps", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "BAD JSON")
	}))

	p := PanamaxAPI{URL: s.URL}
	apps, err := p.GetApps()
	assert.Empty(t, apps)
	assert.EqualError(t, err, "error decoding JSON")
}

func TestErroredMissingURLGetApps(t *testing.T) {
	p := PanamaxAPI{URL: ""}
	apps, err := p.GetApps()
	assert.Empty(t, apps)
	assert.NotNil(t, err)
}

func TestGetApp(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.Handle("/apps/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		json := `{
				"categories": [
				{
					"id": 1,
					"name": "Web Tier",
					"position": null
				},
				{
					"id": 2,
					"name": "DB Tier",
					"position": null
				}
				],
				"documentation": "Wordpress with MySQL\nFoo",
				"errors": null,
				"from": "Template: Wordpress with MySQL",
				"id": 1,
				"name": "Wordpress with MySQL"
		}`
		fmt.Fprintf(w, json)
	}))

	p := PanamaxAPI{URL: s.URL}
	app, err := p.GetApp(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, app.ID)
	assert.Equal(t, "Wordpress with MySQL", app.Name)
}

func TestErroredRequestDoGet(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.Handle("/apps", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}))

	res, err := doGet(s.URL + "/apps")
	assert.Empty(t, res)
	assert.EqualError(t, err, "unexpected status '500 Internal Server Error'")
}
