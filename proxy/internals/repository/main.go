// Package repository provides a simple key-value store.
//
// The repository package provides a simple interface for creating and accessing the
// repository. For now, it only provides a memory implementation.
//
// Example:
//
//	repo := repository.MustCreateRepository("memory")
//
//	repo.Set("foo", []byte("bar"))
//	repo.Set("baz", []byte("quux"))
//
//	fmt.Println(string(repo.Get("foo")))
//	fmt.Println(string(repo.Get("baz")))
//
//	repo.Delete("foo")
//	repo.Delete("baz")
//
//	fmt.Println(string(repo.Get("foo"))) // nil
//	fmt.Println(string(repo.Get("baz"))) // nil
package repository

import (
	"globe-and-citizen/layer8/proxy/constants"
	"time"
)

type (
	// Repository is a simple key-value store.
	Repository interface {
		// Get returns the value for the given key, or nil if it does not exist.
		Get(key string) []byte
		// Pop gets the value for the given key and deletes it. Returns nil if the
		// key does not exist.
		Pop(key string) []byte
		// Set sets the value for the given key.
		Set(key string, value []byte) error
		// SetTTL sets the value for the given key with a short TTL.
		SetTTL(key string, value []byte, ttl time.Duration) error
		// Delete deletes the value for the given key.
		Delete(key string) error
		// Incr increments the value for the given key.
		// 	Note: It is good practice to begin the key with an underscore to avoid
		// 		  conflicts with other keys.
		Incr(key string) (int64, error)
		// Keys returns all keys matching the given pattern.
		Keys(pattern string) ([]string, error)
	}

	// RepositoryFactory is a factory for creating repositories.
	RepositoryFactory interface {
		// Create creates a new repository.
		Create() (Repository, error)
	}

	// RepositoryFactoryFunc is a function that implements RepositoryFactory.
	RepositoryFactoryFunc func() (Repository, error)

	// RepositoryFactoryRegistry is a registry of repository factories.
	RepositoryFactoryRegistry struct {
		factories map[string]RepositoryFactory
	}
)

// Create implements RepositoryFactory.
func (f RepositoryFactoryFunc) Create() (Repository, error) {
	return f()
}

// NewRepositoryFactoryRegistry creates a new repository factory registry.
func NewRepositoryFactoryRegistry() *RepositoryFactoryRegistry {
	return &RepositoryFactoryRegistry{
		factories: make(map[string]RepositoryFactory),
	}
}

// Register registers a repository factory.
func (r *RepositoryFactoryRegistry) Register(name string, factory RepositoryFactory) {
	r.factories[name] = factory
}

// Create creates a new repository.
func (r *RepositoryFactoryRegistry) Create(name string) (Repository, error) {
	factory, ok := r.factories[name]
	if !ok {
		return nil, constants.ErrUnknownRepositoryFactory
	}
	return factory.Create()
}

var (
	repositoryFactoryRegistry = NewRepositoryFactoryRegistry()
)

// RegisterRepositoryFactory registers a repository factory.
func RegisterRepositoryFactory(name string, factory RepositoryFactory) {
	repositoryFactoryRegistry.Register(name, factory)
}

// CreateRepository creates a new repository.
func CreateRepository(name string) (Repository, error) {
	return repositoryFactoryRegistry.Create(name)
}

// MustCreateRepository creates a new repository, panicking if an error occurs.
func MustCreateRepository(name string) Repository {
	repository, err := CreateRepository(name)
	if err != nil {
		panic(err)
	}
	return repository
}

func init() {
	// Register the memory repository factory
	RegisterRepositoryFactory("memory", RepositoryFactoryFunc(func() (Repository, error) {
		return NewMemoryRepository(), nil
	}))
}
