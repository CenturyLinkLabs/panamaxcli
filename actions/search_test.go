package actions

import (
	"errors"
	"testing"

	"github.com/CenturyLinkLabs/panamaxcli/client"
	"github.com/stretchr/testify/assert"
)

func (p FakePanamax) SearchTemplates(terms string) ([]client.Template, error) {
	p.SearchTerms = terms
	if p.ErrorForSearch != nil {
		return nil, p.ErrorForSearch
	}
	return []client.Template{{ID: 1, Name: "Wordpress with MySQL"}}, nil
}

func TestSearch(t *testing.T) {
	p := FakePanamax{}
	output, err := Search(&p, []string{"wordpress", "mysql"})

	//assert.Equal(t, "wordpress mysql", p.SearchTerms)
	assert.NoError(t, err)
	assert.Contains(t, output, "Results")
	assert.Contains(t, output, "Wordpress with MySQL")
}

func TestErroredNoTermsSearch(t *testing.T) {
	p := FakePanamax{}
	output, err := Search(&p, []string{})
	//assert.Equal(t, "", p.SearchTerms)
	assert.EqualError(t, err, "Empty search")
	assert.Equal(t, "", output)
}

func TestErroredSearchErrorSearch(t *testing.T) {
	err := errors.New("test error")
	p := FakePanamax{ErrorForSearch: err}
	output, err := Search(&p, []string{"test"})
	assert.EqualError(t, err, "test error")
	assert.Equal(t, "", output)
}
