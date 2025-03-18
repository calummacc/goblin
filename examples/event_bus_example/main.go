package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"goblin/events"

	"go.uber.org/fx"
)

// UserCreatedEvent represents an event triggered when a user is created
type UserCreatedEvent struct {
	UserID   string
	Username string
	Email    string
	Time     time.Time
}

// Name returns the event name
func (e UserCreatedEvent) Name() string {
	return "user.created"
}

// UserUpdatedEvent represents an event triggered when a user is updated
type UserUpdatedEvent struct {
	UserID      string
	OldUsername string
	NewUsername string
	Time        time.Time
}

// Name returns the event name
func (e UserUpdatedEvent) Name() string {
	return "user.updated"
}

// NotificationService sends notifications based on events
type NotificationService struct {
	mu            sync.Mutex
	notifications []string
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifications: make([]string, 0),
	}
}

// GetNotifications returns all recorded notifications
func (s *NotificationService) GetNotifications() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy to avoid concurrent modification
	result := make([]string, len(s.notifications))
	copy(result, s.notifications)
	return result
}

// AddNotification adds a notification
func (s *NotificationService) AddNotification(notification string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notifications = append(s.notifications, notification)
	log.Printf("Notification: %s", notification)
}

// HandleUserCreated handles user created events
func (s *NotificationService) HandleUserCreated(ctx context.Context, event events.Event) error {
	e, ok := event.(UserCreatedEvent)
	if !ok {
		return fmt.Errorf("expected UserCreatedEvent, got %T", event)
	}

	notification := fmt.Sprintf("Welcome %s! Your account was created with email %s", e.Username, e.Email)
	s.AddNotification(notification)
	return nil
}

// HandleUserUpdated handles user updated events
func (s *NotificationService) HandleUserUpdated(ctx context.Context, event events.Event) error {
	e, ok := event.(UserUpdatedEvent)
	if !ok {
		return fmt.Errorf("expected UserUpdatedEvent, got %T", event)
	}

	notification := fmt.Sprintf("Hi %s! Your username was changed from %s to %s", e.NewUsername, e.OldUsername, e.NewUsername)
	s.AddNotification(notification)
	return nil
}

// AnalyticsService records analytics based on events
type AnalyticsService struct {
	mu     sync.Mutex
	events map[string]int
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		events: make(map[string]int),
	}
}

// RecordEvent records an event
func (s *AnalyticsService) RecordEvent(eventName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[eventName]++
	log.Printf("Analytics: Recorded event %s (count: %d)", eventName, s.events[eventName])
}

// GetEventCounts returns the event counts
func (s *AnalyticsService) GetEventCounts() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy to avoid concurrent modification
	result := make(map[string]int)
	for k, v := range s.events {
		result[k] = v
	}
	return result
}

// HandleEvent handles any event for analytics purposes
func (s *AnalyticsService) HandleEvent(ctx context.Context, event events.Event) error {
	s.RecordEvent(event.Name())

	// Simulate slow processing for async handlers
	time.Sleep(100 * time.Millisecond)

	return nil
}

// AuditLogger logs audit events
type AuditLogger struct {
	mu   sync.Mutex
	logs []string
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		logs: make([]string, 0),
	}
}

// LogEvent logs an event
func (l *AuditLogger) LogEvent(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, message)
	log.Printf("Audit: %s", message)
}

// GetLogs returns the logs
func (l *AuditLogger) GetLogs() []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Return a copy to avoid concurrent modification
	result := make([]string, len(l.logs))
	copy(result, l.logs)
	return result
}

// HandleUserCreated handles user created events
func (l *AuditLogger) HandleUserCreated(ctx context.Context, event events.Event) error {
	e, ok := event.(UserCreatedEvent)
	if !ok {
		return fmt.Errorf("expected UserCreatedEvent, got %T", event)
	}

	l.LogEvent(fmt.Sprintf("User created: %s (%s) at %s", e.Username, e.UserID, e.Time.Format(time.RFC3339)))
	return nil
}

// HandleUserUpdated handles user updated events
func (l *AuditLogger) HandleUserUpdated(ctx context.Context, event events.Event) error {
	e, ok := event.(UserUpdatedEvent)
	if !ok {
		return fmt.Errorf("expected UserUpdatedEvent, got %T", event)
	}

	l.LogEvent(fmt.Sprintf("User updated: %s (%s) changed username from %s to %s at %s",
		e.NewUsername, e.UserID, e.OldUsername, e.NewUsername, e.Time.Format(time.RFC3339)))
	return nil
}

// RegisterServices registers all services with the event bus
func RegisterServices(eventBus *events.EventBus, notificationService *NotificationService, analyticsService *AnalyticsService, auditLogger *AuditLogger) {
	// Register notification service handlers (synchronous)
	events.OnEvent("user.created", notificationService.HandleUserCreated).
		WithName("notification_user_created").
		Register(eventBus)

	events.OnEvent("user.updated", notificationService.HandleUserUpdated).
		WithName("notification_user_updated").
		Register(eventBus)

	// Register analytics service handler (asynchronous with retries)
	events.OnEvent("user.created", analyticsService.HandleEvent).
		WithName("analytics_user_created").
		WithAsync().
		WithRetries(3).
		Register(eventBus)

	events.OnEvent("user.updated", analyticsService.HandleEvent).
		WithName("analytics_user_updated").
		WithAsync().
		WithRetries(3).
		Register(eventBus)

	// Register audit logger handlers (synchronous)
	events.OnEvent("user.created", auditLogger.HandleUserCreated).
		WithName("audit_user_created").
		Register(eventBus)

	events.OnEvent("user.updated", auditLogger.HandleUserUpdated).
		WithName("audit_user_updated").
		Register(eventBus)
}

// UserService manages user operations and emits events
type UserService struct {
	eventBus *events.EventBus
	users    map[string]UserInfo
	mu       sync.Mutex
}

// UserInfo holds user information
type UserInfo struct {
	ID       string
	Username string
	Email    string
}

// NewUserService creates a new user service
func NewUserService(eventBus *events.EventBus) *UserService {
	return &UserService{
		eventBus: eventBus,
		users:    make(map[string]UserInfo),
	}
}

// CreateUser creates a new user and emits an event
func (s *UserService) CreateUser(ctx context.Context, id, username, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store user
	s.users[id] = UserInfo{
		ID:       id,
		Username: username,
		Email:    email,
	}

	// Emit event
	event := UserCreatedEvent{
		UserID:   id,
		Username: username,
		Email:    email,
		Time:     time.Now(),
	}

	log.Printf("Emitting user.created event for user %s", username)
	errors := s.eventBus.Publish(ctx, event)

	if len(errors) > 0 {
		// In this example we'll just log errors, but in real-world applications
		// you might want to handle them differently
		log.Printf("Errors publishing user.created event: %v", errors)
	}

	return nil
}

// UpdateUsername updates a user's username and emits an event
func (s *UserService) UpdateUsername(ctx context.Context, id, newUsername string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get user
	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user %s not found", id)
	}

	oldUsername := user.Username

	// Update user
	user.Username = newUsername
	s.users[id] = user

	// Emit event
	event := UserUpdatedEvent{
		UserID:      id,
		OldUsername: oldUsername,
		NewUsername: newUsername,
		Time:        time.Now(),
	}

	log.Printf("Emitting user.updated event for user %s", newUsername)
	errors := s.eventBus.Publish(ctx, event)

	if len(errors) > 0 {
		log.Printf("Errors publishing user.updated event: %v", errors)
	}

	return nil
}

// Module provides the dependencies
func Module() fx.Option {
	return fx.Options(
		events.NewEventBusModule(),
		fx.Provide(
			NewUserService,
			NewNotificationService,
			NewAnalyticsService,
			NewAuditLogger,
		),
		fx.Invoke(RegisterServices),
	)
}

func main() {
	app := fx.New(
		Module(),
		fx.Invoke(func(userService *UserService) {
			// Create a user
			ctx := context.Background()
			userService.CreateUser(ctx, "user1", "johndoe", "john@example.com")

			// Give time for async handlers to complete
			time.Sleep(200 * time.Millisecond)

			// Update username
			userService.UpdateUsername(ctx, "user1", "john.doe")

			// Give time for async handlers to complete
			time.Sleep(200 * time.Millisecond)
		}),
		fx.Invoke(func(notificationService *NotificationService, analyticsService *AnalyticsService, auditLogger *AuditLogger) {
			// Wait to ensure all async handlers have completed
			time.Sleep(500 * time.Millisecond)

			// Print summary
			fmt.Println("\n--- Summary ---")

			fmt.Println("\nNotifications:")
			for _, notification := range notificationService.GetNotifications() {
				fmt.Printf("- %s\n", notification)
			}

			fmt.Println("\nAnalytics:")
			eventCounts := analyticsService.GetEventCounts()
			for event, count := range eventCounts {
				fmt.Printf("- %s: %d\n", event, count)
			}

			fmt.Println("\nAudit Logs:")
			for _, log := range auditLogger.GetLogs() {
				fmt.Printf("- %s\n", log)
			}
		}),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
