package actions

import (
	"errors"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/stretchr/testify/assert"
)

func TestListDeployments(t *testing.T) {
	setupFactory()
	r := config.Remote{Name: "Test"}
	fakeClient.Deployments = []agent.DeploymentResponseLite{
		{Name: "Test", ID: 1, ServiceIDs: []string{"wp", "db"}},
	}
	o, err := ListDeployments(r)

	assert.NoError(t, err)

	assert.Len(t, fakeFactory.NewedRemotes, 1)
	lo, ok := o.(*ListOutput)
	if assert.True(t, ok) && assert.Len(t, lo.Rows, 1) {
		assert.Equal(t, "Test", lo.Rows[0]["Name"])
		assert.Equal(t, "1", lo.Rows[0]["ID"])
		assert.Equal(t, "2", lo.Rows[0]["Services"])
	}
}

func TestListDeploymentsErrored(t *testing.T) {
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
			Services: []agent.Service{{ID: "wp", ActualState: "running"}},
		},
	}
	o, err := DescribeDeployment(r, "1")

	assert.NoError(t, err)

	assert.Equal(t, "1", fakeClient.DescribedDeployment)
	assert.Len(t, fakeFactory.NewedRemotes, 1)

	co, ok := o.(*CombinedOutput)
	if assert.True(t, ok) && assert.Len(t, co.Outputs, 2) {
		do, ok := co.Outputs[0].Output.(DetailOutput)
		if assert.True(t, ok) {
			assert.Equal(t, "Test", do.Details["Name"])
			assert.Equal(t, "1", do.Details["ID"])
		}

		lo, ok := co.Outputs[1].Output.(ListOutput)
		if assert.True(t, ok) && assert.Len(t, lo.Rows, 1) {
			assert.Equal(t, "wp", lo.Rows[0]["ID"])
			assert.Equal(t, "running", lo.Rows[0]["State"])
		}
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

func TestRedeployDeployment(t *testing.T) {
	setupFactory()
	fakeClient.RedeploymentResponse = agent.DeploymentResponseLite{
		ID:         1,
		Name:       "Test Name",
		ServiceIDs: []string{"wp", "db"},
	}
	r := config.Remote{}
	o, err := RedeployDeployment(r, "1")
	assert.Equal(t, "1", fakeClient.RedeployedDeployment)
	assert.NoError(t, err)
	assert.Equal(t, "Redeployed 'Test Name', services: wp, db", o.ToPrettyOutput())
}

func TestRedeployDeploymentErrored(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	fakeClient.ErrorForDeploymentRedeploy = errors.New("Errored Redeploy")
	o, err := RedeployDeployment(r, "Bad ID")

	assert.Equal(t, PlainOutput{}, o)
	assert.EqualError(t, err, "Errored Redeploy")
}

func TestDeleteDeployment(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	o, err := DeleteDeployment(r, "1")

	assert.Equal(t, "1", fakeClient.DeletedDeployment)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully deleted deployment '1'", o.ToPrettyOutput())
}

func TestErroredDeleteDeployment(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	fakeClient.ErrorForDeploymentDelete = errors.New("Delete Error")
	o, err := DeleteDeployment(r, "1")

	assert.EqualError(t, err, "Delete Error")
	assert.Equal(t, "", o.ToPrettyOutput())
}
