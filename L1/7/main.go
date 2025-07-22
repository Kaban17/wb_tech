package main

import (
	"fmt"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu sync.Mutex
	mp map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		mp: make(map[K]V),
	}
}
func (s *SafeMap[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mp[key] = value
}
func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.mp[key]
	return value, ok
}
func (s *SafeMap[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.mp, key)
}
func main() {
	sm := NewSafeMap[string, int]()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			sm.Set(key, i)
		}(i)
	}
	wg.Wait()
	for i := range 10 {
		key := fmt.Sprintf("key%d", i)
		val, ok := sm.Get(key)
		if ok {
			fmt.Printf("%s has value %d\n", key, val)
		} else {
			fmt.Printf("%s not found\n", key)
		}
	}
}
