package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config interface {
	Save(name string, token string) error
	Exists(name string) bool
	Remotes() []Agent
	SetActive(name string) error
	Active() *Agent
}

type FileConfig struct {
	Path  string
	store Store
}

type Store struct {
	Active string
	Agents []Agent
}

type Agent struct {
	Name  string
	Token string
}

func (c *FileConfig) Save(name string, token string) error {
	a := Agent{name, token}
	c.store.Agents = append(c.Remotes(), a)
	return c.saveAll()
}

func (c *FileConfig) Exists(name string) bool {
	for _, a := range c.Remotes() {
		if a.Name == name {
			return true
		}
	}
	return false
}

func (c *FileConfig) SetActive(name string) error {
	if !c.Exists(name) {
		return fmt.Errorf("remote '%s' does not exist", name)
	}
	c.store.Active = name
	c.saveAll()
	return nil
}

func (c *FileConfig) Active() *Agent {
	activeName := c.store.Active
	if activeName == "" {
		return nil
	}

	for _, r := range c.Remotes() {
		if r.Name == activeName {
			return &r
		}
	}

	return nil
}

func (c *FileConfig) Load() error {
	f, err := os.Open(c.Path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil
		}

		return err
	}

	d := json.NewDecoder(f)
	if err := d.Decode(&c.store); err != nil {
		return err
	}

	return nil
}

func (c *FileConfig) Remotes() []Agent {
	return c.store.Agents
}

func (c *FileConfig) saveAll() error {
	b, err := json.MarshalIndent(c.store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Path, b, 0600)
}
