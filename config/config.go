package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config interface {
	Save(name string, token string) error
	Remove(name string) error
	Get(name string) (Remote, error)
	Remotes() []Remote
	SetActive(name string) error
	Active() *Remote
}

type FileConfig struct {
	Path  string
	store Store
}

type Store struct {
	Active  string   `json:"active"`
	Remotes []Remote `json:"remotes"`
}

type Remote struct {
	Name       string `json:"name"`
	Token      string `json:"token"`
	Endpoint   string `json:"endpoint"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
}

func (c *FileConfig) Save(name string, token string) error {
	r := Remote{Name: name, Token: token}
	if err := r.DecodeToken(); err != nil {
		return err
	}

	c.store.Remotes = append(c.Remotes(), r)
	return c.saveAll()
}

func (c *FileConfig) Remove(name string) error {
	if _, err := c.Get(name); err != nil {
		return err
	}

	if c.Active() != nil && c.Active().Name == name {
		c.store.Active = ""
	}

	var newRemotes []Remote
	for _, r := range c.store.Remotes {
		if r.Name != name {
			newRemotes = append(newRemotes, r)
		}
	}
	c.store.Remotes = newRemotes
	return c.saveAll()
}

func (c *FileConfig) Get(name string) (Remote, error) {
	for _, r := range c.Remotes() {
		if r.Name == name {
			return r, nil
		}
	}
	return Remote{}, fmt.Errorf("remote '%s' does not exist", name)
}

func (c *FileConfig) SetActive(name string) error {
	r, err := c.Get(name)
	if err != nil {
		return err
	}
	c.store.Active = r.Name
	c.saveAll()
	return nil
}

func (c *FileConfig) Active() *Remote {
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
		return fmt.Errorf("Error parsing configuration file: %s", err.Error())
	}

	return nil
}

func (c *FileConfig) Remotes() []Remote {
	return c.store.Remotes
}

func (c *FileConfig) saveAll() error {
	b, err := json.MarshalIndent(c.store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Path, b, 0600)
}

func (r *Remote) DecodeToken() error {
	if r.Token == "" {
		return errors.New("Missing token")
	}
	bs, err := base64.StdEncoding.DecodeString(r.Token)
	if err != nil {
		return fmt.Errorf("There was a problem with your token: %s", err.Error())
	}

	data := strings.Split(string(bs), "|")
	if len(data) != 4 {
		return errors.New("There was a problem with your token: incorrect number of fields")
	}

	r.Endpoint = data[0]
	r.Username = data[1]
	r.Password = data[2]
	r.PrivateKey = data[3]

	return nil
}
