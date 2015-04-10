package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
	"github.com/stretchr/testify/assert"
)

var (
	server        *httptest.Server
	adapterServer *httptest.Server
	baseURI       string
	router        http.Handler
	ps            agent.DeploymentStore
)

func init() {
	ps, _ = repo.MakeDeploymentStore("../db/agent_test.db")
}

func setup(hdlr http.Handler) {
	adapterServer = httptest.NewServer(hdlr)
	c := adapter.MakeClient(adapterServer.URL)
	dm := agent.MakeDeploymentManager(ps, c, "v1")
	router = MakeServer(dm, "user", "pass", "", "").newRouter()
	server = httptest.NewServer(router)
	baseURI = server.URL
}

func teardown() {
	server.Close()
}

func getAllDeployments() []agent.DeploymentResponseLite {

	res, err := doGET(baseURI + "/deployments")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var drs []agent.DeploymentResponseLite
	jd := json.NewDecoder(res.Body)
	if err := jd.Decode(&drs); err != nil {
		panic(err)
	}

	return drs
}

func doGET(url string) (*http.Response, error) {
	c := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("user", "pass")

	res, errr := c.Do(req)
	if errr != nil {
		return nil, errr
	}

	return res, nil
}

func doPOST(url string, r io.Reader) (*http.Response, error) {
	c := &http.Client{}

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("user", "pass")

	res, errr := c.Do(req)
	if errr != nil {
		return nil, errr
	}

	return res, nil
}

func doDELETE(url string) (*http.Response, error) {
	c := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("user", "pass")

	res, errr := c.Do(req)
	if errr != nil {
		return nil, errr
	}

	return res, nil
}

func removeAll() {
	drs := getAllDeployments()

	for _, dr := range drs {
		url := fmt.Sprintf("%s/deployments/%d", baseURI, dr.ID)
		doDELETE(url)
	}
}

func TestListDeploymentsWhenNoDeploymentsExist(t *testing.T) {
	setup(nil)
	defer teardown()
	removeAll()

	res, _ := doGET(baseURI + "/deployments")
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, "[]", strings.TrimSpace(string(body)))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header["Content-Type"][0])
}

func TestCreateDeployment(t *testing.T) {
	var resBody []byte
	var rURL string
	setup(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody, _ = ioutil.ReadAll(r.Body)
		rURL = r.URL.Path

		drs := []adapter.Service{
			{ID: "wp-pod"},
			{ID: "mysql-pod"},
			{ID: "honey-pod"},
		}

		drsj, err := json.Marshal(drs)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(drsj)

	}))
	defer teardown()

	buf := strings.NewReader(`{
		"override":{
			"images":[
				{
					"name":"wp",
					"environment":[
						{ "variable":"DB_PASSWORD", "value":"pass@word02" }
					],
					"deployment":{ "count":1 }
				},
				{
					"name":"mysql",
					"environment":[
						{ "variable":"MYSQL_ROOT_PASSWORD", "value":"pass@word02" }
					]
				}
			]
		},
		"template":{
			"name": "fooya",
			"images":[
				{
					"name":"wp",
					"source":"centurylink/wordpress:3.9.1",
					"links":[
						{ "service":"mysql", "alias":"DB_1" }
					],
					"ports":[
						{ "host_port":"8000", "container_port":"80" }
					],
					"environment":[
						{ "variable":"DB_PASSWORD", "value":"pass@word01" },
						{ "variable":"DB_NAME", "value":"wordpress" }
					],
					"command":"./run.sh"
				},
				{
					"name":"mysql",
					"source":"centurylink/mysql:5.5",
					"environment":[
						{ "variable":"MYSQL_ROOT_PASSWORD", "value":"pass@word01" }
					],
					"ports":[
						{ "host_port":3306, "container_port":3306 }
					],
					"expose": [1234, 5678],
					"volumes": [
						{"host_path":"foo/bar", "container_path":"/var/bar"}
					],
					"volumesFrom":["wp"],
					"command":"./run.sh"
				},
				{
					"name":"honeybadger",
					"source":"honey/badger"
				}
			]
		}
	}`)

	res, err := doPOST(baseURI+"/deployments", buf)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	jd := json.NewDecoder(res.Body)
	dr := &agent.DeploymentResponseLite{}
	if err := jd.Decode(dr); err != nil {
		panic(err)
	}

	b := []byte(`[
		{
			"command":"./run.sh",
			"deployment":{"count":1},
			"environment":[{"variable":"DB_PASSWORD","value":"pass@word02"},{"variable":"DB_NAME","value":"wordpress"}],
			"links":[{"alias":"DB_1","name":"mysql"}],
			"name":"wp",
			"ports":[{"containerPort":80,"hostPort":8000}],
			"source":"centurylink/wordpress:3.9.1"
		},
		{
			"command":"./run.sh",
			"environment":[{"variable":"MYSQL_ROOT_PASSWORD","value":"pass@word02"}],
			"expose":[1234,5678],
			"name":"mysql",
			"ports":[{"containerPort":3306,"hostPort":3306}],
			"source":"centurylink/mysql:5.5",
			"volumes":[{"containerPath":"/var/bar","hostPath":"foo/bar"}],
			"volumesFrom":["wp"]
		},
		{
			"name":"honeybadger",
			"source":"honey/badger"
		}
	]`)

	bb := bytes.Buffer{}
	json.Compact(&bb, b)
	expBody, _ := ioutil.ReadAll(&bb)

	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header["Content-Type"][0])
	assert.Equal(t, "/v1/services", rURL)
	assert.NotNil(t, dr.ID)
	assert.Equal(t, "fooya", dr.Name)
	assert.Equal(t, true, dr.Redeployable)
	assert.Equal(t, []string{"wp-pod", "mysql-pod", "honey-pod"}, dr.ServiceIDs)
	assert.Equal(t, string(expBody), strings.TrimSpace(string(resBody)))
}

func TestListDeploymentsWhenOneExists(t *testing.T) {
	setup(nil)
	defer teardown()

	res, _ := doGET(baseURI + "/deployments")
	defer res.Body.Close()

	var drs []agent.DeploymentResponseLite
	jd := json.NewDecoder(res.Body)
	if err := jd.Decode(&drs); err != nil {
		panic(err)
	}

	dr := drs[0]

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header["Content-Type"][0])
	assert.Equal(t, 1, len(drs))
	assert.Equal(t, "fooya", dr.Name)
	assert.Equal(t, true, dr.Redeployable)
	assert.Equal(t, []string{"wp-pod", "mysql-pod", "honey-pod"}, dr.ServiceIDs)
}

func TestGetDeployment(t *testing.T) {
	setup(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var srvc agent.Service
		var st int
		if r.URL.Path == "/v1/services/wp-pod" {
			srvc = agent.Service{
				ActualState: "Running",
				ID:          "wp-pod",
			}
			st = http.StatusOK
		} else if r.URL.Path == "/v1/services/honey-pod" {
			st = http.StatusInternalServerError
		} else {
			st = http.StatusNotFound
		}

		srvcj, err := json.Marshal(srvc)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(st)
		w.Write(srvcj)
	}))

	defer teardown()

	drs := getAllDeployments()

	resp, _ := doGET(fmt.Sprintf("%v/deployments/%d", baseURI, drs[0].ID))
	defer resp.Body.Close()

	dr := &agent.DeploymentResponseFull{}
	jdd := json.NewDecoder(resp.Body)
	if err := jdd.Decode(dr); err != nil {
		panic(err)
	}

	var sis []string
	var sas []string
	for _, s := range dr.Status.Services {
		sis = append(sis, s.ID)
		sas = append(sas, s.ActualState)
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header["Content-Type"][0])
	assert.Equal(t, "fooya", dr.Name)
	assert.Equal(t, true, dr.Redeployable)
	assert.Equal(t, 3, len(dr.Status.Services))

	assert.Contains(t, dr.Status.Services, agent.Service{ID: "wp-pod", ActualState: "Running"})
	assert.Contains(t, dr.Status.Services, agent.Service{ID: "mysql-pod", ActualState: "not found"})
	assert.Contains(t, dr.Status.Services, agent.Service{ID: "honey-pod", ActualState: "error"})
}

func TestReDeploy(t *testing.T) {
	setup(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			drs := []adapter.Service{
				{ID: "wp-pod"},
				{ID: "mysql-pod"},
				{ID: "honey-pod"},
			}

			drsj, err := json.Marshal(drs)
			if err != nil {
				panic(err)
			}

			w.WriteHeader(http.StatusCreated)
			w.Write(drsj)
		}
	}))

	defer teardown()

	drsPreRedeploy := getAllDeployments()

	ogID := drsPreRedeploy[0].ID
	resp, err := doPOST(fmt.Sprintf("%s/deployments/%d/redeploy", baseURI, ogID), &bytes.Buffer{})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dr := &agent.DeploymentResponseLite{}
	jdd := json.NewDecoder(resp.Body)
	if err := jdd.Decode(dr); err != nil {
		panic(err)
	}

	drsPostRedeploy := getAllDeployments()

	assert.Equal(t, 1, len(drsPreRedeploy))
	assert.Equal(t, 1, len(drsPostRedeploy))
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header["Content-Type"][0])
	assert.Equal(t, "fooya", dr.Name)
	assert.NotEqual(t, ogID, dr.ID)
	assert.Equal(t, true, dr.Redeployable)
	assert.Equal(t, 3, len(dr.ServiceIDs))
}

func TestDeleteDeployment(t *testing.T) {
	var calledURIs []string
	var calledMethods []string

	setup(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledMethods = append(calledMethods, r.Method)
		calledURIs = append(calledURIs, r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))

	defer teardown()

	drs := getAllDeployments()

	url := fmt.Sprintf("%s/deployments/%d", baseURI, drs[0].ID)
	doDELETE(url)

	drsAfterDelete := getAllDeployments()

	assert.Equal(t, 1, len(drs))
	assert.Equal(t, 0, len(drsAfterDelete))
	assert.Equal(t, []string{"DELETE", "DELETE", "DELETE"}, calledMethods)
	assert.Equal(t, len(calledURIs), 3)
	assert.Contains(t, calledURIs, "/v1/services/wp-pod")
	assert.Contains(t, calledURIs, "/v1/services/mysql-pod")
	assert.Contains(t, calledURIs, "/v1/services/honey-pod")
}

func TestGetMetadata(t *testing.T) {
	setup(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adMeta := struct {
			Boo  string
			Flee string
		}{
			Boo:  "yah",
			Flee: "foo",
		}

		jsn, _ := json.Marshal(adMeta)

		w.Write(jsn)
		w.WriteHeader(http.StatusOK)
	}))

	defer teardown()

	res, err := doGET(baseURI + "/metadata")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	ex := `{"agent":{"version":"v1"},"adapter":{"Boo":"yah","Flee":"foo"}}`

	assert.Equal(t, ex, strings.TrimSpace(string(body)))
}
