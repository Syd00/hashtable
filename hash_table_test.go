package main

import (
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
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

// Benchmark della nostra Custom Hash Table
func BenchmarkCustomHashTable(b *testing.B) {
	ht := NewHashTable[string, int](1024, StringHash)

	b.ResetTimer() // Azzera il tempo di creazione della tabella
	for i := 0; i < b.N; i++ {
		key := "key_" + string(rune(i%100))
		ht.Put(key, i)
		_, _ = ht.Get(key)
	}
}

// Benchmark della 'map' nativa di Go (per confronto!)
func BenchmarkNativeMap(b *testing.B) {
	m := make(map[string]int, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key_" + string(rune(i%100))
		m[key] = i
		_ = m[key]
	}
}
