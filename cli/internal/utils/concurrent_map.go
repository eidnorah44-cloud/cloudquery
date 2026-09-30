package utils

import (
	"maps"
	"sync"
)

// ConcurrentMap is a thread-safe map of values.
type ConcurrentMap[T comparable, U any] struct {
	mu sync.RWMutex
	m  map[T]U
}

func NewConcurrentMap[T comparable, U any]() *ConcurrentMap[T, U] {
	return &ConcurrentMap[T, U]{
		m: make(map[T]U),
	}
}

func (s *ConcurrentMap[T, U]) Add(key T, value U) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// SetIfGreater updates the value for key under a single write lock if the key is not present
// or if isGreater returns true for (newValue, existingValue).
func (s *ConcurrentMap[T, U]) SetIfGreater(key T, value U, isGreater func(newVal, existingVal U) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[key]
	if !ok || isGreater(value, v) {
		s.m[key] = value
	}
}

func (s *ConcurrentMap[T, U]) Get(key T) (U, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *ConcurrentMap[T, U]) GetAll() map[T]U {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.m)
}
