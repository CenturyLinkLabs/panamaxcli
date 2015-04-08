package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
)

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
	err := c.doRequest("GET", api.URLForDeployments(), &deployments)
	return deployments, err
}

func (c APIClient) GetMetadata() (agent.Metadata, error) {
	var metadata agent.Metadata
	err := c.doRequest("GET", api.URLForMetadata(), &metadata)
	return metadata, err
}

func (c APIClient) DescribeDeployment(id string) (agent.DeploymentResponseFull, error) {
	var resp agent.DeploymentResponseFull
	err := c.doRequest("GET", api.URLForDeploymentID(id), &resp)
	return resp, err
}

func (c APIClient) RedeployDeployment(id string) (agent.DeploymentResponseLite, error) {
	var deployment agent.DeploymentResponseLite
	err := c.doRequest("POST", api.RedeploymentURLForDeploymentID(id), &deployment)
	return deployment, err
}

func (c APIClient) DeleteDeployment(id string) error {
	return c.doRequest("DELETE", api.URLForDeploymentID(id), nil)
}

// TODO You need to be verifying the server cert and not ignoring it, but this
// keeps us working.
func (c APIClient) doRequest(method string, urn string, o interface{}) error {
	//pool := x509.NewCertPool()
	//pool.AppendCertsFromPEM(CertBytes)
	nonverifyingSSL := &http.Transport{
		TLSClientConfig: &tls.Config{
			//ServerName:         "X.X.X.X",
			//RootCAs:            hardcodedPool,
			InsecureSkipVerify: true,
		},
	}
	insecureHTTPClient := &http.Client{Transport: nonverifyingSSL}

	req, err := http.NewRequest(method, c.Endpoint+urn, strings.NewReader(""))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := insecureHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return RequestError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	if o == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(o)
}
