package actions

import (
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentsList(t *testing.T) {
	setupFactory()
	r := config.Remote{Name: "Test"}
	fakeClient.Deployments = []agent.DeploymentResponseLite{{Name: "Test", ID: 1}}
	o, err := ListDeployments(r)

	assert.NoError(t, err)

	assert.Len(t, fakeFactory.NewedRemotes, 1)
	lo, ok := o.(*ListOutput)
	if assert.True(t, ok) && assert.Len(t, lo.Rows, 1) {
		assert.Equal(t, "Test", lo.Rows[0]["Name"])
		assert.Equal(t, "1", lo.Rows[0]["ID"])
	}
}
