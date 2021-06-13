package common

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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
	return json.Marshal(s.Load())
}

func (s *String) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(s.Load())
}