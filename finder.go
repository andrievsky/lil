package main

import (
	"time"
)

type Finder struct {
	timeout    time.Duration
	query      []rune
	lastUpdate time.Time
	now        func() time.Time
}

func NewFinder(timeout time.Duration, now func() time.Time) *Finder {
	return &Finder{
		timeout,
		make([]rune, 0),
		time.Now(),
		now,
	}
}

func (s *Finder) Update(key rune) {
	s.query = append(s.query, key)
	s.lastUpdate = time.Now()
}

func (s *Finder) Find(list []Path, key rune) int {
	time := s.now()

	if time.Sub(s.lastUpdate) > s.timeout {
		s.query = s.query[:0]
	}
	s.lastUpdate = time
	s.query = append(s.query, key)

	return indexOf(list, s.query)
}

func indexOf(list []Path, query []rune) int {
	queryString := string(query)
	maxIndex := -1
	maxScore := 0
	for i, item := range list {
		score := match(item.Label(), queryString)
		if score > maxScore {
			maxIndex = i
			maxScore = score
		}
	}
	return maxIndex
}

func match(source string, query string) int {
	bound := Min(len(source), len(query))
	for i := 0; i < bound; i++ {
		if source[i] != query[i] {
			return i
		}
	}
	return bound
}
