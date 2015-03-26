package actions

import (
	"errors"
	"testing"

	"github.com/CenturyLinkLabs/panamaxcli/client"
	"github.com/stretchr/testify/assert"
)

type FakePanamax struct {
	ErrorForGetApps error
}

func (p FakePanamax) GetApps() ([]client.App, error) {
	if p.ErrorForGetApps != nil {
		return nil, p.ErrorForGetApps
	}

	a := client.App{ID: 123, Name: "Foo"}
	return []client.App{a}, nil
}

func TestListApps(t *testing.T) {
	output, err := ListApps(FakePanamax{})

	assert.NoError(t, err)
	assert.Contains(t, output, "Running Apps")
	assert.Contains(t, output, "123")
	assert.Contains(t, output, "Foo")
}

func TestErorredListApps(t *testing.T) {
	err := errors.New("GetApps Error")
	p := FakePanamax{ErrorForGetApps: err}
	output, err := ListApps(p)

	assert.Equal(t, "", output)
	assert.EqualError(t, err, "GetApps Error")
}
