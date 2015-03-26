package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CenturyLinkLabs/panamaxcli/client"
)

func Search(p client.PanamaxClient, terms []string) (string, error) {
	s := strings.Join(terms, " ")
	if s == "" {
		return "", errors.New("Empty search")
	}
	results, err := p.SearchTemplates(s)
	if err != nil {
		return "", err
	}
	out := "Search Results\n"
	for _, result := range results {
		out += fmt.Sprintf("Id:%d, Name:%s\n", result.ID, result.Name)
	}
	return out, nil
}
