package store

import (
	"sync"
	"time"
)

type item struct {
	value     string
	expiresAt time.Time // Represents the exact moment this key dies
}

type Store struct {
	mu   sync.RWMutex
	data map[string]item
}

func New() *Store {
	return &Store{
		data: make(map[string]item),
	}
}

func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = item{
		value: value,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return "", false
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(s.data, key)
		return "", false
	}

	return item.value, true

}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]
	return ok
}
