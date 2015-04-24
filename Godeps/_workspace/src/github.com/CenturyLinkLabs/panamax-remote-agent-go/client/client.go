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

var (
	// DefaultHTTPTimeout does exactly what you think it does.
	DefaultHTTPTimeout = 10
	// SkipSSLVerify allows SSL certificate verification to be skipped. This is
	// frowned upon.
	SkipSSLVerify = false
)

func init() {
	log.SetLevel(log.ErrorLevel)
}

// A RequestError is a special type of error that will be returned for
// unexpected HTTP responses.
type RequestError struct {
	StatusCode int
	Body       string
}

func (e RequestError) Error() string {
	return fmt.Sprintf("unexpected status '%d'", e.StatusCode)
}

// Client represents any struct that can communicate with the Agent.
type Client interface {
	ListDeployments() ([]agent.DeploymentResponseLite, error)
	DescribeDeployment(id string) (agent.DeploymentResponseFull, error)
	CreateDeployment(b agent.DeploymentBlueprint) (agent.DeploymentResponseLite, error)
	RedeployDeployment(id string) (agent.DeploymentResponseLite, error)
	DeleteDeployment(id string) error
	GetMetadata() (agent.Metadata, error)
}

// APIClient implements the Client interface and communicates with the Agent
// over HTTPS.
type APIClient struct {
	Endpoint   string
	Username   string
	Password   string
	PrivateKey string
}

// ListDeployments fetches a list of deployments.
func (c APIClient) ListDeployments() ([]agent.DeploymentResponseLite, error) {
	var deployments []agent.DeploymentResponseLite
	err := c.doRequest("GET", api.URLForDeployments(), &deployments, nil)
	return deployments, err
}

// GetMetadata fetches metadata for the agent and adapter.
func (c APIClient) GetMetadata() (agent.Metadata, error) {
	var metadata agent.Metadata
	err := c.doRequest("GET", api.URLForMetadata(), &metadata, nil)
	return metadata, err
}

// DescribeDeployment fetches details for a specific deployment.
func (c APIClient) DescribeDeployment(id string) (agent.DeploymentResponseFull, error) {
	var resp agent.DeploymentResponseFull
	err := c.doRequest("GET", api.URLForDeploymentID(id), &resp, nil)
	return resp, err
}

// CreateDeployment creates a new deployment from a blueprint.
func (c APIClient) CreateDeployment(b agent.DeploymentBlueprint) (agent.DeploymentResponseLite, error) {
	var resp agent.DeploymentResponseLite
	err := c.doRequest("POST", api.URLForDeployments(), &resp, b)
	return resp, err
}

// RedeployDeployment redeploys a specific deployment.
func (c APIClient) RedeployDeployment(id string) (agent.DeploymentResponseLite, error) {
	var deployment agent.DeploymentResponseLite
	err := c.doRequest("POST", api.RedeploymentURLForDeploymentID(id), &deployment, nil)
	return deployment, err
}

// DeleteDeployment deletes a specific deployment.
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
			InsecureSkipVerify: SkipSSLVerify,
		},
	}

	return &http.Client{
		Timeout:   time.Duration(DefaultHTTPTimeout) * time.Second,
		Transport: verifyingTLS,
	}
}
