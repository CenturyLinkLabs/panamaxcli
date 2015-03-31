package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainOutput(t *testing.T) {
	assert.Equal(t, "Test", PlainOutput{"Test"}.ToPrettyOutput())
}

func TestListOutput(t *testing.T) {
	lo := ListOutput{Labels: []string{"ID"}}
	lo.AddRow(map[string]string{"ID": "10"})
	s := lo.ToPrettyOutput()

	assert.Contains(t, s, "ID")
	assert.Contains(t, s, "10")
}
