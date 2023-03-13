package syncMap

import "sync"

// SyncMap is like a Go map[K]V but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// Loads, stores, and deletes run in amortized constant time.
//
// The SyncMap type is specialized. Most code should use a plain Go map instead,
// with separate locking or coordination, for better type safety and to make it
// easier to maintain other invariants along with the map content.
//
// The SyncMap type is optimized for two common use cases: (1) when the entry for a given
// key is only ever written once but read many times, as in caches that only grow,
// or (2) when multiple goroutines read, write, and overwrite entries for disjoint
// sets of keys. In these two cases, use of a SyncMap may significantly reduce lock
// contention compared to a Go map paired with a separate Mutex or RWMutex.
//
// The zero SyncMap is empty and ready for use. A SyncMap must not be copied after first use.
//
// In the terminology of the Go memory model, SyncMap arranges that a write operation
// “synchronizes before” any read operation that observes the effect of the write, where
// read and write operations are defined as follows.
// Load, LoadAndDelete, LoadOrStore are read operations;
// Delete, LoadAndDelete, and Store are write operations;
// and LoadOrStore is a write operation when it returns loaded set to false.
type SyncMap[K comparable, V any] struct {
	sync.Map
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *SyncMap[K, V]) Load(key K) (val V, ok bool) {
	res, ok := m.Map.Load(key)
	if res != nil {
		val = res.(V)
	}
	return val, ok
}

// Store sets the value for a key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

// Delete deletes the value for a key.
func (m *SyncMap[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (val V, ok bool) {
	res, ok := m.Map.LoadAndDelete(key)
	if res != nil {
		val = res.(V)
	}
	return val, ok
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	act, loaded := m.Map.LoadOrStore(key, value)
	if act != nil {
		actual = act.(V)
	}
	return actual, loaded
}
