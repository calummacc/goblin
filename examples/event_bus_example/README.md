# Event Bus Example

This example demonstrates how to use the event bus in Goblin Framework to implement a publish-subscribe pattern for event-driven architecture.

## Overview

The example showcases:

1. **Event Publishing**: Services emitting domain events
2. **Event Subscription**: Multiple services consuming the same events
3. **Sync & Async Handling**: Both synchronous and asynchronous event processing
4. **Fx Integration**: Dependency injection for event handlers
5. **Retry Logic**: Automatic retrying of failed event handlers
6. **Error Handling**: Custom error handling for events

## Architecture

The example includes these components:

- **UserService**: Creates and updates users, publishes events
- **NotificationService**: Sends notifications to users (synchronous handlers)
- **AnalyticsService**: Records event statistics (asynchronous handlers with retry)
- **AuditLogger**: Logs audit trails for user actions (synchronous handlers)

## Events

Two types of events are demonstrated:

1. **UserCreatedEvent**: Emitted when a new user is created
2. **UserUpdatedEvent**: Emitted when a user's information is updated

## Running the Example

To run this example:

```bash
go run main.go
```

You'll see the output showing:
- Event emissions
- Synchronous handler execution
- Asynchronous handler execution
- A summary of all notifications, analytics, and audit logs

## Key Concepts

### Defining Events

Events are defined as structs that implement the `Event` interface:

```go
type UserCreatedEvent struct {
    UserID   string
    Username string
    Email    string
    Time     time.Time
}

func (e UserCreatedEvent) Name() string {
    return "user.created"
}
```

### Publishing Events

Events are published using the EventBus's Publish method:

```go
event := UserCreatedEvent{
    UserID:   id,
    Username: username,
    Email:    email,
    Time:     time.Now(),
}
errors := eventBus.Publish(ctx, event)
```

### Subscribing to Events

Event handlers are registered using the fluent API:

```go
// Synchronous handler
events.OnEvent("user.created", handleUserCreated).
    WithName("notification_user_created").
    Register(eventBus)

// Asynchronous handler with retry
events.OnEvent("user.created", handleEvent).
    WithName("analytics_user_created").
    WithAsync().
    WithRetries(3).
    Register(eventBus)
```

### Handling Events

Event handlers are functions that receive a context and an event:

```go
func (s *NotificationService) HandleUserCreated(ctx context.Context, event events.Event) error {
    e, ok := event.(UserCreatedEvent)
    if !ok {
        return fmt.Errorf("expected UserCreatedEvent, got %T", event)
    }
    
    // Handle the event...
    return nil
}
```

## Fx Integration

The event bus can be integrated with Uber's Fx dependency injection framework:

```go
app := fx.New(
    events.NewEventBusModule(),
    fx.Provide(
        NewUserService,
        NewNotificationService,
        // ... other services
    ),
    fx.Invoke(RegisterServices),
)
``` 