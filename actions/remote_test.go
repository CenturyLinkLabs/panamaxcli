package actions

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeConfig struct {
	ExistingNames []string
	SavedName     string
	SavedToken    string
	ErrorForSave  error
}

func (c *FakeConfig) Save(name string, token string) error {
	c.SavedName = name
	c.SavedToken = token
	return c.ErrorForSave
}

func (c *FakeConfig) Exists(name string) bool {
	for _, s := range c.ExistingNames {
		if s == name {
			return true
		}
	}

	return false
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
	assert.Equal(t, "Success!", output)
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
	fc := FakeConfig{ExistingNames: []string{"name"}}
	output, err := AddRemote(&fc, "name", "unused")

	assert.Equal(t, "", output)
	assert.EqualError(t, err, "Name already exists")
}

func TestErroredInvalidNameAddRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "bad name", "unused")
	assert.Equal(t, "", output)
	assert.EqualError(t, err, "Invalid name")

	_, err = AddRemote(&fc, "bad!", "unused")
	assert.EqualError(t, err, "Invalid name")

	_, err = AddRemote(&fc, "bad/", "unused")
	assert.EqualError(t, err, "Invalid name")
}

func TestErroredMissingFileAddRemote(t *testing.T) {
	fc := FakeConfig{}
	output, err := AddRemote(&fc, "name", "bad/file")
	assert.Equal(t, "", output)
	assert.EqualError(t, err, "open bad/file: no such file or directory")
}

func TestErroredConfigSaveAddRemote(t *testing.T) {
	tokenFilePath := setupTokenFile(t, "token data")
	defer os.Remove(tokenFilePath)
	fc := FakeConfig{ErrorForSave: errors.New("test error")}
	output, err := AddRemote(&fc, "name", tokenFilePath)

	assert.Equal(t, "", output)
	assert.EqualError(t, err, "test error")
}
