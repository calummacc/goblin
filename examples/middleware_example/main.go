package main

import (
	"context"
	"log"
	"time"

	"goblin/middleware"

	"github.com/gin-gonic/gin"
)

// CustomMiddleware demonstrates a custom middleware implementation
type CustomMiddleware struct {
	middleware.BaseMiddleware
	name string
}

func NewCustomMiddleware(name string) *CustomMiddleware {
	return &CustomMiddleware{
		name: name,
	}
}

func (m *CustomMiddleware) Handle(c *gin.Context) {
	log.Printf("[%s] Before request", m.name)
	c.Next()
	log.Printf("[%s] After request", m.name)
}

func (m *CustomMiddleware) OnRegister(ctx context.Context) error {
	log.Printf("[%s] Middleware registered", m.name)
	return nil
}

func (m *CustomMiddleware) OnShutdown(ctx context.Context) error {
	log.Printf("[%s] Middleware shutting down", m.name)
	return nil
}

func main() {
	// Create a new Gin engine
	engine := gin.Default()

	// Create middleware manager
	manager := middleware.NewManager()

	// Register global middlewares
	if err := manager.Register(middleware.Config{
		Name: "logger",
		Options: middleware.Options{
			Global:   true,
			Priority: 1,
		},
		Middleware: middleware.NewLoggerMiddleware(),
	}); err != nil {
		log.Fatal(err)
	}

	if err := manager.Register(middleware.Config{
		Name: "recovery",
		Options: middleware.Options{
			Global:   true,
			Priority: 0,
		},
		Middleware: middleware.NewRecoveryMiddleware(),
	}); err != nil {
		log.Fatal(err)
	}

	// Register auth middleware for specific paths
	if err := manager.Register(middleware.Config{
		Name: "auth",
		Options: middleware.Options{
			Path:     "/protected/*",
			Methods:  []string{"GET", "POST"},
			Priority: 2,
		},
		Middleware: middleware.NewAuthMiddleware(),
	}); err != nil {
		log.Fatal(err)
	}

	// Register custom middleware
	if err := manager.Register(middleware.Config{
		Name: "custom",
		Options: middleware.Options{
			Global:   true,
			Priority: 3,
		},
		Middleware: NewCustomMiddleware("custom"),
	}); err != nil {
		log.Fatal(err)
	}

	// Register middleware group
	if err := manager.RegisterGroup(middleware.Group{
		Name: "api",
		Middlewares: []interface{}{
			middleware.NewLoggerMiddleware(),
			middleware.NewAuthMiddleware(),
			NewCustomMiddleware("api"),
		},
	}); err != nil {
		log.Fatal(err)
	}

	// Apply all middlewares to the engine
	if err := manager.Use(engine); err != nil {
		log.Fatal(err)
	}

	// Initialize all lifecycle middlewares
	if err := manager.OnRegister(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Add some routes
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	engine.GET("/protected/resource", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "protected resource",
		})
	})

	// Start the server
	go func() {
		if err := engine.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for a while
	time.Sleep(time.Second * 30)

	// Shutdown all lifecycle middlewares
	if err := manager.OnShutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
