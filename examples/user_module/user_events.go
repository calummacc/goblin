// goblin/examples/user_module/user_events.go
package user_module

import (
	"context"
)

// UserCreatedEvent represents a user creation event
type UserCreatedEvent struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Name returns the event name
func (e *UserCreatedEvent) Name() string {
	return "user.created"
}

// UserUpdatedEvent represents a user update event
type UserUpdatedEvent struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Name returns the event name
func (e *UserUpdatedEvent) Name() string {
	return "user.updated"
}

// UserDeletedEvent represents a user deletion event
type UserDeletedEvent struct {
	ID uint `json:"id"`
}

// Name returns the event name
func (e *UserDeletedEvent) Name() string {
	return "user.deleted"
}

// UserEventHandler handles user events
type UserEventHandler struct{}

// NewUserEventHandler creates a new user event handler
func NewUserEventHandler() *UserEventHandler {
	return &UserEventHandler{}
}

// HandleUserCreated handles user created events
func (h *UserEventHandler) HandleUserCreated(ctx context.Context, event *UserCreatedEvent) error {
	// Handle the event (e.g., send email, update cache, etc.)
	return nil
}

// HandleUserUpdated handles user updated events
func (h *UserEventHandler) HandleUserUpdated(ctx context.Context, event *UserUpdatedEvent) error {
	// Handle the event
	return nil
}

// HandleUserDeleted handles user deleted events
func (h *UserEventHandler) HandleUserDeleted(ctx context.Context, event *UserDeletedEvent) error {
	// Handle the event
	return nil
}
