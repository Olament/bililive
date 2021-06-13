package common

import (
	"go.mongodb.org/mongo-driver/bson"
	"sync/atomic"
)

// atomic String
type String struct {
	v atomic.Value
}

var zeroString string

func NewString(v string) *String {
	s := &String{}
	if v != zeroString {
		s.Store(v)
	}
	return s
}

func (s *String) Load() string {
	if v := s.v.Load(); v != nil {
		return v.(string)
	}
	return zeroString
}

func (s *String) Store(v string) {
	s.v.Store(v)
}

func (s *String) MarshalJSON() ([]byte, error) {
	return []byte(s.Load()), nil
}

func (s *String) MarshalBSON() ([]byte, error) {
	return bson.Marshal(s.Load())
}