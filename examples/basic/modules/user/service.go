package user

import (
	"time"
)

type Service interface {
	GetAllUsers() ([]User, error)
	GetUserByID(id uint) (*User, error)
	CreateUser(username, email string) (*User, error)
	UpdateUser(id uint, username, email string) (*User, error)
	DeleteUser(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllUsers() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) GetUserByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) CreateUser(username, email string) (*User, error) {
	now := time.Now()
	user := &User{
		ID:        uint(now.UnixNano()),
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateUser(id uint, username, email string) (*User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
