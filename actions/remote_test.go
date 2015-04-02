package actions

import (
	"errors"
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
	SavedToken          string
	ErrorForSave        error
	ActiveRemote        *config.Remote
	ActivatedRemoteName string
	ErrorForSetActive   error
}

func (c *FakeConfig) Save(name string, token string) error {
	c.SavedName = name
	c.SavedToken = token
	return c.ErrorForSave
}

func (c *FakeConfig) Exists(name string) bool {
	for _, a := range c.Agents {
		if a.Name == name {
			return true
		}
	}

	return false
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
	assert.Equal(t, "Success!", output.ToPrettyOutput())
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

	assert.Equal(t, "", output.ToPrettyOutput())
	assert.EqualError(t, err, "Name already exists")
}

func TestErroredInvalidNameAddRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "bad name", "unused")
	assert.Equal(t, "", output.ToPrettyOutput())
	assert.EqualError(t, err, "Invalid name")

	_, err = AddRemote(&fc, "bad!", "unused")
	assert.EqualError(t, err, "Invalid name")

	_, err = AddRemote(&fc, "bad/", "unused")
	assert.EqualError(t, err, "Invalid name")
}

func TestErroredMissingFileAddRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "name", "bad/file")
	assert.Equal(t, "", output.ToPrettyOutput())
	assert.EqualError(t, err, "open bad/file: no such file or directory")
}

func TestErroredConfigSaveAddRemote(t *testing.T) {
	tokenFilePath := setupTokenFile(t, "token data")
	defer os.Remove(tokenFilePath)
	fc := FakeConfig{ErrorForSave: errors.New("test error")}
	output, err := AddRemote(&fc, "name", tokenFilePath)

	assert.Equal(t, "", output.ToPrettyOutput())
	assert.EqualError(t, err, "test error")
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
		assert.Equal(t, "", lo.Rows[0]["Active"])
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
	assert.Equal(t, "Success!", output.ToPrettyOutput())
}

func TestDescribeRemote(t *testing.T) {
	setupFactory()
	fakeClient.Metadata = agent.Metadata{
		Agent: struct {
			Version string `json:"version"`
		}{Version: "0.1"},
		Adapter: struct {
			Version   string `json:version`
			Type      string `json:type`
			IsHealthy bool   `json:isHealthy`
		}{"0.2", "Test", true},
	}
	r := config.Remote{Name: "Test", Endpoint: "http://example.com"}
	fc := FakeConfig{Agents: []config.Remote{r}}
	output, err := DescribeRemote(&fc, "Test")

	assert.NoError(t, err)
	if assert.Len(t, fakeFactory.NewedRemotes, 1) {
		assert.Equal(t, "Test", fakeFactory.NewedRemotes[0].Name)
	}

	do, ok := output.(*DetailOutput)
	if assert.True(t, ok) {
		assert.Equal(t, "false", do.Details["Active"])
		assert.Equal(t, "Test", do.Details["Name"])
		assert.Equal(t, "http://example.com", do.Details["Endpoint"])
		assert.Equal(t, "0.1", do.Details["Agent Version"])
		assert.Equal(t, "0.2", do.Details["Adapter Version"])
		assert.Equal(t, "Test", do.Details["Adapter Type"])
		assert.Equal(t, "true", do.Details["Adapter Is Healthy"])
	}
}

func TestErroredClientDescribeRemote(t *testing.T) {
	setupFactory()
	fakeClient.ErrorForMetadata = errors.New("test error")
	r := config.Remote{Name: "Test", Endpoint: "http://example.com"}
	fc := FakeConfig{Agents: []config.Remote{r}}

	output, err := DescribeRemote(&fc, "Test")
	assert.Equal(t, "", output.ToPrettyOutput())
	assert.EqualError(t, err, "test error")
}

func TestErroredNonexistantDescribeRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := DescribeRemote(&fc, "Nonexistant")

	assert.EqualError(t, err, "the remote 'Nonexistant' does not exist")
	assert.Equal(t, "", output.ToPrettyOutput())
}

func TestErroredSetActiveRemote(t *testing.T) {
	fc := FakeConfig{ErrorForSetActive: errors.New("Name Not Found")}
	output, err := SetActiveRemote(&fc, "Bad")

	assert.Error(t, err)
	assert.Equal(t, "", output.ToPrettyOutput())
}
