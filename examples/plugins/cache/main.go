package cache

import (
	"context"
	"fmt"
	"goblin/plugin"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// CachePlugin implements the Plugin interface for caching
type CachePlugin struct {
	config    *plugin.PluginConfig
	logger    *log.Logger
	container *fx.App
}

// NewCachePlugin creates a new cache plugin
func NewCachePlugin(config *plugin.PluginConfig, logger *log.Logger) *CachePlugin {
	return &CachePlugin{
		config: config,
		logger: logger,
	}
}

// Name returns the plugin name
func (p *CachePlugin) Name() string {
	return "cache"
}

// Version returns the plugin version
func (p *CachePlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *CachePlugin) Description() string {
	return "Caching plugin for Goblin Framework"
}

// Dependencies returns the plugin dependencies
func (p *CachePlugin) Dependencies() []string {
	return []string{} // No dependencies
}

// OnRegister is called when the plugin is registered
func (p *CachePlugin) OnRegister(ctx context.Context) error {
	p.logger.Printf("Registering cache plugin...")
	return nil
}

// OnStart is called when the application starts
func (p *CachePlugin) OnStart(ctx context.Context) error {
	p.logger.Printf("Starting cache plugin...")
	return nil
}

// OnStop is called when the application stops
func (p *CachePlugin) OnStop(ctx context.Context) error {
	p.logger.Printf("Stopping cache plugin...")
	return nil
}

// RegisterRoutes registers the plugin's routes
func (p *CachePlugin) RegisterRoutes(router *gin.Engine) error {
	cache := router.Group("/cache")
	{
		cache.GET("/:key", p.handleGet)
		cache.POST("/:key", p.handleSet)
		cache.DELETE("/:key", p.handleDelete)
		cache.GET("/stats", p.handleStats)
	}
	return nil
}

// RegisterDependencies registers the plugin's dependencies
func (p *CachePlugin) RegisterDependencies(app *fx.App) error {
	// Create a module with all dependencies
	module := fx.Module("cache",
		fx.Provide(
			fx.Annotate(
				NewCacheService,
				fx.As(new(*CacheService)),
			),
		),
	)

	// Create a new app with the module
	p.container = fx.New(module)

	// Start the app to initialize dependencies
	if err := p.container.Start(context.Background()); err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	return nil
}

// CacheItem represents a cached item
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// CacheService handles caching operations
type CacheService struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
	stats struct {
		hits   int64
		misses int64
		evicts int64
	}
}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	cs := &CacheService{
		items: make(map[string]*CacheItem),
	}

	// Start cleanup goroutine
	go cs.cleanup()

	return cs
}

// Get retrieves a value from the cache
func (s *CacheService) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[key]
	if !exists {
		s.stats.misses++
		return nil, false
	}

	if time.Now().After(item.Expiration) {
		s.stats.misses++
		return nil, false
	}

	s.stats.hits++
	return item.Value, true
}

// Set stores a value in the cache
func (s *CacheService) Set(key string, value interface{}, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = &CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration),
	}
}

// Delete removes a value from the cache
func (s *CacheService) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}

// Stats returns cache statistics
func (s *CacheService) Stats() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]int64{
		"hits":   s.stats.hits,
		"misses": s.stats.misses,
		"evicts": s.stats.evicts,
		"items":  int64(len(s.items)),
	}
}

// cleanup removes expired items from the cache
func (s *CacheService) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, item := range s.items {
			if now.After(item.Expiration) {
				delete(s.items, key)
				s.stats.evicts++
			}
		}
		s.mu.Unlock()
	}
}

// handleGet handles cache get requests
func (p *CachePlugin) handleGet(c *gin.Context) {
	key := c.Param("key")

	// Get cache service from container
	var cacheService *CacheService
	if err := fx.Populate(p.container, &cacheService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache service"})
		return
	}

	value, exists := cacheService.Get(key)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

// handleSet handles cache set requests
func (p *CachePlugin) handleSet(c *gin.Context) {
	key := c.Param("key")
	var req struct {
		Value    interface{} `json:"value" binding:"required"`
		Duration string      `json:"duration"` // e.g., "1h", "30m"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	duration := time.Hour // Default duration
	if req.Duration != "" {
		var err error
		duration, err = time.ParseDuration(req.Duration)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format"})
			return
		}
	}

	// Get cache service from container
	var cacheService *CacheService
	if err := fx.Populate(p.container, &cacheService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache service"})
		return
	}

	cacheService.Set(key, req.Value, duration)

	c.JSON(http.StatusOK, gin.H{
		"message": "Value cached successfully",
		"key":     key,
	})
}

// handleDelete handles cache delete requests
func (p *CachePlugin) handleDelete(c *gin.Context) {
	key := c.Param("key")

	// Get cache service from container
	var cacheService *CacheService
	if err := fx.Populate(p.container, &cacheService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache service"})
		return
	}

	cacheService.Delete(key)

	c.JSON(http.StatusOK, gin.H{
		"message": "Key deleted successfully",
		"key":     key,
	})
}

// handleStats handles cache statistics requests
func (p *CachePlugin) handleStats(c *gin.Context) {
	// Get cache service from container
	var cacheService *CacheService
	if err := fx.Populate(p.container, &cacheService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache service"})
		return
	}

	c.JSON(http.StatusOK, cacheService.Stats())
}
