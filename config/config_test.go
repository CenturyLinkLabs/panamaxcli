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
	c.agents = make([]Agent, 0)
	err = c.Load()
	assert.NoError(t, err)
	if assert.Len(t, c.Agents(), 1) {
		a := c.Agents()[0]
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
	assert.Empty(t, c.Agents())
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
	c := FileConfig{}
	c.agents = append(c.Agents(), Agent{Name: "Test"})
	assert.True(t, c.Exists("Test"))
	assert.False(t, c.Exists("BadName"))
}

func TestConfigAgents(t *testing.T) {
	c := FileConfig{agents: []Agent{{Name: "Test"}}}
	agents := c.Agents()
	if assert.Len(t, agents, 1) {
		assert.Equal(t, "Test", agents[0].Name)
	}
}
