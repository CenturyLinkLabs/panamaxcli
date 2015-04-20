package main

import (
	"flag"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

var optionalFn = actionRequiresArgument("first", "optional:second")
var requiredFn = actionRequiresArgument("first", "second")

func contextWithFlags(args ...string) *cli.Context {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	flags.Parse(args)
	return cli.NewContext(nil, flags, nil)
}

func TestSuccessfulIncludingOptionalActionRequiresArgument(t *testing.T) {
	c := contextWithFlags("one", "two")
	assert.NoError(t, optionalFn(c))
}

func TestSuccessfulExcludingOptionalActionRequiresArgument(t *testing.T) {
	c := contextWithFlags("one")
	assert.NoError(t, optionalFn(c))
}

func TestErroredTooFewOptionalActionRequiresArgument(t *testing.T) {
	c := contextWithFlags()
	assert.EqualError(t, optionalFn(c), "This command requires the following arguments: first, optional:second")
}

func TestErroredTooManyActionRequiresArgument(t *testing.T) {
	c := contextWithFlags("one", "two", "three")
	assert.EqualError(t, requiredFn(c), "This command requires the following arguments: first, second")
}
