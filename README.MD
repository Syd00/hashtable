# HashTable

A thread-safe, generic Hash Table implementation in Go using separate chaining and dynamic resizing.

## 🚀 Features

- **Generics** (`[K comparable, V any]`) for full type-safety.
- **Thread-safe** operations via `sync.RWMutex` (concurrent reads with exclusive writes).
- **Dynamic Resizing** and collision resolution via Separate Chaining.
- **Idiomatic Iteration** with early-exit support (`Range`).

## 💡 Usage

### Basic Operations

```go
ht := hashtable.NewHashTable[string, int](16, hashtable.StringHash)

ht.Put("mario", 30)
val, ok := ht.Get("mario")
ht.Delete("mario")

// Iterate over key-value pairs (supports early exit)
ht.Range(func(key string, value int) bool {
    return true 
})

// Get a slice of all keys
keys := ht.Keys()
```

## ⚡ Benchmarks

Tested on **Intel Core i7-1185G7 @ 3.00GHz** running Windows `amd64`.

| Implementation | Operations (`b.N`) | Speed (`ns/op`) | Memory (`B/op`) | Allocs (`allocs/op`) |
| :--- | :--- | :--- | :--- | :--- |
| **Custom HashTable** | **19,074,261** | **61.26 ns/op** | **5 B/op** | **1 alloc/op** |
| Native Go `map` | 24,175,124 | 42.32 ns/op | 5 B/op | 1 alloc/op |

> The custom implementation runs within **1.4x** the speed of Go's native, highly optimized built-in map while providing explicit thread safety (`RWMutex`).