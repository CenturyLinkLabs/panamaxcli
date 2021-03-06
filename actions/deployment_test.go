package actions

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/CenturyLinkLabs/prettycli"
	"github.com/stretchr/testify/assert"
)

func setupTemplateFile(t *testing.T, contents string) string {
	f, err := ioutil.TempFile("", "test-template.pmx")
	f.WriteString(contents)
	assert.NoError(t, err)
	return f.Name()
}

func TestListDeployments(t *testing.T) {
	setupFactory()
	r := config.Remote{Name: "Test"}
	fakeClient.Deployments = []agent.DeploymentResponseLite{
		{Name: "Test", ID: 1, ServiceIDs: []string{"wp", "db"}},
	}
	o, err := ListDeployments(r)

	assert.NoError(t, err)

	assert.Len(t, fakeFactory.NewedRemotes, 1)
	lo, ok := o.(*prettycli.ListOutput)
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
	assert.Equal(t, prettycli.PlainOutput{}, o)
}

func TestDeploymentListEmpty(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	o, err := ListDeployments(r)

	assert.NoError(t, err)
	assert.Equal(t, prettycli.PlainOutput{"No Deployments"}, o)
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

	co, ok := o.(*prettycli.CombinedOutput)
	if assert.True(t, ok) && assert.Len(t, co.Outputs, 2) {
		do, ok := co.Outputs[0].Output.(prettycli.DetailOutput)
		if assert.True(t, ok) {
			assert.Equal(t, "Test", do.Details["Name"])
			assert.Equal(t, "1", do.Details["ID"])
		}

		lo, ok := co.Outputs[1].Output.(prettycli.ListOutput)
		if assert.True(t, ok) && assert.Len(t, lo.Rows, 1) {
			assert.Equal(t, "wp", lo.Rows[0]["ID"])
			assert.Equal(t, "running", lo.Rows[0]["State"])
		}
	}
}

func TestSuccessfulCreateDeployment(t *testing.T) {
	setupFactory()
	template := setupTemplateFile(t, wordpressTemplate)
	defer os.Remove(template)
	r := config.Remote{Name: "Test"}
	fakeClient.DeployedDeployment = agent.DeploymentResponseLite{ID: 1}

	o, err := CreateDeployment(r, template)
	assert.NoError(t, err)
	images := fakeClient.DeployedBlueprint.Template.Images
	if assert.Len(t, images, 2) {
		assert.Equal(t, "WP", images[0].Name)
	}
	assert.Equal(t, "Template successfully deployed as '1'", o.ToPrettyOutput())
}

func TestErroredMissingFileCreateDeployment(t *testing.T) {
	setupFactory()
	r := config.Remote{Name: "Test"}
	o, err := CreateDeployment(r, "Bad Path")
	assert.Contains(t, err.Error(), "no such file")
	assert.Equal(t, prettycli.PlainOutput{}, o)
}

func TestErroredBadYAMLCreateDeployment(t *testing.T) {
	setupFactory()
	template := setupTemplateFile(t, "!!!?!@@@BAD YAML@@@?!?!?")
	defer os.Remove(template)
	r := config.Remote{Name: "Test"}

	o, err := CreateDeployment(r, template)
	assert.Contains(t, err.Error(), "cannot unmarshal")
	assert.Empty(t, fakeClient.DeployedBlueprint.Template.Images)
	assert.Equal(t, prettycli.PlainOutput{}, o)
}

func TestErroredClientCreateDeployment(t *testing.T) {
	setupFactory()
	fakeClient.ErrorForDeploymentCreate = errors.New("test error")
	template := setupTemplateFile(t, wordpressTemplate)
	defer os.Remove(template)
	r := config.Remote{Name: "Test"}

	o, err := CreateDeployment(r, template)
	assert.EqualError(t, err, "test error")
	assert.Equal(t, prettycli.PlainOutput{}, o)
}

func TestDescribeDeploymentErrored(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	fakeClient.ErrorForDeploymentDescription = errors.New("Errored Deployment List")
	o, err := DescribeDeployment(r, "Bad ID")

	assert.EqualError(t, err, "Errored Deployment List")
	assert.Equal(t, prettycli.PlainOutput{}, o)
}

func TestRedeployDeployment(t *testing.T) {
	setupFactory()
	fakeClient.RedeploymentResponse = agent.DeploymentResponseLite{
		ID:         2,
		Name:       "Test Name",
		ServiceIDs: []string{"wp", "db"},
	}
	r := config.Remote{}
	o, err := RedeployDeployment(r, "1")
	assert.Equal(t, "1", fakeClient.RedeployedDeployment)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"Redeployed 'Test Name' as Deployment ID 2",
		o.ToPrettyOutput(),
	)
}

func TestRedeployDeploymentErrored(t *testing.T) {
	setupFactory()
	r := config.Remote{}
	fakeClient.ErrorForDeploymentRedeploy = errors.New("Errored Redeploy")
	o, err := RedeployDeployment(r, "Bad ID")

	assert.Equal(t, prettycli.PlainOutput{}, o)
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
