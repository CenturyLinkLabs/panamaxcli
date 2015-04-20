package actions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamaxcli/config"
	"github.com/stretchr/testify/assert"
)

type FakeConfig struct {
	Agents              []config.Remote
	SavedName           string
	RemovedName         string
	SavedToken          string
	ErrorForSave        error
	ActiveRemote        *config.Remote
	ActivatedRemoteName string
	ErrorForSetActive   error
	ErrorForRemove      error
}

func (c *FakeConfig) Save(name string, token string) error {
	c.SavedName = name
	c.SavedToken = token
	c.Agents = append(c.Agents, config.Remote{Name: name, Token: token})
	return c.ErrorForSave
}

func (c *FakeConfig) Remove(name string) error {
	c.RemovedName = name
	return c.ErrorForRemove
}

func (c *FakeConfig) Get(name string) (config.Remote, error) {
	for _, r := range c.Agents {
		if r.Name == name {
			return r, nil
		}
	}

	return config.Remote{}, fmt.Errorf("the remote '%s' does not exist", name)
}

func (c *FakeConfig) Remotes() []config.Remote {
	return c.Agents
}

func (c *FakeConfig) SetActive(name string) error {
	c.ActivatedRemoteName = name
	if c.ErrorForSetActive != nil {
		return c.ErrorForSetActive
	}
	return nil
}

func (c *FakeConfig) Active() *config.Remote {
	return c.ActiveRemote
}

func setupTokenFile(t *testing.T, data string) string {
	tokenFile, err := ioutil.TempFile("", "pmx-test-token")
	tokenFile.WriteString(data)
	assert.NoError(t, err)
	return tokenFile.Name()
}

func TestAddRemote(t *testing.T) {
	tokenFilePath := setupTokenFile(t, "token data")
	defer os.Remove(tokenFilePath)
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "testname", tokenFilePath)

	assert.Equal(t, "testname", fc.SavedName)
	assert.Equal(t, "token data", fc.SavedToken)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully added! 'testname' is now your active remote.", output.ToPrettyOutput())
	assert.Equal(t, "testname", fc.ActivatedRemoteName)
}

func TestInactiveSecondAddRemote(t *testing.T) {
	tokenFilePath := setupTokenFile(t, "token data")
	defer os.Remove(tokenFilePath)
	fc := FakeConfig{Agents: []config.Remote{{Name: "Test"}}}
	output, err := AddRemote(&fc, "testname", tokenFilePath)

	assert.NoError(t, err)
	assert.Equal(t, "Successfully added!", output.ToPrettyOutput())
	assert.Empty(t, fc.ActivatedRemoteName)
}

func TestStripsWhitespaceAddRemote(t *testing.T) {
	tokenFilePath := setupTokenFile(t, "\n token data \n\n ")
	defer os.Remove(tokenFilePath)
	fc := FakeConfig{}
	_, err := AddRemote(&fc, "testname", tokenFilePath)

	assert.NoError(t, err)
	assert.Equal(t, "token data", fc.SavedToken)
}

func TestErroredExistingNameAddRemote(t *testing.T) {
	fc := FakeConfig{Agents: []config.Remote{{Name: "name"}}}
	output, err := AddRemote(&fc, "name", "unused")

	assert.Empty(t, output.ToPrettyOutput())
	assert.EqualError(t, err, "Name already exists")
}

func TestErroredInvalidNameAddRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "bad name", "unused")
	assert.Empty(t, output.ToPrettyOutput())
	assert.EqualError(t, err, "Invalid name")

	_, err = AddRemote(&fc, "bad!", "unused")
	assert.EqualError(t, err, "Invalid name")

	_, err = AddRemote(&fc, "bad/", "unused")
	assert.EqualError(t, err, "Invalid name")
}

func TestErroredMissingFileAddRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "name", "bad/file")
	assert.Empty(t, output.ToPrettyOutput())
	assert.EqualError(t, err, "open bad/file: no such file or directory")
}

func TestErroredConfigSaveAddRemote(t *testing.T) {
	tokenFilePath := setupTokenFile(t, "token data")
	defer os.Remove(tokenFilePath)
	fc := FakeConfig{ErrorForSave: errors.New("test error")}
	output, err := AddRemote(&fc, "name", tokenFilePath)

	assert.Empty(t, output.ToPrettyOutput())
	assert.EqualError(t, err, "test error")
}

func TestRemoveRemote(t *testing.T) {
	fc := FakeConfig{}
	o, err := RemoveRemote(&fc, "test")

	assert.Equal(t, "test", fc.RemovedName)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully removed remote 'test' from configuration!", o.ToPrettyOutput())
}

func TestRemoveRemoteErrored(t *testing.T) {
	fc := FakeConfig{ErrorForRemove: errors.New("Test Error")}
	o, err := RemoveRemote(&fc, "test")

	assert.EqualError(t, err, "Test Error")
	assert.Empty(t, o.ToPrettyOutput())
}

func TestListRemotes(t *testing.T) {
	active := config.Remote{Name: "Active"}
	fc := FakeConfig{
		Agents: []config.Remote{
			{Name: "Test", Endpoint: "http://example.com"},
			active,
		},
		ActiveRemote: &active,
	}
	output := ListRemotes(&fc)

	lo, ok := output.(*ListOutput)
	if assert.True(t, ok) && assert.Len(t, lo.Rows, 2) {
		assert.Empty(t, lo.Rows[0]["Active"])
		assert.Equal(t, "Test", lo.Rows[0]["Name"])
		assert.Equal(t, "http://example.com", lo.Rows[0]["Endpoint"])

		assert.Equal(t, "*", lo.Rows[1]["Active"])
		assert.Equal(t, "Active", lo.Rows[1]["Name"])
	}
}

func TestNoActiveListRemotes(t *testing.T) {
	fc := FakeConfig{Agents: []config.Remote{{Name: "Test"}}}
	assert.NotPanics(t, func() {
		ListRemotes(&fc)
	})
}

func TestListRemotesNoRemotes(t *testing.T) {
	fc := FakeConfig{}
	output := ListRemotes(&fc)
	assert.Equal(t, "No remotes", output.ToPrettyOutput())
}

func TestSetActiveRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := SetActiveRemote(&fc, "Test")

	assert.Equal(t, "Test", fc.ActivatedRemoteName)
	assert.NoError(t, err)
	assert.Equal(t, "'Test' is now your active remote!", output.ToPrettyOutput())
}

func TestDescribeRemote(t *testing.T) {
	setupFactory()
	fakeClient.Metadata = agent.Metadata{
		Agent: agent.AgentMetadata{Version: "0.1"},
		Adapter: struct {
			Version   string `json:"version"`
			Type      string `json:"type"`
			IsHealthy bool   `json:"isHealthy"`
		}{"0.2", "Test", true},
	}
	fakeClient.Deployments = []agent.DeploymentResponseLite{
		{ID: 17, Name: "Test Deployment", ServiceIDs: []string{"1", "2"}},
	}
	r := config.Remote{Name: "Test", Endpoint: "http://example.com"}
	fc := FakeConfig{Agents: []config.Remote{r}}
	output, err := DescribeRemote(&fc, "Test")

	assert.NoError(t, err)
	if assert.Len(t, fakeFactory.NewedRemotes, 2) {
		assert.Equal(t, "Test", fakeFactory.NewedRemotes[0].Name)
	}

	co, ok := output.(*CombinedOutput)
	if assert.True(t, ok) && assert.Len(t, co.Outputs, 2) {
		do, ok := co.Outputs[0].Output.(DetailOutput)
		if assert.True(t, ok) {
			assert.Equal(t, "false", do.Details["Active"])
			assert.Equal(t, "Test", do.Details["Name"])
			assert.Equal(t, "http://example.com", do.Details["Endpoint"])
			assert.Equal(t, "0.1", do.Details["Agent Version"])
			assert.Equal(t, "0.2", do.Details["Adapter Version"])
			assert.Equal(t, "Test", do.Details["Adapter Type"])
			assert.Equal(t, "true", do.Details["Adapter Is Healthy"])
		}

		lo, ok := co.Outputs[1].Output.(*ListOutput)
		if assert.True(t, ok) && assert.Len(t, lo.Rows, 1) {
			r := lo.Rows[0]
			assert.Equal(t, "17", r["ID"])
			assert.Equal(t, "Test Deployment", r["Name"])
		}
	}
}

func TestErroredClientMetadataDescribeRemote(t *testing.T) {
	setupFactory()
	fakeClient.ErrorForMetadata = errors.New("test error")
	r := config.Remote{Name: "Test", Endpoint: "http://example.com"}
	fc := FakeConfig{Agents: []config.Remote{r}}

	output, err := DescribeRemote(&fc, "Test")
	assert.Empty(t, output.ToPrettyOutput())
	assert.EqualError(t, err, "test error")
}

func TestErroredClientDeploymentListDescribeRemote(t *testing.T) {
	setupFactory()
	fakeClient.ErrorForDeploymentList = errors.New("test error deployment list")
	r := config.Remote{Name: "Test", Endpoint: "http://example.com"}
	fc := FakeConfig{Agents: []config.Remote{r}}

	output, err := DescribeRemote(&fc, "Test")
	assert.Empty(t, output.ToPrettyOutput())
	assert.EqualError(t, err, "test error deployment list")
}

func TestErroredNonexistantDescribeRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := DescribeRemote(&fc, "Nonexistant")

	assert.EqualError(t, err, "the remote 'Nonexistant' does not exist")
	assert.Empty(t, output.ToPrettyOutput())
}

func TestErroredSetActiveRemote(t *testing.T) {
	fc := FakeConfig{ErrorForSetActive: errors.New("Name Not Found")}
	output, err := SetActiveRemote(&fc, "Bad")

	assert.Error(t, err)
	assert.Empty(t, output.ToPrettyOutput())
}

func TestGetRemoteToken(t *testing.T) {
	fc := FakeConfig{
		Agents: []config.Remote{{Name: "Test", Token: "Token\nText"}},
	}
	output, err := GetRemoteToken(&fc, "Test")

	assert.NoError(t, err)
	assert.Equal(t, "Token\nText", output.ToPrettyOutput())
}

func TestErroredGetRemoteToken(t *testing.T) {
	fc := FakeConfig{}
	output, err := GetRemoteToken(&fc, "nonexistant")

	assert.EqualError(t, err, "the remote 'nonexistant' does not exist")
	assert.Empty(t, output.ToPrettyOutput())
}
