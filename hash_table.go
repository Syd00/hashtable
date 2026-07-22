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
			ht.Put(curr.key, curr.value)
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
	// Partiamo con una tabella PICCOLISSIMA: solo 2 bucket!
	ht := NewHashTable[string, int](2, StringHash)

	fmt.Println("Dimensione iniziale bucket:", len(ht.nodes)) // Stampava: 2

	// Inseriamo 5 elementi (supererà di molto il Load Factor 0.75!)
	ht.Put("mario", 30)
	ht.Put("luigi", 25)
	//ht.Put("peach", 28)
	//ht.Put("bowser", 40)
	ht.Put("yoshi", 12)

	fmt.Println("Elementi totali (count):", ht.count)
	fmt.Println("Nuova dimensione bucket dopo resize:", len(ht.nodes)) // Dovrebbe essere raddoppiata a 4 o 8!

	// Verifichiamo che tutti i dati siano ancora accessibili dopo il rehash
	val, ok := ht.Get("mario")
	fmt.Printf("Get('mario'): %d, trovato: %t\n", val, ok)

	val, ok = ht.Get("yoshi")
	fmt.Printf("Get('yoshi'): %d, trovato: %t\n", val, ok)
}
