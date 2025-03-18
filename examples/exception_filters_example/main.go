package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"goblin/core"
	"goblin/filter"

	"github.com/gin-gonic/gin"
)

// BusinessException represents an application-specific exception
type BusinessException struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *BusinessException) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// UserNotFoundException is a specialized business exception
type UserNotFoundException struct {
	UserID string
	BusinessException
}

// NewUserNotFoundException creates a new UserNotFoundException
func NewUserNotFoundException(userID string) *UserNotFoundException {
	return &UserNotFoundException{
		UserID: userID,
		BusinessException: BusinessException{
			Code:    "USER_NOT_FOUND",
			Message: fmt.Sprintf("User with ID %s was not found", userID),
		},
	}
}

// HttpExceptionFilter is a global filter to handle HTTP exceptions
type HttpExceptionFilter struct {
	filter.BaseExceptionFilter
}

// Catch handles exceptions by converting them to HTTP responses
func (f *HttpExceptionFilter) Catch(exception error, ctx *filter.ExceptionContext) {
	c := ctx.GinContext

	// Handle built-in HttpException
	if httpEx, ok := exception.(*filter.HttpException); ok {
		c.JSON(httpEx.StatusCode, httpEx)
		return
	}

	// Handle business exceptions
	if businessEx, ok := exception.(*BusinessException); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"error":      "Bad Request",
			"message":    businessEx.Message,
			"code":       businessEx.Code,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}

	// Handle user not found exception
	if userNotFoundEx, ok := exception.(*UserNotFoundException); ok {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"error":      "Not Found",
			"message":    userNotFoundEx.Message,
			"code":       userNotFoundEx.Code,
			"userId":     userNotFoundEx.UserID,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}

	// Default handling for unknown exceptions
	c.JSON(http.StatusInternalServerError, gin.H{
		"statusCode": http.StatusInternalServerError,
		"error":      "Internal Server Error",
		"message":    exception.Error(),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// CanHandle tells if this filter can handle the given exception
func (f *HttpExceptionFilter) CanHandle(exception error) bool {
	// This filter can handle all exceptions
	return true
}

// UserController handles user-related requests
type UserController struct {
	core.BaseController
}

// NewUserController creates a new UserController
func NewUserController() *UserController {
	controller := &UserController{}

	// Initialize metadata
	metadata := &core.ControllerMetadata{
		Prefix:     "/users",
		Routes:     make([]core.RouteMetadata, 0),
		Guards:     make([]core.Guard, 0),
		Middleware: make([]gin.HandlerFunc, 0),
	}

	// Set metadata directly
	controller.SetMetadata(metadata)

	return controller
}

// GetUser handles GET /users/:id
func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	// Simulate user lookup
	if userID == "404" {
		// Throw custom exception
		panic(NewUserNotFoundException(userID))
	}

	if userID == "400" {
		// Throw business exception
		panic(&BusinessException{
			Code:    "INVALID_USER_ID",
			Message: "The provided user ID is invalid",
		})
	}

	if userID == "500" {
		// Throw generic exception
		panic("Internal server error occurred during user lookup")
	}

	if userID == "403" {
		// Throw HTTP exception
		panic(filter.NewHttpException(http.StatusForbidden, "You don't have permission to access this user"))
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": fmt.Sprintf("user_%s", userID),
		"email":    fmt.Sprintf("user%s@example.com", userID),
	})
}

func main() {
	// Create Gin router
	r := gin.Default()

	// Create exception filter manager
	filterManager := filter.NewExceptionFilterManager()

	// Register global exception filter
	filterManager.RegisterGlobalFilter(&HttpExceptionFilter{})

	// Apply exception handler middleware
	r.Use(filterManager.GetExceptionHandlerMiddleware())

	// Create and register the user controller
	userController := NewUserController()

	// Register routes
	r.GET("/users/:id", userController.GetUser)

	// Add some informational routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Exception Filters Example",
			"routes": []string{
				"/users/1 - Returns a normal user",
				"/users/404 - Demonstrates UserNotFoundException",
				"/users/400 - Demonstrates BusinessException",
				"/users/500 - Demonstrates generic error handling",
				"/users/403 - Demonstrates HttpException",
			},
		})
	})

	// Start the server
	log.Println("Starting server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
