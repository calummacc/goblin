package user

import (
	"errors"
	"fmt"
)

type UserService struct {
	users []string
}

func NewUserService() *UserService {
	return &UserService{users: []string{"Alice", "Bob", "Charlie"}}
}

func (s *UserService) FindAll() []string {
	return s.users
}

func (s *UserService) FindOne(id string) (string, error) {
	for i, user := range s.users {
		if user == id {
			return s.users[i], nil
		}
	}
	return "", errors.New(fmt.Sprintf("user with id %s is not found", id))
}

func (s *UserService) Create(name string) error {
	for _, user := range s.users {
		if user == name {
			return errors.New(fmt.Sprintf("user %s is already exist", name))
		}
	}
	s.users = append(s.users, name)
	return nil
}

