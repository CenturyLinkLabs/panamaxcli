package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config interface {
	Save(name string, token string) error
	Exists(name string) bool
}

type FileConfig struct {
	Path   string
	Agents []Agent
}

type Agent struct {
	Name  string
	Token string
}

func (c *FileConfig) Save(name string, token string) error {
	a := Agent{name, token}
	c.Agents = append(c.Agents, a)
	b, err := json.MarshalIndent(c.Agents, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.Path, b, 0600)
	return err
}

func (c *FileConfig) Exists(name string) bool {
	for _, a := range c.Agents {
		if a.Name == name {
			return true
		}
	}
	return false
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
	if err := d.Decode(&c.Agents); err != nil {
		return err
	}

	return nil
}
