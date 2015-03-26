package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakePanamax struct{}

func (p FakePanamax) GetApps() []App {
	a := App{ID: 123, Name: "Foo"}
	return []App{a}
}

func TestListApps(t *testing.T) {
	output, err := ListApps(FakePanamax{})

	assert.NoError(t, err)
	assert.Contains(t, output, "Running Apps")
	assert.Contains(t, output, "123")
	assert.Contains(t, output, "Foo")
}
