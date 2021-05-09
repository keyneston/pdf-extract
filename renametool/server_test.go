package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
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

func TestSafeFile(t *testing.T) {
	type testCase struct {
		root     string
		name     string
		expected error
	}

	testCases := []testCase{
		{"/tmp/tokens", "foo123.png", nil},
		{"/tmp/tokens", "../../foo123.png",
			fmt.Errorf(`"/foo123.png" does not seem to be a subpath of "/tmp/tokens"`)},
	}

	for _, c := range testCases {
		out := safeFile(c.root, c.name)
		assert.Equal(t, c.expected, out)
	}
}
