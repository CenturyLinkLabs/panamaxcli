package actions

import (
	"fmt"
	"strconv"
)

type Output interface {
	ToPrettyOutput() string
}

type ListOutput struct {
	Labels []string
	Rows   []map[string]string
}

func (o *ListOutput) AddRow(m map[string]string) {
	o.Rows = append(o.Rows, m)
}

func (o *ListOutput) ToPrettyOutput() string {
	var output string
	for _, l := range o.Labels {
		output += l
	}
	output += "\n---------\n"
	for i, r := range o.Rows {
		output += strconv.Itoa(i) + ") "
		for _, l := range o.Labels {
			output += r[l]
		}
		output += "\n"
	}
	return output
}

type PlainOutput struct {
	Output string
}

func (o PlainOutput) ToPrettyOutput() string {
	return o.Output
}

type DetailOutput struct {
	Details map[string]string
}

func (o DetailOutput) ToPrettyOutput() string {
	var output string
	for k, v := range o.Details {
		output += fmt.Sprintf("%s: %s\n", k, v)
	}

	return output
}
