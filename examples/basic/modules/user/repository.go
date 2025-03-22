package user

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type Repository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type repository struct {
	// In a real application, you would have your database connection here
	users map[uint]*User
}

func NewRepository() Repository {
	return &repository{
		users: make(map[uint]*User),
	}
}

func (r *repository) FindAll() ([]User, error) {
	users := make([]User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, *user)
	}
	return users, nil
}

func (r *repository) FindByID(id uint) (*User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (r *repository) Create(user *User) error {
	if _, exists := r.users[user.ID]; exists {
		return ErrUserExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *repository) Update(user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *repository) Delete(id uint) error {
	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}
