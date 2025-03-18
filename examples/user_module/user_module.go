// goblin/examples/user_module/user_module.go
package user_module

import (
	"context"
	"goblin/core"
	"goblin/events"
)

// UserModule represents the user module
type UserModule struct {
	*core.BaseModule
}

// NewUserModule creates a new user module
func NewUserModule() *UserModule {
	module := &UserModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			NewUserRepository,
			NewUserService,
			NewUserController,
			NewUserEventHandler,
		},
		Controllers: []interface{}{
			RegisterEventHandlers,
		},
		// Export UserService để các module khác có thể sử dụng
		Exports: []interface{}{
			NewUserService,
		},
	})
	return module
}

// OnModuleInit is called when the module is initialized
func (m *UserModule) OnModuleInit(ctx context.Context) error {
	// Khởi tạo module resources (ví dụ: kết nối database)
	return nil
}

// OnModuleDestroy is called when the module is destroyed
func (m *UserModule) OnModuleDestroy(ctx context.Context) error {
	// Cleanup module resources
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
