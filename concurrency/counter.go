package concurrency

import "sync"

type Counter interface {
	Inc()
	Value() int
}

type SimpleCounter struct {
	mu    sync.Mutex
	count int
}

func (s *SimpleCounter) Inc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
}

func (s *SimpleCounter) Value() int {
	return s.count
}

func NewCounter() Counter {
	return &SimpleCounter{}
}
