package actions

import (
	"errors"
	"testing"

	"github.com/CenturyLinkLabs/panamaxcli/client"
	"github.com/stretchr/testify/assert"
)

type FakePanamax struct {
	SearchTerms     string
	ErrorForGetApps error
	ErrorForGetApp  error
	ErrorForSearch  error
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

func (p FakePanamax) GetApp(id int) (client.App, error) {
	if p.ErrorForGetApp != nil {
		return client.App{}, p.ErrorForGetApp
	}

	return client.App{ID: 1, Name: "Test"}, nil
}

func TestDescribeApp(t *testing.T) {
	output, err := DescribeApp(FakePanamax{}, 1)
	assert.Contains(t, output, "App Details")
	assert.Contains(t, output, "1")
	assert.Contains(t, output, "Test")
	assert.NoError(t, err)
}

func TestErroredDescribeApp(t *testing.T) {
	err := errors.New("Test Error")
	p := FakePanamax{ErrorForGetApp: err}
	output, err := DescribeApp(p, 1)
	assert.Equal(t, "", output)
	assert.EqualError(t, err, "Test Error")
}
