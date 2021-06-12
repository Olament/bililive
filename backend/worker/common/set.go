package common

type void struct{}

// regular set implementation
type Set struct {
	m map[int64]void
}

func NewSet() *Set {
	s := &Set{
		m: make(map[int64]void),
	}
	return s
}

func (s *Set) Add(element int64) {
	s.m[element] = void{}
}

func (s *Set) Len() int64 {
	return int64(len(s.m))
}