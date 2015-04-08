package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLForDeploymentID(t *testing.T) {
	assert.Equal(t, "/deployments/17", URLForDeploymentID("17"))
}

func TestRedeploymentURLForDeploymentID(t *testing.T) {
	assert.Equal(t, "/deployments/17/redeploy", RedeploymentURLForDeploymentID("17"))
}

func TestURLForDeployments(t *testing.T) {
	assert.Equal(t, "/deployments", URLForDeployments())
}

func TestURLForMetadata(t *testing.T) {
	assert.Equal(t, "/metadata", URLForMetadata())
}
