package user

import (
	"errors"
	"fmt"
)

type UserService struct {
	// Add dependencies (e.g., database) here
	users []string
}

func NewUserService() *UserService {
	return &UserService{users: []string{"Alice", "Bob", "Charlie"}}
}

func (s *UserService) FindAll() []string {
	return s.users
}

func (s *UserService) FindOne(id string) (string, error) {
	// Simulate finding a user by ID (replace with your actual database logic)
	i, err := findUserIndex(id, s.users)
	if err != nil {
		return "", err
	}
	return s.users[i], nil
}

func (s *UserService) Create(name string) error {
	//Simulate creating user (replace with database logic)
	for _, user := range s.users {
		if user == name {
			return errors.New(fmt.Sprintf("user %s is already exist", name))
		}
	}
	s.users = append(s.users, name)
	return nil
}

func findUserIndex(id string, users []string) (int, error) {
	for i, user := range users {
		if user == id {
			return i, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("user with id %s is not found", id))
}
