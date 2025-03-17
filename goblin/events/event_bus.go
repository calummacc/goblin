// goblin/events/event_bus.go
package events

import (
	"context"
	"sync"
)

// Event represents an event in the system
type Event interface {
	Name() string
}

// Handler is a function that handles an event
type Handler func(ctx context.Context, event Event) error

// EventBus manages event publishing and subscription
type EventBus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for an event
func (b *EventBus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.handlers[eventName]; !exists {
		b.handlers[eventName] = []Handler{}
	}
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish publishes an event to all registered handlers
func (b *EventBus) Publish(ctx context.Context, event Event) []error {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Name()]
	b.mu.RUnlock()

	if !exists {
		return nil
	}

	var errors []error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// EventListener creates an event listener
type EventListener struct {
	eventName string
	handler   Handler
}

// OnEvent creates a new event listener
func OnEvent(eventName string, handler Handler) *EventListener {
	return &EventListener{
		eventName: eventName,
		handler:   handler,
	}
}

// Register registers the listener with the event bus
func (l *EventListener) Register(bus *EventBus) {
	bus.Subscribe(l.eventName, l.handler)
}
