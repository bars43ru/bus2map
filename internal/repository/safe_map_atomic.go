package repository

import (
	"sync/atomic"
)

type SafeMapAtomic[Key comparable, Value any] struct {
	data atomic.Value
}

func (s *SafeMapAtomic[Key, Value]) Get(key Key) (Value, bool) {
	m := s.data.Load().(map[Key]Value)
	value, exists := m[key]
	return value, exists
}

func (s *SafeMapAtomic[Key, Value]) Replace(data map[Key]Value) {
	s.data.Store(data)
}

func NewSafeMapAtomic[Key comparable, Value any]() SafeMapAtomic[Key, Value] {
	var sm SafeMapAtomic[Key, Value]
	sm.Replace(make(map[Key]Value))
	return sm
}
