// goblin/examples/user_module/user_module.go
package user_module

import (
	"context"
	"goblin/core"
	"goblin/events"

	"github.com/gin-gonic/gin"
)

// UserModule represents the user feature module.
// It provides user-related functionality including user management,
// authentication, and authorization.
type UserModule struct {
	*core.BaseModule
}

// NewUserModule creates and initializes a new user module with its providers,
// controllers, and exports. It sets up:
// - User service for business logic
// - User controller for HTTP endpoints
// - User repository for data access
func NewUserModule() *UserModule {
	module := &UserModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			// Provide user service
			NewUserService,
			// Provide user controller
			NewUserController,
		},
		Controllers: []interface{}{
			// Register user routes
			func(engine *gin.Engine, controller *UserController) {
				controller.RegisterRoutes(engine)
			},
		},
		// Export user service for use in other modules
		Exports: []interface{}{
			NewUserService,
		},
	})
	return module
}

// OnModuleInit initializes the user module and its resources.
// This method is called when the application starts up.
//
// Parameters:
//   - ctx: The context for the initialization process
//
// Returns:
//   - error: Any error that occurred during initialization
func (m *UserModule) OnModuleInit(ctx context.Context) error {
	// Initialize user module resources (e.g., database tables)
	return nil
}

// OnModuleDestroy performs cleanup operations for the user module.
// This method is called when the application is shutting down.
//
// Parameters:
//   - ctx: The context for the cleanup process
//
// Returns:
//   - error: Any error that occurred during cleanup
func (m *UserModule) OnModuleDestroy(ctx context.Context) error {
	// Cleanup user module resources
	return nil
}

// RegisterEventHandlers registers event handlers
func RegisterEventHandlers(handler *UserEventHandler, eventBus *events.EventBus) {
	eventBus.Subscribe("user.created", func(ctx context.Context, event events.Event) error {
		if e, ok := event.(*UserCreatedEvent); ok {
			return handler.HandleUserCreated(ctx, e)
		}
		return nil
	})

	eventBus.Subscribe("user.updated", func(ctx context.Context, event events.Event) error {
		if e, ok := event.(*UserUpdatedEvent); ok {
			return handler.HandleUserUpdated(ctx, e)
		}
		return nil
	})

	eventBus.Subscribe("user.deleted", func(ctx context.Context, event events.Event) error {
		if e, ok := event.(*UserDeletedEvent); ok {
			return handler.HandleUserDeleted(ctx, e)
		}
		return nil
	})
}
