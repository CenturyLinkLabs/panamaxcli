package actions

import (
	"errors"
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

func TestDeploymentListErrored(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	fakeClient.ErrorForDeploymentList = errors.New("Errored Deployment List")
	o, err := ListDeployments(r)

	assert.EqualError(t, err, "Errored Deployment List")
	assert.Equal(t, PlainOutput{}, o)
}

func TestDeploymentListEmpty(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	o, err := ListDeployments(r)

	assert.NoError(t, err)
	assert.Equal(t, PlainOutput{"No Deployments"}, o)
}

func TestDescribeDeployment(t *testing.T) {
	setupFactory()
	r := config.Remote{Name: "Test"}
	fakeClient.DeploymentDescription = agent.DeploymentResponseFull{
		Name: "Test",
		ID:   1,
		Status: agent.Status{
			Services: []agent.Service{
				{ID: "1", ActualState: "running"},
				{ID: "2", ActualState: "exploded"},
			},
		},
	}
	o, err := DescribeDeployment(r, "1")

	assert.NoError(t, err)

	assert.Len(t, fakeFactory.NewedRemotes, 1)
	do, ok := o.(*DetailOutput)
	if assert.True(t, ok) {
		assert.Equal(t, "Test", do.Details["Name"])
		assert.Equal(t, "1", do.Details["ID"])
		assert.Equal(t, "running, exploded", do.Details["Service Statuses"])
	}
}

func TestDescribeDeploymentErrored(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	fakeClient.ErrorForDeploymentDescription = errors.New("Errored Deployment List")
	o, err := DescribeDeployment(r, "Bad ID")

	assert.EqualError(t, err, "Errored Deployment List")
	assert.Equal(t, PlainOutput{}, o)
}
