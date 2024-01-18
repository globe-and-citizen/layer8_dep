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
	"globe-and-citizen/layer8/proxy/models"
	"time"
)

type (
	// Repository is a simple key-value store.
	Repository interface {

		// Get the salt from db using the username
		LoginUserPrecheck(username string) (string, error)

		// Get user from db by username
		GetUser(username string) (*models.User, error)

		// GetUserByID gets a user by ID.
		GetUserByID(id int64) (*models.User, error)

		// GetUserMetadata gets a user metadata by key.
		GetUserMetadata(userID int64, key string) (*models.UserMetadata, error)

		// Set a client for testing purposes
		SetClient(client *models.Client) error

		// Get a client by ID.
		GetClient(id string) (*models.Client, error)

		// SetTTL sets the value for the given key with a short TTL.
		SetTTL(key string, value []byte, ttl time.Duration) error

		// GetTTL gets the value for the given key which has a short TTL.
		GetTTL(key string) ([]byte, error)
	}

	// RepositoryFactory is a factory for creating repositories.
	RepositoryFactory interface {
		// Create creates a new repository.
		Create() (Repository, error)
	}

	// RepositoryFactoryFunc is an adapter that implements RepositoryFactory.
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
	RegisterRepositoryFactory("in_memory", RepositoryFactoryFunc(func() (Repository, error) {
		return TheInMemoryRepository, nil
	}))

	RegisterRepositoryFactory("postgres", RepositoryFactoryFunc(func() (Repository, error) {
		return NewPostgresRepository(), nil
	}))

}
