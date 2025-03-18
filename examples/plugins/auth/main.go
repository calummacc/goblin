package auth

import (
	"context"
	"fmt"
	"goblin/plugin"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/fx"
)

// Dependencies represents the dependencies required by the auth plugin
type Dependencies struct {
	fx.In

	Config *plugin.PluginConfig
	Logger *log.Logger
}

// Result represents the dependencies provided by the auth plugin
type Result struct {
	fx.Out

	Secret      []byte
	UserService *UserService
}

// AuthPlugin implements the Plugin interface for authentication
type AuthPlugin struct {
	config *plugin.PluginConfig
	logger *log.Logger
	deps   Dependencies
}

// NewAuthPlugin creates a new auth plugin
func NewAuthPlugin(deps Dependencies) *AuthPlugin {
	return &AuthPlugin{
		config: deps.Config,
		logger: deps.Logger,
		deps:   deps,
	}
}

// Name returns the plugin name
func (p *AuthPlugin) Name() string {
	return "auth"
}

// Version returns the plugin version
func (p *AuthPlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *AuthPlugin) Description() string {
	return "Authentication plugin for Goblin Framework"
}

// Dependencies returns the plugin dependencies
func (p *AuthPlugin) Dependencies() []string {
	return []string{} // No dependencies
}

// OnRegister is called when the plugin is registered
func (p *AuthPlugin) OnRegister(ctx context.Context) error {
	p.logger.Printf("Registering auth plugin...")
	return nil
}

// OnStart is called when the application starts
func (p *AuthPlugin) OnStart(ctx context.Context) error {
	p.logger.Printf("Starting auth plugin...")
	return nil
}

// OnStop is called when the application stops
func (p *AuthPlugin) OnStop(ctx context.Context) error {
	p.logger.Printf("Stopping auth plugin...")
	return nil
}

// RegisterRoutes registers the plugin's routes
func (p *AuthPlugin) RegisterRoutes(router *gin.Engine) error {
	auth := router.Group("/auth")
	{
		auth.POST("/login", p.handleLogin)
		auth.POST("/register", p.handleRegister)
		auth.GET("/profile", p.authMiddleware(), p.handleProfile)
	}
	return nil
}

// RegisterDependencies registers the plugin's dependencies
func (p *AuthPlugin) RegisterDependencies(app *fx.App) error {
	// Create a module with all dependencies
	module := fx.Module("auth",
		fx.Provide(
			// Provide JWT secret
			fx.Annotate(
				func() []byte {
					return []byte("your-secret-key") // In production, use environment variable
				},
				fx.As(new([]byte)),
			),
			// Provide user service
			fx.Annotate(
				NewUserService,
				fx.As(new(*UserService)),
			),
		),
	)

	// Create a new app with the module
	newApp := fx.New(module)

	// Start the app to initialize dependencies
	if err := newApp.Start(context.Background()); err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	// Stop the app
	if err := newApp.Stop(context.Background()); err != nil {
		return fmt.Errorf("failed to stop app: %w", err)
	}

	return nil
}

// User represents a user entity
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Password is not sent in responses
}

// UserService handles user-related operations
type UserService struct {
	users map[string]*User
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, password string) (*User, error) {
	if _, exists := s.users[username]; exists {
		return nil, fmt.Errorf("user already exists")
	}

	user := &User{
		ID:       fmt.Sprintf("user_%d", len(s.users)+1),
		Username: username,
		Password: password, // In production, hash the password
	}

	s.users[username] = user
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// handleLogin handles user login
func (p *AuthPlugin) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a module with all dependencies
	module := fx.Module("auth",
		fx.Provide(
			// Provide JWT secret
			fx.Annotate(
				func() []byte {
					return []byte("your-secret-key") // In production, use environment variable
				},
				fx.As(new([]byte)),
			),
			// Provide user service
			fx.Annotate(
				NewUserService,
				fx.As(new(*UserService)),
			),
		),
	)

	// Create a new app with the module
	app := fx.New(module)

	// Start the app to initialize dependencies
	if err := app.Start(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start app"})
		return
	}
	defer app.Stop(context.Background())

	// Get user service and secret from app
	var result struct {
		fx.In
		UserService *UserService
		Secret      []byte
	}

	if err := fx.Populate(app, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dependencies"})
		return
	}

	user, err := result.UserService.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Password != req.Password { // In production, use proper password comparison
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(result.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// handleRegister handles user registration
func (p *AuthPlugin) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a module with all dependencies
	module := fx.Module("auth",
		fx.Provide(
			// Provide user service
			fx.Annotate(
				NewUserService,
				fx.As(new(*UserService)),
			),
		),
	)

	// Create a new app with the module
	app := fx.New(module)

	// Start the app to initialize dependencies
	if err := app.Start(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start app"})
		return
	}
	defer app.Stop(context.Background())

	// Get user service from app
	var result struct {
		fx.In
		UserService *UserService
	}

	if err := fx.Populate(app, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dependencies"})
		return
	}

	user, err := result.UserService.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// handleProfile handles user profile retrieval
func (p *AuthPlugin) handleProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	username := c.GetString("username")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
		},
	})
}

// authMiddleware is a middleware that validates JWT tokens
func (p *AuthPlugin) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Create a module with all dependencies
		module := fx.Module("auth",
			fx.Provide(
				// Provide JWT secret
				fx.Annotate(
					func() []byte {
						return []byte("your-secret-key") // In production, use environment variable
					},
					fx.As(new([]byte)),
				),
			),
		)

		// Create a new app with the module
		app := fx.New(module)

		// Start the app to initialize dependencies
		if err := app.Start(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start app"})
			c.Abort()
			return
		}
		defer app.Stop(context.Background())

		// Get JWT secret from app
		var result struct {
			fx.In
			Secret []byte
		}

		if err := fx.Populate(app, &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get JWT secret"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return result.Secret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
