package gocron

import (
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	safeslice := NewSafeSlice(0)
	task := Task{
		job: func() {},
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		// add
		for i := 0; i < 1000000; i++ {
			safeslice.Add(&task)
		}

		wg.Done()
	}()
	go func() {
		// add
		for i := 0; i < 1000000; i++ {
			safeslice.Add(&task)
		}

		wg.Done()
	}()

	wg.Wait()

	if safeslice.Len() != 2000000 {
		t.Error("Expected 2000000, got", safeslice.Len())
	}
}
