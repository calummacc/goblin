package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"goblin/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// UserService represents a service for user operations
type UserService struct {
	db *Database
}

// NewUserService creates a new UserService
func NewUserService(db *Database) *UserService {
	return &UserService{db: db}
}

// GetUsers returns a list of users
func (s *UserService) GetUsers() []string {
	return []string{"user1", "user2", "user3"}
}

// Database represents a database connection
type Database struct {
	connected bool
}

// NewDatabase creates a new database connection
func NewDatabase() *Database {
	return &Database{connected: false}
}

// Connect connects to the database
func (d *Database) Connect() error {
	time.Sleep(time.Second) // Simulate connection
	d.connected = true
	return nil
}

// Disconnect disconnects from the database
func (d *Database) Disconnect() error {
	d.connected = false
	return nil
}

// UserController handles user-related HTTP requests
type UserController struct {
	userService *UserService
}

// NewUserController creates a new UserController
func NewUserController(userService *UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUsers handles GET /users request
func (c *UserController) GetUsers(ctx *gin.Context) {
	users := c.userService.GetUsers()
	ctx.JSON(200, gin.H{"users": users})
}

// UserModule represents the user module
type UserModule struct {
	*core.BaseModule
}

// NewUserModule creates a new UserModule
func NewUserModule() *UserModule {
	module := &UserModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			NewDatabase,
			NewUserService,
		},
		Controllers: []interface{}{
			NewUserController,
		},
	})
	return module
}

// OnModuleInit initializes the user module
func (m *UserModule) OnModuleInit(ctx context.Context) error {
	// Get the database instance from the module's providers
	db := m.GetMetadata().Providers[0].(*Database)
	return db.Connect()
}

// OnModuleDestroy cleans up the user module
func (m *UserModule) OnModuleDestroy(ctx context.Context) error {
	// Get the database instance from the module's providers
	db := m.GetMetadata().Providers[0].(*Database)
	return db.Disconnect()
}

// AuthModule represents the authentication module
type AuthModule struct {
	*core.BaseModule
}

// NewAuthModule creates a new AuthModule
func NewAuthModule() *AuthModule {
	module := &AuthModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			func() string { return "jwt-secret" },
		},
	})
	return module
}

// AppModule represents the main application module
type AppModule struct {
	*core.BaseModule
}

// NewAppModule creates a new AppModule
func NewAppModule() *AppModule {
	module := &AppModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Imports: []core.Module{
			NewUserModule(),
			NewAuthModule(),
		},
	})
	return module
}

func main() {
	// Create the module manager
	moduleManager := core.NewModuleManager()

	// Create and register the app module
	appModule := NewAppModule()
	if err := moduleManager.RegisterModule(appModule); err != nil {
		log.Fatalf("Failed to register app module: %v", err)
	}

	// Create Gin engine
	engine := gin.Default()

	// Create Fx application
	app := fx.New(
		fx.Provide(
			func() *gin.Engine { return engine },
		),
		fx.Provide(moduleManager.GetModuleProviders(appModule)...),
		fx.Invoke(moduleManager.GetModuleControllers(appModule)...),
	)

	// Initialize modules
	ctx := context.Background()
	if err := moduleManager.InitializeModules(ctx); err != nil {
		log.Fatalf("Failed to initialize modules: %v", err)
	}

	// Start the application
	go func() {
		if err := app.Start(ctx); err != nil {
			log.Printf("Error starting application: %v", err)
		}
	}()

	// Wait for a moment to ensure the server is running
	time.Sleep(time.Second)

	// Make a test request
	resp, err := http.Get("http://localhost:8080/users")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	fmt.Printf("Response: %s\n", string(body))

	// Stop the application
	if err := app.Stop(ctx); err != nil {
		log.Printf("Error stopping application: %v", err)
	}

	// Destroy modules
	if err := moduleManager.DestroyModules(ctx); err != nil {
		log.Printf("Error destroying modules: %v", err)
	}
}
