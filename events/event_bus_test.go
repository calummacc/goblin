package events

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEvent is a simple event for testing
type TestEvent struct {
	data string
}

func (e TestEvent) Name() string {
	return "test.event"
}

// TestEventWithPayload is an event with additional data
type TestEventWithPayload struct {
	ID      string
	Message string
	Time    time.Time
}

func (e TestEventWithPayload) Name() string {
	return "test.event.with.payload"
}

// TestSyncEventHandling tests that synchronous event handlers are executed correctly
func TestSyncEventHandling(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Create a channel to signal handler execution
	handlerExecuted := make(chan bool, 1)
	handlerError := make(chan error, 1)

	// Register a handler
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		// Verify the event
		e, ok := event.(TestEvent)
		if !ok {
			handlerError <- errors.New("event type mismatch")
			return errors.New("event type mismatch")
		}

		// Verify the data
		if e.data != "test data" {
			handlerError <- errors.New("data mismatch")
			return errors.New("data mismatch")
		}

		// Signal execution
		handlerExecuted <- true
		return nil
	})

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	errors := bus.Publish(ctx, event)

	// Verify no errors
	assert.Empty(t, errors)

	// Verify handler execution
	select {
	case <-handlerExecuted:
		// Success
	case err := <-handlerError:
		t.Fatalf("Handler error: %v", err)
	case <-time.After(time.Second):
		t.Fatal("Handler not executed within timeout")
	}
}

// TestAsyncEventHandling tests that asynchronous event handlers are executed correctly
func TestAsyncEventHandling(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Create a channel to signal handler execution
	handlerExecuted := make(chan bool, 1)

	// Create a WaitGroup to wait for async handlers
	var wg sync.WaitGroup
	wg.Add(1)

	// Register an async handler
	OnEvent("test.event", func(ctx context.Context, event Event) error {
		defer wg.Done()

		// Verify the event
		e, ok := event.(TestEvent)
		if !ok {
			t.Errorf("Event type mismatch")
			return errors.New("event type mismatch")
		}

		// Verify the data
		if e.data != "test data" {
			t.Errorf("Data mismatch")
			return errors.New("data mismatch")
		}

		// Signal execution
		handlerExecuted <- true
		return nil
	}).WithAsync().Register(bus)

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	errs := bus.Publish(ctx, event)

	// Verify no errors in publish (async errors won't be returned)
	assert.Empty(t, errs)

	// Wait for async handler to complete
	wg.Wait()

	// Verify handler execution
	select {
	case <-handlerExecuted:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Handler not executed within timeout")
	}
}

// TestEventWithPayloadHandling tests handling events with complex payloads
func TestEventWithPayloadHandling(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Create a channel to receive the event
	receivedEvent := make(chan TestEventWithPayload, 1)

	// Register a handler
	bus.Subscribe("test.event.with.payload", func(ctx context.Context, event Event) error {
		e, ok := event.(TestEventWithPayload)
		if !ok {
			return errors.New("event type mismatch")
		}

		receivedEvent <- e
		return nil
	})

	// Create and publish an event
	ctx := context.Background()
	eventTime := time.Now()
	originalEvent := TestEventWithPayload{
		ID:      "123",
		Message: "Hello, World!",
		Time:    eventTime,
	}

	errs := bus.Publish(ctx, originalEvent)
	assert.Empty(t, errs)

	// Verify the received event
	select {
	case received := <-receivedEvent:
		assert.Equal(t, "123", received.ID)
		assert.Equal(t, "Hello, World!", received.Message)
		assert.Equal(t, eventTime.Unix(), received.Time.Unix())
	case <-time.After(time.Second):
		t.Fatal("Event not received within timeout")
	}
}

// TestMultipleHandlers tests that multiple handlers for the same event are all executed
func TestMultipleHandlers(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Create channels to signal handler execution
	handler1Executed := make(chan bool, 1)
	handler2Executed := make(chan bool, 1)
	handler3Executed := make(chan bool, 1)

	// Register handlers
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		handler1Executed <- true
		return nil
	})

	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		handler2Executed <- true
		return nil
	})

	OnEvent("test.event", func(ctx context.Context, event Event) error {
		handler3Executed <- true
		return nil
	}).Register(bus)

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	errs := bus.Publish(ctx, event)

	// Verify no errors
	assert.Empty(t, errs)

	// Verify all handlers executed
	select {
	case <-handler1Executed:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Handler 1 not executed within timeout")
	}

	select {
	case <-handler2Executed:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Handler 2 not executed within timeout")
	}

	select {
	case <-handler3Executed:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Handler 3 not executed within timeout")
	}
}

// TestHandlerError tests that handler errors are properly returned
func TestHandlerError(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Register a handler that returns an error
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		return errors.New("handler error")
	})

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	errs := bus.Publish(ctx, event)

	// Verify the error
	assert.Len(t, errs, 1)
	assert.EqualError(t, errs[0], "handler error")
}

// TestUnsubscribe tests that unsubscribing handlers works correctly
func TestUnsubscribe(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Counter for executions
	executions := 0
	var mu sync.Mutex

	// Create the handler function
	handler := func(ctx context.Context, event Event) error {
		mu.Lock()
		executions++
		mu.Unlock()
		return nil
	}

	// Register the handler
	bus.Subscribe("test.event", handler)

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	bus.Publish(ctx, event)

	// Verify the handler was executed
	mu.Lock()
	assert.Equal(t, 1, executions)
	mu.Unlock()

	// Unsubscribe the handler
	bus.Unsubscribe("test.event", handler)

	// Publish another event
	bus.Publish(ctx, event)

	// Verify the handler was not executed again
	mu.Lock()
	assert.Equal(t, 1, executions)
	mu.Unlock()
}

// TestRetry tests that retry logic works for handlers
func TestRetry(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Attempt counter
	attempts := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(1)

	// Register a handler that fails twice then succeeds
	OnEvent("test.event", func(ctx context.Context, event Event) error {
		defer func() {
			mu.Lock()
			attempts++
			if attempts == 3 {
				wg.Done()
			}
			mu.Unlock()
		}()

		mu.Lock()
		currentAttempt := attempts
		mu.Unlock()

		if currentAttempt < 2 {
			return errors.New("temporary error")
		}

		return nil
	}).WithRetries(2).WithAsync().Register(bus)

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	errs := bus.Publish(ctx, event)

	// No synchronous errors
	assert.Empty(t, errs)

	// Wait for retries to complete
	wg.Wait()

	// Verify attempts
	mu.Lock()
	assert.Equal(t, 3, attempts)
	mu.Unlock()
}

// TestBaseEvent tests the BaseEvent implementation
func TestBaseEvent(t *testing.T) {
	// Create a base event
	event := NewEvent("custom.event", map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	})

	// Verify the event name
	assert.Equal(t, "custom.event", event.Name())

	// Verify the event payload
	baseEvent, ok := event.(BaseEvent)
	assert.True(t, ok)

	payload, ok := baseEvent.Payload.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value1", payload["key1"])
	assert.Equal(t, 123, payload["key2"])
}

// TestEventHandlerModes tests that handler modes work correctly
func TestEventHandlerModes(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Channels to signal execution
	syncExecuted := make(chan bool, 1)
	asyncExecuted := make(chan bool, 1)

	// Register a sync handler
	OnEvent("test.event", func(ctx context.Context, event Event) error {
		syncExecuted <- true
		return nil
	}).WithSync().Register(bus)

	// Register an async handler
	var wg sync.WaitGroup
	wg.Add(1)
	OnEvent("test.event", func(ctx context.Context, event Event) error {
		defer wg.Done()
		asyncExecuted <- true
		return nil
	}).WithAsync().Register(bus)

	// Publish an event
	ctx := context.Background()
	event := TestEvent{data: "test data"}
	bus.Publish(ctx, event)

	// Verify sync handler executed immediately
	select {
	case <-syncExecuted:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Sync handler not executed within timeout")
	}

	// Wait for async handler
	wg.Wait()

	// Verify async handler executed
	select {
	case <-asyncExecuted:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Async handler not executed within timeout")
	}
}

// TestShutdown tests that the event bus can be shut down
func TestShutdown(t *testing.T) {
	// Create an event bus
	bus := NewEventBus(EventBusParams{})

	// Shutdown the bus
	bus.Shutdown()

	// This is mostly a smoke test to ensure no panics
	assert.NotNil(t, bus.ctx)
	assert.NotNil(t, bus.cancel)
}
