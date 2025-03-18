package interceptor

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

// Manager manages interceptor registration and execution
type Manager struct {
	mu sync.RWMutex
	// interceptors holds all registered interceptor configs
	interceptors []Config
	// lifecycleInterceptors holds interceptors that implement LifecycleInterceptor
	lifecycleInterceptors []LifecycleInterceptor
}

// NewManager creates a new interceptor manager
func NewManager() *Manager {
	return &Manager{}
}

// Register registers an interceptor with the given configuration
func (m *Manager) Register(config Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate config
	if config.Name == "" {
		return fmt.Errorf("interceptor name is required")
	}
	if config.Interceptor == nil {
		return fmt.Errorf("interceptor implementation is required")
	}

	// Check for lifecycle interceptor
	if li, ok := config.Interceptor.(LifecycleInterceptor); ok {
		m.lifecycleInterceptors = append(m.lifecycleInterceptors, li)
	}

	m.interceptors = append(m.interceptors, config)
	return nil
}

// Use applies interceptors to a gin.Engine instance
func (m *Manager) Use(engine *gin.Engine) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Sort interceptors by priority
	sorted := make([]Config, len(m.interceptors))
	copy(sorted, m.interceptors)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority < sorted[j].Priority
	})

	// Create a middleware that wraps all interceptors
	engine.Use(func(c *gin.Context) {
		// Create execution context
		execCtx := &ExecutionContext{
			GinContext: c,
			Handler:    c.Handler(),
			Path:       c.FullPath(),
			Method:     c.Request.Method,
			Data:       make(map[string]interface{}),
		}

		// Run Before interceptors
		for _, config := range sorted {
			if !m.shouldRunInterceptor(config, c) {
				continue
			}

			var interceptor Interceptor
			switch i := config.Interceptor.(type) {
			case Interceptor:
				interceptor = i
			default:
				continue
			}

			if err := interceptor.Before(execCtx); err != nil {
				c.Abort()
				return
			}
		}

		// If the request wasn't aborted, process the request
		c.Next()

		// Run After interceptors in reverse order
		for i := len(sorted) - 1; i >= 0; i-- {
			config := sorted[i]
			if !m.shouldRunInterceptor(config, c) {
				continue
			}

			var interceptor Interceptor
			switch i := config.Interceptor.(type) {
			case Interceptor:
				interceptor = i
			default:
				continue
			}

			if err := interceptor.After(execCtx); err != nil {
				// Log error but don't abort since we're in the After phase
				fmt.Printf("Error in After interceptor %s: %v\n", config.Name, err)
			}
		}
	})

	return nil
}

// shouldRunInterceptor checks if an interceptor should run for the current request
func (m *Manager) shouldRunInterceptor(config Config, c *gin.Context) bool {
	// If no path or methods specified, run for all requests
	if config.Path == "" && len(config.Methods) == 0 {
		return true
	}

	// Check path match
	if config.Path != "" && config.Path != c.FullPath() {
		return false
	}

	// Check method match
	if len(config.Methods) > 0 {
		methodMatch := false
		for _, method := range config.Methods {
			if method == c.Request.Method {
				methodMatch = true
				break
			}
		}
		if !methodMatch {
			return false
		}
	}

	return true
}

// OnRegister calls OnRegister for all lifecycle interceptors
func (m *Manager) OnRegister(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, interceptor := range m.lifecycleInterceptors {
		if err := interceptor.OnRegister(ctx); err != nil {
			return fmt.Errorf("failed to register interceptor: %w", err)
		}
	}
	return nil
}

// OnShutdown calls OnShutdown for all lifecycle interceptors
func (m *Manager) OnShutdown(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, interceptor := range m.lifecycleInterceptors {
		if err := interceptor.OnShutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown interceptor: %w", err)
		}
	}
	return nil
}
