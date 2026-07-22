package main

import "fmt"

type Node[K comparable, V any] struct {
	key   K
	value V
	next  *Node[K, V]
}

type HashTable[K comparable, V any] struct {
	nodes    []*Node[K, V]
	hashFunc func(key K) int
	count    int
}

func (ht HashTable[K, V]) hash(key K) int {
	h := ht.hashFunc(key)

	if h < 0 {
		h = -h
	}

	return h % len(ht.nodes)
}

func NewHashTable[K comparable, V any](size int, hashFunc func(key K) int) *HashTable[K, V] {
	return &HashTable[K, V]{
		nodes:    make([]*Node[K, V], size),
		hashFunc: hashFunc,
		count:    0,
	}
}

func (ht *HashTable[K, V]) Put(key K, value V) {
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

func (ht *HashTable[K, V]) Get(key K) (V, bool) {
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

func (ht *HashTable[K, V]) Delete(key K) bool {
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
