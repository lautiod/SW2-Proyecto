package users

import (
	"fmt"
	"time"
	"users-api/dao/users"

	"github.com/karlseguin/ccache"
)

type CacheConfig struct {
	TTL time.Duration // Cache expiration time
}

type Cache struct {
	client *ccache.Cache //Cliente que interactua con la cache
	ttl    time.Duration //tiempo de vida de los objetos en la cache
}

func NewCache(config CacheConfig) Cache {
	// Initialize ccache with default settings
	cache := ccache.New(ccache.Configure())
	return Cache{
		client: cache,
		ttl:    config.TTL,
	}
}

func (repository Cache) GetAll() ([]users.User, error) {
	// Since it's not typical to cache all users in one request, you might skip caching here
	// Alternatively, you can cache a summary list if needed
	return nil, fmt.Errorf("GetAll not implemented in cache")
}

func (repository Cache) GetByID(id int64) (users.User, error) {
	// Convert ID to string for cache key
	idKey := fmt.Sprintf("user:id:%d", id)

	// Try to get from cache
	item := repository.client.Get(idKey)
	if item != nil && !item.Expired() {
		// Return cached value
		user, ok := item.Value().(users.User)
		if !ok {
			return users.User{}, fmt.Errorf("failed to cast cached value to user")
		}
		return user, nil
	}

	// If not found, return cache miss error
	return users.User{}, fmt.Errorf("cache miss for user ID %d", id)
}

func (repository Cache) GetByEmail(email string) (users.User, error) {
	// Use username as cache key
	emailKey := fmt.Sprintf("user:email:%s", email)

	// Try to get from cache
	item := repository.client.Get(emailKey)
	if item != nil && !item.Expired() {
		// Return cached value
		user, ok := item.Value().(users.User)
		if !ok {
			return users.User{}, fmt.Errorf("failed to cast cached value to email")
		}
		return user, nil
	}

	// If not found, return cache miss error
	return users.User{}, fmt.Errorf("cache miss for email %s", email)
}

func (repository Cache) Create(user users.User) (int64, error) {
	// Cache user by ID and by email after creation
	idKey := fmt.Sprintf("user:id:%d", user.ID)
	userKey := fmt.Sprintf("user:email:%s", user.Email)

	// Set user in cache
	repository.client.Set(idKey, user, repository.ttl)
	repository.client.Set(userKey, user, repository.ttl)

	// Return the user ID as if it was created successfully
	return user.ID, nil
}

func (repository Cache) Update(user users.User) error {
	// Update both the ID and email keys in cache
	idKey := fmt.Sprintf("user:id:%d", user.ID)
	emailKey := fmt.Sprintf("user:email:%s", user.Email)

	// Set the updated user in cache
	repository.client.Set(idKey, user, repository.ttl)
	repository.client.Set(emailKey, user, repository.ttl)

	return nil
}

func (repository Cache) Delete(id int64) error {
	// Delete user by ID and email from cache
	idKey := fmt.Sprintf("user:id:%d", id)

	// Try to get user by ID to retrieve email
	user, err := repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("error retrieving user by ID for deletion: %w", err)
	}

	// Delete by ID
	repository.client.Delete(idKey)

	// Delete by username
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	repository.client.Delete(emailKey)

	return nil
}
