package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigPersistence(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{Path: dir + "/agent"}
	err = c.Save("Test Agent", "Token Data")
	assert.NoError(t, err)

	// To make sure it *really* loaded...
	c.store = Store{}
	err = c.Load()
	assert.NoError(t, err)
	if assert.Len(t, c.Remotes(), 1) {
		a := c.Remotes()[0]
		assert.Equal(t, "Test Agent", a.Name)
		assert.Equal(t, "Token Data", a.Token)
	}
}

func TestSuccessfulNonExistantLoad(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{Path: dir + "/agent"}
	err = c.Load()
	assert.NoError(t, err)
	assert.Empty(t, c.Remotes())
}

func TestErroredBadFormatLoad(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)
	c := FileConfig{Path: dir + "/agent"}
	err = ioutil.WriteFile(dir+"/agent", []byte("BAD"), 0600)
	assert.NoError(t, err)

	err = c.Load()
	assert.Contains(t, err.Error(), "invalid character")
}

func TestConfigExists(t *testing.T) {
	c := FileConfig{store: Store{Agents: []Agent{{Name: "Test"}}}}
	assert.True(t, c.Exists("Test"))
	assert.False(t, c.Exists("BadName"))
}

func TestConfigAgents(t *testing.T) {
	c := FileConfig{store: Store{Agents: []Agent{{Name: "Test"}}}}
	if assert.Len(t, c.Remotes(), 1) {
		assert.Equal(t, "Test", c.Remotes()[0].Name)
	}
}

func TestConfigSetActive(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{
		Path: dir + "/agent",
		store: Store{
			Active: "Test",
			Agents: []Agent{{Name: "Test"}, {Name: "Test2"}},
		},
	}
	assert.NoError(t, c.SetActive("Test2"))
	// To make sure it really got persisted...
	c.store = Store{}
	assert.NoError(t, c.Load())

	assert.Equal(t, "Test2", c.Active().Name)
}

func TestErroredNonexistantRemoteConfigSetActive(t *testing.T) {
	c := FileConfig{}
	err := c.SetActive("nonexistant")
	assert.EqualError(t, err, "remote 'nonexistant' does not exist")
}

func TestConfigActive(t *testing.T) {
	agent := Agent{Name: "Test"}
	c := FileConfig{store: Store{Agents: []Agent{agent}}}
	assert.Nil(t, c.Active())

	assert.NoError(t, c.SetActive("Test"))
	assert.Equal(t, &agent, c.Active())
}
