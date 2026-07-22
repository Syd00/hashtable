package main

import (
	"sync"
	"testing"
)

func TestConcorre(t *testing.T) {
	ht := NewHashTable[int, int](10, func(k int) int { return k })
	var wg sync.WaitGroup

	numGoroutines := 10
	elementsPerRoutine := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		go func(goRoutineID int) {
			defer wg.Done()

			for j := 0; j < elementsPerRoutine; j++ {
				key := goRoutineID*elementsPerRoutine + j
				ht.Put(key, i)
			}
		}(i)
	}
	wg.Wait()

	if ht.count != 500 {
		t.Errorf("Expected 500, got %d", ht.count)
	}
}
