package repository

import (
	"regexp"
	"strconv"
	"time"
)

type MemoryRepository struct {
	storage map[string][]byte
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{make(map[string][]byte)}
}

// Get returns the value of the key
func (r *MemoryRepository) Get(key string) []byte {
	value, ok := r.storage[key]
	if !ok {
		return nil
	}
	return value
}

// Pop returns the value of the key and deletes it
func (r *MemoryRepository) Pop(key string) []byte {
	value, ok := r.storage[key]
	if !ok {
		return nil
	}
	delete(r.storage, key)
	return value
}

// Set sets the key to hold the value
func (r *MemoryRepository) Set(key string, value []byte) error {
	r.storage[key] = value
	return nil
}

// SetTTL sets the key to hold the value for a limited time
func (r *MemoryRepository) SetTTL(key string, value []byte, ttl time.Duration) error {
	r.storage[key] = value
	go func() {
		time.Sleep(ttl)
		delete(r.storage, key)
	}()
	return nil
}

// Delete deletes the key from the storage
func (r *MemoryRepository) Delete(key string) error {
	delete(r.storage, key)
	return nil
}

// Incr increments the key by 1 and returns the new value.
//
//	Note:
//		It is good practice to begin the key with an underscore
//		to avoid conflicts with other keys.
func (r *MemoryRepository) Incr(key string) (int64, error) {
	current, ok := r.storage[key]
	if !ok {
		current = []byte("0")
	}
	index, err := strconv.ParseInt(string(current), 10, 64)
	if err != nil {
		return 0, err
	}
	index++
	r.storage[key] = []byte(strconv.FormatInt(index, 10))
	return index, nil
}

// Keys returns all the keys in the storage using the pattern
func (r *MemoryRepository) Keys(pattern string) ([]string, error) {
	keys := make([]string, 0)
	for key := range r.storage {
		matched, err := regexp.MatchString(pattern, key)
		if err != nil {
			return nil, err
		}
		if matched {
			keys = append(keys, key)
		}
	}
	return keys, nil
}
