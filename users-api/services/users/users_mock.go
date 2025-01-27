package users

import (
	"errors"
	"sync"
	dao "users-api/dao/users"
)

var (
	mockUsers = make(map[int64]dao.User)
	mockMutex sync.RWMutex
	lastID    int64 = 0
)

type Mock struct{}

func NewMock() Repository {
	return &Mock{}
}

func (m *Mock) GetAll() ([]dao.User, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	var users []dao.User
	for _, user := range mockUsers {
		users = append(users, user)
	}
	return users, nil
}

func (m *Mock) GetByID(id int64) (dao.User, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	if user, exists := mockUsers[id]; exists {
		return user, nil
	}
	return dao.User{}, errors.New("user not found")
}

func (m *Mock) GetByEmail(email string) (dao.User, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	for _, user := range mockUsers {
		if user.Email == email {
			return user, nil
		}
	}
	return dao.User{}, errors.New("user not found")
}

func (m *Mock) Create(user dao.User) (int64, error) {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	if user.ID > 0 {
		mockUsers[user.ID] = user
		if user.ID > lastID {
			lastID = user.ID
		}
		return user.ID, nil
	}

	lastID++
	user.ID = lastID
	mockUsers[lastID] = user
	return lastID, nil
}

func (m *Mock) Update(user dao.User) error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	if _, exists := mockUsers[user.ID]; !exists {
		return errors.New("user not found")
	}

	mockUsers[user.ID] = user
	return nil
}

func (m *Mock) Delete(id int64) error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	delete(mockUsers, id)
	return nil
}
