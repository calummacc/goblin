// goblin/examples/user_module/user_module.go
package user_module

import (
	"context"
	"goblin/core"
	"goblin/events"

	"go.uber.org/fx"
)

// Module is the user module
func Module() core.GoblinModule {
	return core.NewModule("UserModule", fx.Options(
		// Providers
		fx.Provide(NewUserRepository),
		fx.Provide(NewUserService),
		fx.Provide(NewUserController),
		fx.Provide(NewUserEventHandler),

		// Event handlers
		fx.Invoke(RegisterEventHandlers),
	))
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
