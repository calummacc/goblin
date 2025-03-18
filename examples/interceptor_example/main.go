package main

import (
	"context"
	"log"
	"time"

	"goblin/interceptor"

	"github.com/gin-gonic/gin"
)

// CustomInterceptor demonstrates a custom interceptor implementation
type CustomInterceptor struct {
	interceptor.BaseInterceptor
	name string
}

func NewCustomInterceptor(name string) *CustomInterceptor {
	return &CustomInterceptor{
		name: name,
	}
}

func (i *CustomInterceptor) Before(ctx *interceptor.ExecutionContext) error {
	log.Printf("[%s] Before request to %s\n", i.name, ctx.Path)
	return nil
}

func (i *CustomInterceptor) After(ctx *interceptor.ExecutionContext) error {
	log.Printf("[%s] After request to %s\n", i.name, ctx.Path)
	return nil
}

func main() {
	// Create a new Gin engine
	engine := gin.Default()

	// Create interceptor manager
	manager := interceptor.NewManager()

	// Register logging interceptor
	if err := manager.Register(interceptor.Config{
		Name:        "logging",
		Priority:    1,
		Interceptor: interceptor.NewLoggingInterceptor(),
	}); err != nil {
		log.Fatal(err)
	}

	// Register metrics interceptor
	metricsInterceptor := interceptor.NewMetricsInterceptor()
	if err := manager.Register(interceptor.Config{
		Name:        "metrics",
		Priority:    2,
		Interceptor: metricsInterceptor,
	}); err != nil {
		log.Fatal(err)
	}

	// Register custom interceptor for specific paths
	if err := manager.Register(interceptor.Config{
		Name:        "custom",
		Priority:    3,
		Path:        "/api/*",
		Methods:     []string{"GET", "POST"},
		Interceptor: NewCustomInterceptor("API"),
	}); err != nil {
		log.Fatal(err)
	}

	// Apply all interceptors to the engine
	if err := manager.Use(engine); err != nil {
		log.Fatal(err)
	}

	// Initialize all lifecycle interceptors
	if err := manager.OnRegister(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Add some routes
	engine.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 100) // Simulate work
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	engine.GET("/api/users", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 200) // Simulate work
		c.JSON(200, gin.H{
			"users": []string{"user1", "user2", "user3"},
		})
	})

	engine.POST("/api/users", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 300) // Simulate work
		c.JSON(400, gin.H{
			"error": "not implemented",
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

	// Print metrics before shutdown
	metrics := metricsInterceptor.GetMetrics()
	log.Printf("\nCurrent Metrics:")
	for route, m := range metrics {
		log.Printf("\n%s:", route)
		log.Printf("  Total Requests: %d", m.TotalRequests)
		log.Printf("  Total Errors: %d", m.TotalErrors)
		log.Printf("  Average Response Time: %v", time.Duration(m.AverageTimeNs))
	}

	// Shutdown all lifecycle interceptors
	if err := manager.OnShutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
