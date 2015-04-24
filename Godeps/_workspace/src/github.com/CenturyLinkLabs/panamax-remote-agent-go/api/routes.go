package api

import (
	"net/http"
	"strings"
)
import "github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(agent.Manager, http.ResponseWriter, *http.Request)
}

const (
	metadataURN    = "/metadata"
	deploymentsURN = "/deployments"
	deploymentURN  = deploymentsURN + "/{id}"
	reDeployURN    = deploymentURN + "/redeploy"
)

var routes = []route{
	{
		"showDeployment",
		"GET",
		deploymentURN,
		ShowDeployment,
	},
	{
		"listDeployments",
		"GET",
		deploymentsURN,
		ListDeployments,
	},
	{
		"createDeployment",
		"POST",
		deploymentsURN,
		CreateDeployment,
	},
	{
		"deleteDeployment",
		"DELETE",
		deploymentURN,
		DeleteDeployment,
	},
	{
		"reDeploy",
		"POST",
		reDeployURN,
		ReDeploy,
	},
	{
		"metadata",
		"GET",
		metadataURN,
		Metadata,
	},
}

// URLForDeploymentID returns the URL for a specific deployment.
func URLForDeploymentID(id string) string {
	return strings.Replace(deploymentURN, "{id}", id, 1)
}

// RedeploymentURLForDeploymentID returns the URL to redeploy a specific
// deployment.
func RedeploymentURLForDeploymentID(id string) string {
	return strings.Replace(reDeployURN, "{id}", id, 1)
}

// URLForDeployments returns the URL for all deployments.
func URLForDeployments() string {
	return deploymentsURN
}

// URLForMetadata returns the URL for metadata.
func URLForMetadata() string {
	return metadataURN
}
