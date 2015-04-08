package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	router *mux.Router
	client APIClient
)

func setup() {
	router = mux.NewRouter()
	server = httptest.NewServer(router)
	client = APIClient{Endpoint: server.URL, Username: "username", Password: "password"}
}

func teardown() {
	server.Close()
}

type FakeManager struct {
	Deployment   agent.DeploymentResponseFull
	Deployments  []agent.DeploymentResponseLite
	Redeployment agent.DeploymentResponseLite
	Metadata     agent.Metadata
}

func (m *FakeManager) ListDeployments() ([]agent.DeploymentResponseLite, error) {
	return m.Deployments, nil
}

func (m *FakeManager) GetFullDeployment(id int) (agent.DeploymentResponseFull, error) {
	if id == m.Deployment.ID {
		return m.Deployment, nil
	}
	return agent.DeploymentResponseFull{}, nil
}

func (m *FakeManager) GetDeployment(id int) (agent.DeploymentResponseLite, error) {
	return agent.DeploymentResponseLite{}, nil
}

func (m *FakeManager) DeleteDeployment(id int) error {
	return nil
}

func (m *FakeManager) CreateDeployment(b agent.DeploymentBlueprint) (agent.DeploymentResponseLite, error) {
	return agent.DeploymentResponseLite{}, nil
}

func (m *FakeManager) ReDeploy(id int) (agent.DeploymentResponseLite, error) {
	return m.Redeployment, nil
}

func (m *FakeManager) FetchMetadata() (agent.Metadata, error) {
	return m.Metadata, nil
}

func TestGetMetadata(t *testing.T) {
	setup()
	defer teardown()
	fakeManager := FakeManager{
		Metadata: agent.Metadata{Agent: agent.AgentMetadata{Version: "1"}},
	}

	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		api.Metadata(&fakeManager, w, r)
	}
	router.Methods("GET").Path(api.URLForMetadata()).Name("test").HandlerFunc(handler)

	m, err := client.GetMetadata()
	assert.NoError(t, err)
	assert.Equal(t, "1", m.Agent.Version)
	assert.True(t, handlerCalled)
}

func TestListDeployments(t *testing.T) {
	setup()
	defer teardown()
	drs := []agent.DeploymentResponseLite{{Name: "Test"}}
	fakeManager := FakeManager{Deployments: drs}

	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		api.ListDeployments(&fakeManager, w, r)
	}
	router.Methods("GET").Path(api.URLForDeployments()).Name("test").HandlerFunc(handler)

	d, err := client.ListDeployments()
	assert.NoError(t, err)
	if assert.Len(t, d, 1) {
		assert.Equal(t, "Test", d[0].Name)
	}
	assert.True(t, handlerCalled)
}

func TestDescribeDeployment(t *testing.T) {
	setup()
	defer teardown()
	dr := agent.DeploymentResponseFull{ID: 1, Name: "Test"}
	fakeManager := FakeManager{Deployment: dr}
	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		api.ShowDeployment(&fakeManager, w, r)
	}
	router.Methods("GET").Path(api.URLForDeploymentID("{id}")).Name("test").HandlerFunc(handler)

	d, err := client.DescribeDeployment("1")
	assert.NoError(t, err)
	assert.Equal(t, "Test", d.Name)
	assert.True(t, handlerCalled)
}

func TestRedeployDeployment(t *testing.T) {
	setup()
	defer teardown()
	dr := agent.DeploymentResponseLite{Name: "Test"}
	fakeManager := FakeManager{Redeployment: dr}
	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		api.ReDeploy(&fakeManager, w, r)
	}
	router.Methods("POST").Path(api.RedeploymentURLForDeploymentID("{id}")).Name("test").HandlerFunc(handler)

	d, err := client.RedeployDeployment("1")
	assert.NoError(t, err)
	assert.Equal(t, "Test", d.Name)
	assert.True(t, handlerCalled)
}

func TestDeleteDeployment(t *testing.T) {
	setup()
	defer teardown()
	fakeManager := FakeManager{}
	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		api.ReDeploy(&fakeManager, w, r)
	}
	router.Methods("DELETE").Path(api.URLForDeploymentID("{id}")).Name("test").HandlerFunc(handler)

	err := client.DeleteDeployment("1")
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func Test_doRequest_ErroredUnexpectedStatus(t *testing.T) {
	setup()
	defer teardown()
	router.Handle("/urn", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", 500)
	}))

	err := client.doRequest("GET", "/urn", &struct{}{})
	assert.EqualError(t, err, "unexpected status '500'")
}

func Test_doRequest_Success(t *testing.T) {
	setup()
	defer teardown()
	router.Handle("/urn", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte("username:password"))
		assert.Equal(t, []string{"application/json"}, r.Header["Content-Type"])
		assert.Equal(t, []string{"application/json"}, r.Header["Accept"])
		assert.Equal(t, []string{"Basic " + encodedAuth}, r.Header["Authorization"])

		fmt.Fprintf(w, "{}")
	}))

	err := client.doRequest("GET", "/urn", &struct{}{})
	assert.NoError(t, err)
}

func Test_doRequest_ErroredBadJSON(t *testing.T) {
	setup()
	defer teardown()
	router.Handle("/urn", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "BAD JSON")
	}))

	err := client.doRequest("GET", "/urn", &struct{}{})
	assert.Contains(t, err.Error(), "invalid character")
}
