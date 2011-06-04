package main

import (
	. "launchpad.net/gocheck"
	"testing"
	"bytes"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type S struct{}
var _ = Suite(&S{})

var gotestTests = map[string]*Stat{
	"Nothing here\nPASS\n": &Stat{num_errors: 0, matched_lines: []string{"PASS"}},
}

func (s *S) TestGoTestLines(c *C) {
	for input, expected := range gotestTests {
		string_reader := bytes.NewBufferString(input)
		actual := GoTestHandler(string_reader)
		c.Check(actual, Equals, expected)
	}
}