package gocron

import "sync"

type SafeSlice struct {
	slice  []*Task
	rwlock sync.RWMutex
}

func NewSafeSlice(cap uint) *SafeSlice {
	return &SafeSlice{
		slice: make([]*Task, cap),
	}
}

func (s *SafeSlice) Add(task *Task) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	s.slice = append(s.slice, task)
}
func (s *SafeSlice) Get(index int) *Task {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	return s.slice[index]
}

func (s *SafeSlice) Remove(index int) bool {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	if index >= len(s.slice) {
		return false
	}

	s.slice = append(s.slice[:index], s.slice[index+1:]...)
	return true

}

func (s *SafeSlice) IndexOf(task *Task) int {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	for i, t := range s.slice {
		if t == task {
			return i
		}
	}

	return -1
}

func (s *SafeSlice) Len() int {
	return len(s.slice)
}
