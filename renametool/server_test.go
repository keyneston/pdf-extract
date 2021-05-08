package main

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestCleanText(t *testing.T) {
	type testCase struct {
		in       []string
		expected []string
	}

	testCases := []testCase{
		{in: []string{"FOO"}, expected: []string{"foo"}},
		{in: []string{"FOO BAR"}, expected: []string{"foo-bar"}},
		{in: []string{"   blah "}, expected: []string{"blah"}},
		{in: []string{"© foo bar"}, expected: []string{"foo-bar"}},
		{in: []string{"  "}, expected: []string{}},
	}

	for _, c := range testCases {
		out := cleanText(c.in)

		if diff := deep.Equal(c.expected, out); diff != nil {
			t.Errorf("cleanText(%v) =\n%v", c.in, strings.Join(diff, "\n"))
		}
	}
}
