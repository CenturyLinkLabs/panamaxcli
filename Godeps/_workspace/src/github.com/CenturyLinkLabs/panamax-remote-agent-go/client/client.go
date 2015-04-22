package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	log "github.com/Sirupsen/logrus"
)

var DefaultHTTPTimeout = 10

func init() {
	log.SetLevel(log.ErrorLevel)
}

type RequestError struct {
	StatusCode int
	Body       string
}

func (e RequestError) Error() string {
	return fmt.Sprintf("unexpected status '%d'", e.StatusCode)
}

type Client interface {
	ListDeployments() ([]agent.DeploymentResponseLite, error)
	DescribeDeployment(id string) (agent.DeploymentResponseFull, error)
	CreateDeployment(b agent.DeploymentBlueprint) (agent.DeploymentResponseLite, error)
	RedeployDeployment(id string) (agent.DeploymentResponseLite, error)
	DeleteDeployment(id string) error
	GetMetadata() (agent.Metadata, error)
}

type APIClient struct {
	Endpoint   string
	Username   string
	Password   string
	PrivateKey string
}

func (c APIClient) ListDeployments() ([]agent.DeploymentResponseLite, error) {
	var deployments []agent.DeploymentResponseLite
	err := c.doRequest("GET", api.URLForDeployments(), &deployments, nil)
	return deployments, err
}

func (c APIClient) GetMetadata() (agent.Metadata, error) {
	var metadata agent.Metadata
	err := c.doRequest("GET", api.URLForMetadata(), &metadata, nil)
	return metadata, err
}

func (c APIClient) DescribeDeployment(id string) (agent.DeploymentResponseFull, error) {
	var resp agent.DeploymentResponseFull
	err := c.doRequest("GET", api.URLForDeploymentID(id), &resp, nil)
	return resp, err
}

func (c APIClient) CreateDeployment(b agent.DeploymentBlueprint) (agent.DeploymentResponseLite, error) {
	var resp agent.DeploymentResponseLite
	err := c.doRequest("POST", api.URLForDeployments(), &resp, b)
	return resp, err
}

func (c APIClient) RedeployDeployment(id string) (agent.DeploymentResponseLite, error) {
	var deployment agent.DeploymentResponseLite
	err := c.doRequest("POST", api.RedeploymentURLForDeploymentID(id), &deployment, nil)
	return deployment, err
}

func (c APIClient) DeleteDeployment(id string) error {
	return c.doRequest("DELETE", api.URLForDeploymentID(id), nil, nil)
}

func (c APIClient) doRequest(method string, urn string, o interface{}, p interface{}) error {
	httpClient := c.getClient()

	var params io.Reader
	var loggedParams string
	params = strings.NewReader("")
	if p != nil {
		j, err := json.Marshal(p)
		if err != nil {
			return err
		}

		loggedParams = string(j)
		params = bytes.NewReader(j)
	}

	url := c.Endpoint + urn
	req, err := http.NewRequest(method, url, params)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	log.WithFields(log.Fields{
		"URL":    url,
		"Method": method,
		"Body":   loggedParams,
	}).Info("Making request")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, " ", "  ")
	log.WithFields(log.Fields{
		"Status": resp.StatusCode,
		"Body":   prettyJSON.String(),
	}).Info("Received Response")

	if resp.StatusCode >= 400 {
		return RequestError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	if o == nil {
		return nil
	}
	return json.Unmarshal(body, &o)
}

func (c *APIClient) getClient() *http.Client {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(c.PrivateKey))
	verifyingTLS := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: false,
		},
	}

	return &http.Client{
		Timeout:   time.Duration(DefaultHTTPTimeout) * time.Second,
		Transport: verifyingTLS,
	}
}
