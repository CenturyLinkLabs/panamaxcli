package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplementsPanamaxAPI(t *testing.T) {
	assert.Implements(t, (*PanamaxClient)(nil), new(PanamaxAPI))
}

func TestGetApps(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.Handle("/apps", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		fmt.Fprintf(w, appsJSON)
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
		fmt.Fprintf(w, appJSON)
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

func TestSearchTemplates(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.Handle("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//\?search_result_set%5Bq%5D\=redis\&search_result_set%5Btype%5D\=template\&search_result_set%5Blimit%5D\=40
		// TODO test terms and types and limit I guess?
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		fmt.Fprintf(w, templateSearchJSON)
	}))

	p := PanamaxAPI{URL: s.URL}
	templates, err := p.SearchTemplates("wordpress mysql")
	assert.NoError(t, err)
	if assert.Len(t, templates, 1) {
		assert.Equal(t, 1, templates[0].ID)
		assert.Equal(t, "Wordpress with MySQL", templates[0].Name)
	}
}
