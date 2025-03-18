// goblin/events/event_bus.go
// Package events provides an event-driven architecture for the Goblin Framework.
// It implements a flexible event bus system that supports both synchronous and
// asynchronous event handling with retry mechanisms and error handling.
package events

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/fx"
)

// Event represents an event in the system.
// It provides a way to identify and carry data for system events.
// Events are used to decouple components and enable reactive programming.
type Event interface {
	// Name returns the unique identifier for this event type.
	// This identifier is used to match events to their subscribers.
	Name() string
}

// BaseEvent provides a basic implementation of the Event interface.
// It can be used as a base for custom event types or used directly
// for simple event scenarios.
type BaseEvent struct {
	// EventName uniquely identifies this event type
	EventName string
	// Payload contains the event data as an arbitrary value
	Payload interface{}
}

// Name returns the name of the event.
// This implementation simply returns the EventName field value.
//
// Returns:
//   - string: The event name
func (e BaseEvent) Name() string {
	return e.EventName
}

// NewEvent creates a new event with the given name and payload.
// This is a convenience factory function for creating BaseEvent instances.
//
// Example:
//
//	event := NewEvent("user.created", userObject)
//	eventBus.Publish(ctx, event)
//
// Parameters:
//   - name: The unique identifier for this event type
//   - payload: The data associated with this event
//
// Returns:
//   - Event: A new event instance
func NewEvent(name string, payload interface{}) Event {
	return BaseEvent{
		EventName: name,
		Payload:   payload,
	}
}

// EventHandlerMode determines how an event handler is executed.
// It supports both synchronous and asynchronous execution modes.
type EventHandlerMode int

const (
	// SyncMode executes handlers synchronously in the same goroutine.
	// This ensures ordered execution but may block the event publisher.
	// Use this mode when the exact order of execution matters or when
	// the publisher needs to know when all handlers have completed.
	SyncMode EventHandlerMode = iota

	// AsyncMode executes handlers in separate goroutines.
	// This provides better performance but may result in out-of-order execution.
	// Use this mode for handlers that don't depend on each other and when
	// the publisher doesn't need to wait for completion.
	AsyncMode
)

// Handler is a function that processes an event.
// It receives the context and event as parameters and may return an error.
// Handlers should be idempotent and prepared to handle duplicate events.
type Handler func(ctx context.Context, event Event) error

// EventHandlerConfig holds configuration options for an event handler.
// It allows customization of handler behavior including execution mode,
// retry policy, and error handling.
type EventHandlerConfig struct {
	// Mode determines if the handler should run synchronously or asynchronously
	Mode EventHandlerMode

	// MaxRetries specifies how many times to retry the handler on error (async only).
	// A value of 0 means no retries, and the handler will only be executed once.
	// This only applies to handlers in AsyncMode.
	MaxRetries int

	// ErrorHandler is called when a handler returns an error.
	// It receives the error, the event that caused it, and the name of the handler.
	// This allows for custom error handling strategies like logging or alerting.
	ErrorHandler func(err error, event Event, handlerName string)
}

// DefaultHandlerConfig returns the default handler configuration.
// It provides sensible defaults for most use cases.
// The default configuration uses SyncMode with no retries and
// a simple error handler that prints errors to stdout.
//
// Returns:
//   - EventHandlerConfig: The default handler configuration
func DefaultHandlerConfig() EventHandlerConfig {
	return EventHandlerConfig{
		Mode:       SyncMode,
		MaxRetries: 0,
		ErrorHandler: func(err error, event Event, handlerName string) {
			fmt.Printf("Error handling event %s in handler %s: %v\n", event.Name(), handlerName, err)
		},
	}
}

// registeredHandler holds a handler with its configuration and metadata.
// It is used internally by the EventBus to manage handler registration.
type registeredHandler struct {
	// handler is the function that will be called when an event occurs
	handler Handler
	// config contains the execution configuration for this handler
	config EventHandlerConfig
	// name is used for debugging and error reporting
	name string
}

// EventBus manages event publishing and subscription.
// It allows components to communicate without direct dependencies,
// enabling a more decoupled and maintainable architecture.
type EventBus struct {
	// handlers maps event names to their registered handlers
	handlers map[string][]registeredHandler
	// errorHandler is the default error handler for events
	errorHandler func(err error, event Event, handlerName string)
	// mu protects concurrent access to the handlers map
	mu sync.RWMutex
	// ctx is the context used for background operations
	ctx context.Context
	// cancel is the cancellation function for the context
	cancel context.CancelFunc
}

// EventBusParams defines the parameters for creating a new EventBus.
// It uses fx.In for dependency injection with Uber Fx.
type EventBusParams struct {
	fx.In

	// DefaultErrorHandler is used when handlers don't provide their own error handler.
	// It's optional and will use a default implementation if not provided.
	DefaultErrorHandler func(err error, event Event, handlerName string) `optional:"true"`
}

// NewEventBus creates a new event bus instance.
// It initializes the event bus with the provided parameters.
// If no DefaultErrorHandler is provided, a default one is used.
//
// Parameters:
//   - params: The parameters for creating the event bus
//
// Returns:
//   - *EventBus: A new event bus instance
func NewEventBus(params EventBusParams) *EventBus {
	errorHandler := params.DefaultErrorHandler
	if errorHandler == nil {
		errorHandler = func(err error, event Event, handlerName string) {
			fmt.Printf("Error handling event %s in handler %s: %v\n", event.Name(), handlerName, err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventBus{
		handlers:     make(map[string][]registeredHandler),
		errorHandler: errorHandler,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// NewEventBusModule provides an Fx module for the event bus.
// This function can be used to include the event bus in an Fx application.
//
// Example:
//
//	app := fx.New(
//	    events.NewEventBusModule(),
//	    // other modules...
//	)
//
// Returns:
//   - fx.Option: An Fx module option for the event bus
func NewEventBusModule() fx.Option {
	return fx.Options(
		fx.Provide(NewEventBus),
	)
}

// Subscribe registers a handler for an event type with default configuration.
// This is a convenience method that uses the default handler configuration.
//
// Example:
//
//	eventBus.Subscribe("user.created", func(ctx context.Context, event Event) error {
//	    user := event.Payload.(*User)
//	    // handle user creation...
//	    return nil
//	})
//
// Parameters:
//   - eventName: The name of the event to subscribe to
//   - handler: The function to call when the event occurs
func (b *EventBus) Subscribe(eventName string, handler Handler) {
	b.SubscribeWithConfig(eventName, handler, DefaultHandlerConfig(), "")
}

// SubscribeWithConfig registers a handler for an event type with custom configuration.
// This method allows full customization of how the handler is executed.
//
// Parameters:
//   - eventName: The name of the event to subscribe to
//   - handler: The function to call when the event occurs
//   - config: Configuration options for the handler
//   - handlerName: Optional name for the handler (used in error reporting)
func (b *EventBus) SubscribeWithConfig(eventName string, handler Handler, config EventHandlerConfig, handlerName string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Use a default error handler if none is provided
	if config.ErrorHandler == nil {
		config.ErrorHandler = b.errorHandler
	}

	// Use the function pointer address as the name if none is provided
	if handlerName == "" {
		handlerName = fmt.Sprintf("%p", handler)
	}

	// Add the handler to the registry
	b.handlers[eventName] = append(b.handlers[eventName], registeredHandler{
		handler: handler,
		config:  config,
		name:    handlerName,
	})
}

// Unsubscribe removes a handler for an event type.
// This method identifies the handler by function pointer comparison.
//
// Parameters:
//   - eventName: The name of the event to unsubscribe from
//   - handler: The handler function to remove
func (b *EventBus) Unsubscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, ok := b.handlers[eventName]
	if !ok {
		return
	}

	// Find and remove the handler
	for i, h := range handlers {
		// Compare function pointers
		if fmt.Sprintf("%p", h.handler) == fmt.Sprintf("%p", handler) {
			// Remove the handler by replacing it with the last one and truncating
			handlers[i] = handlers[len(handlers)-1]
			b.handlers[eventName] = handlers[:len(handlers)-1]
			break
		}
	}

	// Remove the event key if no handlers remain
	if len(b.handlers[eventName]) == 0 {
		delete(b.handlers, eventName)
	}
}

// Publish dispatches an event to all registered handlers.
// It executes handlers according to their configured execution mode.
// For SyncMode handlers, it waits for all handlers to complete.
// For AsyncMode handlers, it starts goroutines and returns immediately.
//
// Parameters:
//   - ctx: The context for the publish operation
//   - event: The event to publish
//
// Returns:
//   - []error: Any errors that occurred during synchronous handler execution
func (b *EventBus) Publish(ctx context.Context, event Event) []error {
	b.mu.RLock()
	handlers, ok := b.handlers[event.Name()]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	var errors []error

	// Execute handlers according to their mode
	for _, h := range handlers {
		switch h.config.Mode {
		case SyncMode:
			// Execute synchronously and collect errors
			if err := b.executeHandler(ctx, h, event); err != nil {
				errors = append(errors, err)

				// Call error handler if provided
				if h.config.ErrorHandler != nil {
					h.config.ErrorHandler(err, event, h.name)
				}
			}

		case AsyncMode:
			// Execute asynchronously with retries
			go func(h registeredHandler, event Event) {
				var err error
				// Try initial execution
				err = b.executeHandler(ctx, h, event)

				// Retry on failure if configured
				retries := 0
				for err != nil && retries < h.config.MaxRetries {
					retries++
					err = b.executeHandler(ctx, h, event)
				}

				// Call error handler if still failed after retries
				if err != nil && h.config.ErrorHandler != nil {
					h.config.ErrorHandler(err, event, h.name)
				}
			}(h, event)
		}
	}

	return errors
}

// executeHandler executes a single event handler with the given event.
// This is an internal helper method used by Publish.
//
// Parameters:
//   - ctx: The context for the handler execution
//   - h: The handler to execute
//   - event: The event to pass to the handler
//
// Returns:
//   - error: Any error returned by the handler
func (b *EventBus) executeHandler(ctx context.Context, h registeredHandler, event Event) error {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.ctx.Done():
		return b.ctx.Err()
	default:
		// Continue execution
	}

	// Execute the handler
	return h.handler(ctx, event)
}

// Shutdown gracefully shuts down the event bus.
// It cancels the internal context, which will stop any ongoing operations.
func (b *EventBus) Shutdown() {
	if b.cancel != nil {
		b.cancel()
	}
}

// EventListener provides a fluent interface for configuring event handlers.
// It allows for a more readable and chainable API for event subscription.
type EventListener struct {
	// eventName is the name of the event to listen for
	eventName string
	// handler is the function to call when the event occurs
	handler Handler
	// config contains the execution configuration for this handler
	config EventHandlerConfig
	// name is an optional identifier for the handler
	name string
}

// OnEvent creates a new event listener for the specified event.
// It returns an EventListener that can be further configured with method chaining.
//
// Example:
//
//	OnEvent("user.created", handleUserCreated).
//	    WithAsync().
//	    WithRetries(3).
//	    Register(eventBus)
//
// Parameters:
//   - eventName: The name of the event to listen for
//   - handler: The function to call when the event occurs
//
// Returns:
//   - *EventListener: A new event listener
func OnEvent(eventName string, handler Handler) *EventListener {
	return &EventListener{
		eventName: eventName,
		handler:   handler,
		config:    DefaultHandlerConfig(),
	}
}

// WithConfig sets the entire configuration for this listener.
// This method allows setting all configuration options at once.
//
// Parameters:
//   - config: The configuration to use
//
// Returns:
//   - *EventListener: The event listener for method chaining
func (l *EventListener) WithConfig(config EventHandlerConfig) *EventListener {
	l.config = config
	return l
}

// WithName sets a custom name for this listener.
// The name is used for debugging and error reporting.
//
// Parameters:
//   - name: The name to use for this listener
//
// Returns:
//   - *EventListener: The event listener for method chaining
func (l *EventListener) WithName(name string) *EventListener {
	l.name = name
	return l
}

// WithAsync configures this listener to execute asynchronously.
// The handler will be executed in a separate goroutine.
//
// Returns:
//   - *EventListener: The event listener for method chaining
func (l *EventListener) WithAsync() *EventListener {
	l.config.Mode = AsyncMode
	return l
}

// WithSync configures this listener to execute synchronously.
// The handler will be executed in the same goroutine as the publisher.
//
// Returns:
//   - *EventListener: The event listener for method chaining
func (l *EventListener) WithSync() *EventListener {
	l.config.Mode = SyncMode
	return l
}

// WithRetries sets the number of retries for this listener.
// This only applies to async handlers.
//
// Parameters:
//   - maxRetries: The maximum number of times to retry on failure
//
// Returns:
//   - *EventListener: The event listener for method chaining
func (l *EventListener) WithRetries(maxRetries int) *EventListener {
	l.config.MaxRetries = maxRetries
	return l
}

// WithErrorHandler sets a custom error handler for this listener.
// The error handler will be called if the handler returns an error.
//
// Parameters:
//   - handler: The error handler function
//
// Returns:
//   - *EventListener: The event listener for method chaining
func (l *EventListener) WithErrorHandler(handler func(err error, event Event, handlerName string)) *EventListener {
	l.config.ErrorHandler = handler
	return l
}

// Register registers this listener with an event bus.
// This finalizes the configuration and subscribes the handler.
//
// Parameters:
//   - bus: The event bus to register with
func (l *EventListener) Register(bus *EventBus) {
	bus.SubscribeWithConfig(l.eventName, l.handler, l.config, l.name)
}

// EventHandlerResult contains the result of registering an event handler.
// It provides a way to unsubscribe the handler later.
type EventHandlerResult struct {
	// Unsubscribe is a function that can be called to remove the handler.
	// This allows for dynamic subscription management.
	Unsubscribe func()
}

// AsEventHandler is a decorator for auto-registering event handlers.
// It can be used to automatically register a method as an event handler
// during dependency injection.
//
// Parameters:
//   - eventName: The name of the event to handle
//   - mode: The execution mode (sync or async)
//
// Returns:
//   - func(interface{}, string): A decorator function for the handler method
func AsEventHandler(eventName string, mode EventHandlerMode) func(interface{}, string) {
	return func(target interface{}, methodName string) {
		// The auto-registration will happen in the RegisterHandlers function
	}
}

// RegisterHandlers registers all event handlers in an instance.
// It uses reflection to find methods decorated with @AsEventHandler
// and registers them with the event bus.
//
// Parameters:
//   - instance: The instance containing handler methods
func (b *EventBus) RegisterHandlers(instance interface{}) {
	// This would normally use reflection to find and register handlers
	// But it's not implemented in this simplified version
}

// RegisterEventHandlers registers event handlers with the event bus during application startup.
// This function is meant to be used with dependency injection frameworks like Uber Fx.
//
// Parameters:
//   - lc: The lifecycle manager for cleanup
//   - bus: The event bus to register with
//   - handler: The instance containing handler methods
func RegisterEventHandlers(lc fx.Lifecycle, bus *EventBus, handler interface{}) {
	// Register handlers
	bus.RegisterHandlers(handler)

	// Register shutdown hook
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			bus.Shutdown()
			return nil
		},
	})
}
