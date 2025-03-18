// goblin/examples/user_module/user_events.go
package user_module

import "context"

// UserCreatedEvent is fired when a user is created
type UserCreatedEvent struct {
	user User
}

// Name returns the event name
func (e *UserCreatedEvent) Name() string {
	return "user.created"
}

// User returns the user
func (e *UserCreatedEvent) User() User {
	return e.user
}

// UserUpdatedEvent is fired when a user is updated
type UserUpdatedEvent struct {
	user User
}

// Name returns the event name
func (e *UserUpdatedEvent) Name() string {
	return "user.updated"
}

// User returns the user
func (e *UserUpdatedEvent) User() User {
	return e.user
}

// UserDeletedEvent is fired when a user is deleted
type UserDeletedEvent struct {
	user User
}

// Name returns the event name
func (e *UserDeletedEvent) Name() string {
	return "user.deleted"
}

// User returns the user
func (e *UserDeletedEvent) User() User {
	return e.user
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
