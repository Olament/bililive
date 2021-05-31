package common

import "time"

type Set struct {
	m    map[int64]*node
	head *node
	tail *node
	diff time.Duration
}

type node struct {
	element int64
	time    time.Time
	next    *node
}

func NewSet(diff time.Duration) *Set {
	s := Set{
		make(map[int64]*node),
		nil,
		nil,
		diff,
	}
	return &s
}

func (s *Set) Add(element int64) {
	n := &node{
		element: element,
		time:    time.Now(),
		next:    nil,
	}
	if s.head == nil { // first element
		s.head = n
		s.tail = n
	} else {
		s.tail.next = n
		s.tail = n
	}
	s.m[element] = n
}

func (s *Set) Len() int64 {
	for ok := s.evict(); ok; ok = s.evict() {
	}
	return int64(len(s.m))
}

func (s *Set) evict() bool {
	if s.head != nil && time.Now().Sub(s.head.time) > s.diff {
		if s.head == s.tail {
			s.tail = nil
		}
		if v, okay := s.m[s.head.element]; okay && v == s.head {
			delete(s.m, s.head.element)
		}
		s.head = s.head.next
		return true
	}
	return false
}
