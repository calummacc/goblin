package middleware

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

// Manager manages middleware registration and execution
type Manager struct {
	mu sync.RWMutex
	// middlewares holds all registered middleware configs
	middlewares []Config
	// groups holds middleware groups by name
	groups map[string]Group
	// lifecycleMiddlewares holds middleware that implement LifecycleMiddleware
	lifecycleMiddlewares []LifecycleMiddleware
}

// NewManager creates a new middleware manager
func NewManager() *Manager {
	return &Manager{
		groups: make(map[string]Group),
	}
}

// Register registers a middleware with the given configuration
func (m *Manager) Register(config Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate config
	if config.Name == "" {
		return fmt.Errorf("middleware name is required")
	}
	if config.Middleware == nil {
		return fmt.Errorf("middleware implementation is required")
	}

	// Check for lifecycle middleware
	if lm, ok := config.Middleware.(LifecycleMiddleware); ok {
		m.lifecycleMiddlewares = append(m.lifecycleMiddlewares, lm)
	}

	m.middlewares = append(m.middlewares, config)
	return nil
}

// RegisterGroup registers a group of middleware
func (m *Manager) RegisterGroup(group Group) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group.Name == "" {
		return fmt.Errorf("group name is required")
	}
	if len(group.Middlewares) == 0 {
		return fmt.Errorf("group must contain at least one middleware")
	}

	m.groups[group.Name] = group
	return nil
}

// Use applies middleware to a gin.Engine instance
func (m *Manager) Use(engine *gin.Engine) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Sort middlewares by priority
	sorted := make([]Config, len(m.middlewares))
	copy(sorted, m.middlewares)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Options.Priority < sorted[j].Options.Priority
	})

	// Apply middlewares
	for _, config := range sorted {
		if err := m.applyMiddleware(engine, config); err != nil {
			return err
		}
	}

	return nil
}

// applyMiddleware applies a single middleware to the engine
func (m *Manager) applyMiddleware(engine *gin.Engine, config Config) error {
	var handler gin.HandlerFunc

	switch mw := config.Middleware.(type) {
	case gin.HandlerFunc:
		handler = mw
	case Middleware:
		handler = mw.Handle
	case func(*gin.Context):
		handler = mw
	default:
		return fmt.Errorf("unsupported middleware type for %s", config.Name)
	}

	if config.Options.Global {
		engine.Use(handler)
		return nil
	}

	// Apply to specific path and methods
	if config.Options.Path != "" {
		methods := config.Options.Methods
		if len(methods) == 0 {
			methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
		}

		for _, method := range methods {
			engine.Handle(method, config.Options.Path, handler)
		}
	}

	return nil
}

// OnRegister calls OnRegister for all lifecycle middlewares
func (m *Manager) OnRegister(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, mw := range m.lifecycleMiddlewares {
		if err := mw.OnRegister(ctx); err != nil {
			return fmt.Errorf("failed to register middleware: %w", err)
		}
	}
	return nil
}

// OnShutdown calls OnShutdown for all lifecycle middlewares
func (m *Manager) OnShutdown(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, mw := range m.lifecycleMiddlewares {
		if err := mw.OnShutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown middleware: %w", err)
		}
	}
	return nil
}
