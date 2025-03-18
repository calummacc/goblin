package di_example

import (
	"fmt"
	"goblin/di"
	"log"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

// Database represents a database connection
type Database struct {
	name string
}

// NewDatabase creates a new database connection (singleton)
func NewDatabase() *Database {
	return &Database{
		name: "MainDB",
	}
}

// UserRepository handles user data access
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new user repository (transient)
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// RequestContext holds request-specific data
type RequestContext struct {
	RequestID string
	UserID    string
}

// NewRequestContext creates a new request context (request-scoped)
func NewRequestContext(c *gin.Context) *RequestContext {
	return &RequestContext{
		RequestID: c.GetHeader("X-Request-ID"),
		UserID:    c.GetHeader("X-User-ID"),
	}
}

// UserService handles user business logic
type UserService struct {
	repo    *UserRepository
	reqCtx  *RequestContext
	counter int
	mutex   sync.Mutex
}

// NewUserService creates a new user service (transient)
func NewUserService(repo *UserRepository, reqCtx *RequestContext) *UserService {
	return &UserService{
		repo:   repo,
		reqCtx: reqCtx,
	}
}

// IncrementCounter increments the counter
func (s *UserService) IncrementCounter() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counter++
	return s.counter
}

func main() {
	// Create DI container
	container := di.NewContainer()

	// Register providers with different scopes
	err := container.Register(NewDatabase, di.Singleton)
	if err != nil {
		log.Fatal(err)
	}

	err = container.Register(NewUserRepository, di.Transient)
	if err != nil {
		log.Fatal(err)
	}

	err = container.Register(NewRequestContext, di.RequestScoped)
	if err != nil {
		log.Fatal(err)
	}

	err = container.Register(NewUserService, di.Transient)
	if err != nil {
		log.Fatal(err)
	}

	// Create Gin router
	router := gin.Default()

	// Example handler using DI
	router.GET("/api/users", func(c *gin.Context) {
		// Resolve UserService for this request
		serviceInstance, err := container.Resolve(reflect.TypeOf(&UserService{}), c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		service := serviceInstance.(*UserService)
		count := service.IncrementCounter()

		c.JSON(200, gin.H{
			"message":    "User service example",
			"requestID":  service.reqCtx.RequestID,
			"userID":     service.reqCtx.UserID,
			"dbName":     service.repo.db.name,
			"counter":    count,
			"servicePtr": fmt.Sprintf("%p", service),
		})
	})

	// Start server
	log.Println("Server starting on http://localhost:8080")
	router.Run(":8080")
}
