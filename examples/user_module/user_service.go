// goblin/examples/user_module/user_service.go
package user_module

import (
	"goblin/events"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserService handles user-related business logic
type UserService struct {
	// repo     *UserRepository
	eventBus *events.EventBus
}

// NewUserService creates a new user service
func NewUserService(eventBus *events.EventBus) *UserService {
	return &UserService{
		eventBus: eventBus,
	}
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() (interface{}, error) {
	return gin.H{"message": "GetAllUsers"}, nil
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id uint) (interface{}, error) {
	return gin.H{"message": "GetUserByID" + strconv.FormatUint(uint64(id), 10)}, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser() gin.H {
	s.eventBus.Publish(nil, &UserCreatedEvent{ID: 1, Username: "John Doe", Email: "john.doe@example.com"})
	return gin.H{"message": "CreateUser"}
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id uint) (interface{}, error) {
	// Validate user data
	s.eventBus.Publish(nil, &UserUpdatedEvent{ID: 1, Username: "John Doe", Email: "john.doe@example.com"})
	return gin.H{"message": "UpdateUser"}, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint) error {
	s.eventBus.Publish(nil, &UserDeletedEvent{ID: 1})
	return nil
}
