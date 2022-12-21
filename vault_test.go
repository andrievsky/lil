package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListParser1(t *testing.T) {
	output := []byte(`Keys
----
a/
b/
c
`)
	expected := []string{"a/", "b/", "c"}
	actual := make([]string, 0)
	err := parseList(output, func(localPath string) {
		actual = append(actual, localPath)
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, actual)
}

func TestIsFinal(t *testing.T) {
	assert.Equal(t, false, isFinal("/"))
	assert.Equal(t, false, isFinal("/secret/"))
	assert.Equal(t, false, isFinal("/secret/a/"))
	assert.Equal(t, true, isFinal("/a"))
	assert.Equal(t, true, isFinal("/secret/b"))
}
