package users

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	dao "users-api/dao/users"
	domain "users-api/domain/users"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetAll() ([]dao.User, error)
	GetByID(id int64) (dao.User, error)
	GetByEmail(email string) (dao.User, error)
	Create(user dao.User) (int64, error)
	Update(user dao.User) error
	Delete(id int64) error
}

type Service struct {
	mainRepository      Repository
	cacheRepository     Repository
	memcachedRepository Repository
}

func NewService(mainRepository Repository, cacheRepository Repository, memcachedRepository Repository) Service {
	return Service{
		mainRepository:      mainRepository,
		cacheRepository:     cacheRepository,
		memcachedRepository: memcachedRepository,
		// tokenizer:           tokenizer,
	}
}

func (service Service) GetAll() ([]domain.User, error) {
	users, err := service.mainRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	result := make([]domain.User, 0)
	for _, user := range users {
		result = append(result, domain.User{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
			IsAdmin:  user.IsAdmin,
		})
	}

	return result, nil
}

func (service Service) GetByID(id int64) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByID(id)
	if err == nil {
		// Log when user is found in cache
		log.Printf("User found in cache for ID %d", id)
		return service.convertUser(user), nil
	}
	log.Printf("Cache miss for ID %d", id)

	// Check in memcached
	user, err = service.memcachedRepository.GetByID(id)
	if err == nil {
		// Log when user is found in memcached
		log.Printf("User found in memcached for ID %d", id)
		if _, err := service.cacheRepository.Create(user); err != nil {
			return domain.User{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.convertUser(user), nil
	}
	log.Printf("Memcached miss for ID %d", id)

	// Check in main repository
	user, err = service.mainRepository.GetByID(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by ID: %w", err)
	}

	// Save in cache and memcached
	log.Printf("User found in main repository for ID %d, saving to cache and memcached", id)
	if _, err := service.cacheRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.convertUser(user), nil
}

func (service Service) GetByEmail(email string) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByEmail(email)
	if err == nil {
		// Log when user is found in cache
		log.Printf("User found in cache for email %s", email)
		return service.convertUser(user), nil
	}
	log.Printf("Cache miss for email %s", email)

	// Check memcached
	user, err = service.memcachedRepository.GetByEmail(email)
	if err == nil {
		// Log when user is found in memcached
		log.Printf("User found in memcached for email %s", email)
		if _, err := service.cacheRepository.Create(user); err != nil {
			return domain.User{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.convertUser(user), nil
	}
	log.Printf("Memcached miss for email %s", email)

	// Check main repository
	user, err = service.mainRepository.GetByEmail(email)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by email: %w", err)
	}

	// Save in cache and memcached
	log.Printf("User found in main repository for email %s, saving to cache and memcached", email)
	if _, err := service.cacheRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.convertUser(user), nil
}

func (service Service) Create(user domain.User) (int64, error) {
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %w", err)
	}

	newUser := dao.User{
		Email:    user.Email,
		Password: string(hash),
		IsAdmin:  user.IsAdmin,
	}

	// Create in main repository
	id, err := service.mainRepository.Create(newUser)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	// Add to cache and memcached
	newUser.ID = id
	if _, err := service.cacheRepository.Create(newUser); err != nil {
		return 0, fmt.Errorf("error caching new user: %w", err)
	}
	if _, err := service.memcachedRepository.Create(newUser); err != nil {
		return 0, fmt.Errorf("error saving new user in memcached: %w", err)
	}

	return id, nil
}

func (service Service) Update(user domain.User) error {
	// Hash the password if provided
	var passwordHash string
	if user.Password != "" {
		passwordHash = Hash(user.Password) //Se hashea la nueva contra
	} else {
		existingUser, err := service.mainRepository.GetByID(user.ID)
		if err != nil {
			return fmt.Errorf("error retrieving existing user: %w", err)
		}
		passwordHash = existingUser.Password //En el caso que no haya contra nueva se sigue usando la vieja
	}

	// Update in main repository
	err := service.mainRepository.Update(dao.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	// Update in cache and memcached
	updatedUser := dao.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: passwordHash,
	}
	if err := service.cacheRepository.Update(updatedUser); err != nil {
		return fmt.Errorf("error updating user in cache: %w", err)
	}
	if err := service.memcachedRepository.Update(updatedUser); err != nil {
		return fmt.Errorf("error updating user in memcached: %w", err)
	}

	return nil
}

func (service Service) Delete(id int64) error {
	// Delete from main repository
	err := service.mainRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	// Delete from cache and memcached
	if err := service.cacheRepository.Delete(id); err != nil {
		return fmt.Errorf("error deleting user from cache: %w", err)
	}
	if err := service.memcachedRepository.Delete(id); err != nil {
		return fmt.Errorf("error deleting user from memcached: %w", err)
	}

	return nil
}

// func (service Service) Login(username string, password string) (domain.LoginResponse, error) {
// Hash the password

func (service Service) Login(body domain.Login_Request) (domain.LoginResponse, string, error) {
	// Try to get user from cache repository first
	user, err := service.cacheRepository.GetByEmail(body.Email)

	if err != nil {
		// If not found in cache, log and try to get user from memcached repository
		log.Printf("Cache miss for email %s, trying memcached repository", body.Email)
		user, err = service.memcachedRepository.GetByEmail(body.Email)
		if err != nil {
			// If not found in memcached, log and try to get user from the main repository (database)
			log.Printf("Memcached miss for email %s, trying main repository", body.Email)
			user, err = service.mainRepository.GetByEmail(body.Email)
			if err != nil {
				return domain.LoginResponse{}, "", fmt.Errorf("error getting user by email from main repository: %w", err)
			}

			// Save the found user in both cache and memcached repositories, log the caching operations
			log.Printf("User found in main repository, saving to cache and memcached for email %s", body.Email)
			if _, err := service.cacheRepository.Create(user); err != nil {
				return domain.LoginResponse{}, "", fmt.Errorf("error caching user in cache repository: %w", err)
			}
			if _, err := service.memcachedRepository.Create(user); err != nil {
				return domain.LoginResponse{}, "", fmt.Errorf("error caching user in memcached repository: %w", err)
			}

		} else {
			// Log when user is found in memcached, and save it to cache repository for future access
			log.Printf("User found in memcached repository for email %s", body.Email)
			if _, err := service.cacheRepository.Create(user); err != nil {
				return domain.LoginResponse{}, "", fmt.Errorf("error caching user in cache repository: %w", err)
			}
		}

	} else {
		// Log when user is found in cache
		log.Printf("User found in cache for email %s", body.Email)
	}

	// Password comparison
	// Verificar la contraseña usando bcrypt.CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return domain.LoginResponse{}, "", fmt.Errorf("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return domain.LoginResponse{}, "", errors.New("failed to create token")
	}

	// Prepare the response
	userDomain := domain.LoginResponse{
		ID:      strconv.FormatInt(user.ID, 10), // Aquí se realiza la conversión
		IsAdmin: user.IsAdmin,
		Email:   user.Email,
	}

	return userDomain, tokenString, nil
}

func (service Service) convertUser(user dao.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}

func Hash(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), 10)
	if err != nil {
		return "" //, fmt.Errorf("error hashing password: %w", err)
	}
	//hash = md5.Sum([]byte(input))
	return string(hash)
}
