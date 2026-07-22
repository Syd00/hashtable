// Implementation of a basic hash table
package main

import (
	"fmt"
	"sync"
)

// Node is a single element of a linked list
// used to manage collisions in the hash table
type Node[K comparable, V any] struct {
	key   K
	value V
	next  *Node[K, V]
}

// HashTable represents a generic hash table.
// Utilizes generics (K must be comparable) and V any
type HashTable[K comparable, V any] struct {
	mu       sync.RWMutex
	nodes    []*Node[K, V]
	hashFunc func(key K) int
	count    int
}

// NewHashTable inizialize an HashTable with initial size
// and an hash function personalized for key type
func NewHashTable[K comparable, V any](size int, hashFunc func(key K) int) *HashTable[K, V] {
	return &HashTable[K, V]{
		nodes:    make([]*Node[K, V], size),
		hashFunc: hashFunc,
		count:    0,
	}
}

// hash is an helper to calculate the index of the bucket.
// Applies hash function to the key
func (ht HashTable[K, V]) hash(key K) int {
	h := ht.hashFunc(key)

	if h < 0 {
		h = -h
	}

	return h % len(ht.nodes)
}

// Put add a key,value couple in the HashTable
// If load factor > 75% automatically resize
func (ht *HashTable[K, V]) Put(key K, value V) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	ht.putUnlocked(key, value)

	if ht.count*4 > len(ht.nodes)*3 {
		ht.resize()
	}
}

func (ht *HashTable[K, V]) putUnlocked(key K, value V) {
	bucket := ht.hash(key)

	curr := ht.nodes[bucket]
	for curr != nil {
		if curr.key == key {
			curr.value = value
			return
		}
		curr = curr.next
	}

	newNode := &Node[K, V]{
		key:   key,
		value: value,
		next:  ht.nodes[bucket],
	}

	ht.nodes[bucket] = newNode
	ht.count++
}

// Search for a value associated with the key
// if key exists return value, true
// if key doesn't exists return zero-value of type V and false
func (ht *HashTable[K, V]) Get(key K) (V, bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	var zero V
	bucket := ht.hash(key)
	curr := ht.nodes[bucket]
	for curr != nil {
		if curr.key == key {
			return curr.value, true
		}
		curr = curr.next
	}
	return zero, false
}

// Delete removes a key,value couple from the HashTable.
// returns true if element is found and deleted.
// returns false if element is not found.
func (ht *HashTable[K, V]) Delete(key K) bool {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	bucket := ht.hash(key)

	var prev *Node[K, V] = nil
	curr := ht.nodes[bucket]
	for curr != nil {
		if curr.key == key {
			if prev == nil {
				ht.nodes[bucket] = curr.next
			} else {
				prev.next = curr.next
			}
			ht.count--
			return true
		}
		prev = curr
		curr = curr.next
	}
	return false
}

// resize is a private method that double the bucket array size.
func (ht *HashTable[K, V]) resize() {
	newSize := len(ht.nodes) * 2
	oldNodes := ht.nodes
	ht.nodes = make([]*Node[K, V], newSize)
	ht.count = 0

	for _, head := range oldNodes {
		curr := head
		for curr != nil {
			ht.putUnlocked(curr.key, curr.value)
			curr = curr.next
		}
	}
}

// Keys returns a slice containing all the keys in the hashtable
func (ht *HashTable[K, V]) Keys() []K {
	ht.mu.RLock()
	defer ht.mu.RUnlock()

	keys := make([]K, 0, ht.count)
	for _, head := range ht.nodes {
		curr := head
		for curr != nil {
			keys = append(keys, curr.key)
			curr = curr.next
		}
	}
	return keys
}

func (ht *HashTable[K, V]) Range(f func(key K, value V) bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()

	for _, head := range ht.nodes {
		curr := head
		for curr != nil {
			keepGoing := f(curr.key, curr.value)
			if !keepGoing {
				return
			}
			curr = curr.next
		}
	}
}

// StringHash calculates the hash of a string using DJB2 algorithm.
func StringHash(s string) int {
	hash := 5381
	for i := 0; i < len(s); i++ {
		hash = ((hash << 5) + hash) + int(s[i]) // hash * 33 + char
	}
	return hash
}

func main() {
	ht := NewHashTable[string, int](4, StringHash)

	ht.Put("mario", 30)
	ht.Put("luigi", 25)
	ht.Put("peach", 28)

	// Estraiamo tutte le chiavi
	chiavi := ht.Keys()
	fmt.Println("Chiavi estratte:", chiavi)

	// Verifica della lunghezza
	if len(chiavi) != 3 {
		fmt.Printf("ERRORE: Ci si aspettavano 3 chiavi, trovate %d\n", len(chiavi))
	}

	// Creiamo un set temporaneo per controllare la presenza di tutte le chiavi
	attese := map[string]bool{"mario": true, "luigi": true, "peach": true}
	for _, k := range chiavi {
		delete(attese, k) // Rimuoviamo la chiave trovata
	}

	// Se la mappa 'attese' è vuota, significa che tutte le chiavi erano presenti!
	if len(attese) == 0 {
		fmt.Println("Test Keys(): SUCCESS!")
	} else {
		fmt.Println("Test Keys(): FALLITO, mancavano delle chiavi!")
	}

	fmt.Println("\n=== TEST RANGE ===")

	// CASO 1: Scorrimento completo e somma dei valori
	sommaEta := 0
	ht.Range(func(key string, value int) bool {
		fmt.Printf("Elemento: %s -> %d\n", key, value)
		sommaEta += value
		return true // Continua l'iterazione
	})

	// 30 (mario) + 25 (luigi) + 28 (peach) = 83
	fmt.Println("Somma età totale:", sommaEta)

	// CASO 2: Interruzione anticipata (si ferma al primo elemento)
	conteggioVisitati := 0
	ht.Range(func(key string, value int) bool {
		conteggioVisitati++
		return false // Interrompe subito l'iterazione alla prima chiave!
	})

	fmt.Println("Elementi visitati prima dello stop:", conteggioVisitati)
	// Stampava 1!
}
