// goblin/examples/user_module/user_service.go
package user_module

import (
	"errors"
	"goblin/database"
	"goblin/events"
	"strconv"
)

// UserService handles user-related business logic
type UserService struct {
	repo      *UserRepository
	eventBus  *events.EventBus
	txManager *database.TransactionManager
}

// NewUserService creates a new user service
func NewUserService(repo *UserRepository, eventBus *events.EventBus, txManager *database.TransactionManager) *UserService {
	return &UserService{
		repo:      repo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// GetUsers returns all users
func (s *UserService) GetUsers() ([]User, error) {
	return s.repo.FindAll()
}

// GetUser returns a user by ID
func (s *UserService) GetUser(id string) (User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return User{}, errors.New("invalid user ID")
	}
	return s.repo.FindByID(uint(idInt))
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user User) (User, error) {
	// Use transaction
	var createdUser User
	err := s.txManager.Transaction(nil, func(tx interface{}) error {
		var err error
		createdUser, err = s.repo.Create(user)
		if err != nil {
			return err
		}

		// Publish event
		s.eventBus.Publish(nil, &UserCreatedEvent{user: createdUser})
		return nil
	})
	// goblin/examples/user_module/user_service.go (continued)
	return createdUser, err
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id string, user User) (User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return User{}, errors.New("invalid user ID")
	}

	// Use transaction
	var updatedUser User
	err = s.txManager.Transaction(nil, func(tx interface{}) error {
		var err error
		updatedUser, err = s.repo.Update(uint(idInt), user)
		if err != nil {
			return err
		}

		// Publish event
		s.eventBus.Publish(nil, &UserUpdatedEvent{user: updatedUser})
		return nil
	})

	return updatedUser, err
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Use transaction
	return s.txManager.Transaction(nil, func(tx interface{}) error {
		user, err := s.repo.FindByID(uint(idInt))
		if err != nil {
			return err
		}

		if err := s.repo.Delete(uint(idInt)); err != nil {
			return err
		}

		// Publish event
		s.eventBus.Publish(nil, &UserDeletedEvent{user: user})
		return nil
	})
}
