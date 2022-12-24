package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFinder_FindDistinct(t *testing.T) {
	finder := NewFinder(0, oneSecondDelayTimeProvider())
	list := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	assert.Equal(t, 0, finder.Find(list, 'a'))
	assert.Equal(t, 1, finder.Find(list, 'b'))
	assert.Equal(t, 2, finder.Find(list, 'c'))
	assert.Equal(t, -1, finder.Find(list, 'd'))
}

func TestFinder_FindSeq(t *testing.T) {
	finder := NewFinder(2*time.Second, oneSecondDelayTimeProvider())
	list := []Path{
		buildPath("abc"),
		buildPath("aab"),
		buildPath("aac"),
		buildPath("acdc"),
	}
	assert.Equal(t, 0, finder.Find(list, 'a'))
	assert.Equal(t, 1, finder.Find(list, 'a'))
	assert.Equal(t, 2, finder.Find(list, 'c'))
	assert.Equal(t, 2, finder.Find(list, 'd'))
}

func oneSecondDelayTimeProvider() func() time.Time {
	startTime := time.Now()
	step := time.Second
	i := 0
	return func() time.Time {
		i++
		return startTime.Add(time.Duration(i) * step)
	}
}
