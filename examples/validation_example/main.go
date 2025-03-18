package validation_example

import (
	"goblin/validation"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// User represents a user entity
type User struct {
	ID          string    `json:"id" validate:"required"`
	Username    string    `json:"username" validate:"required|min:3|max:20" msg:"Username must be between 3 and 20 characters"`
	Email       string    `json:"email" validate:"required|email" msg:"Invalid email format"`
	Password    string    `json:"password" validate:"required" msg:"Password is required"`
	PhoneNumber string    `json:"phone" validate:"required" msg:"Phone number is required"`
	BirthDate   string    `json:"birthDate" validate:"required" msg:"Birth date is required"`
	Website     string    `json:"website" validate:"required" msg:"Website is required"`
	Avatar      string    `json:"avatar" validate:"required" msg:"Avatar is required"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserController handles user-related requests
type UserController struct {
	validationPipe *validation.ValidationPipe
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	pipe := validation.NewValidationPipe(true)

	// Register custom validators
	pipe.RegisterValidator("password", validation.NewPasswordValidator(8, true, true, true, true))
	pipe.RegisterValidator("phone", validation.NewPhoneNumberValidator("+84", `^\+84\d{9}$`))
	pipe.RegisterValidator("birthDate", validation.NewDateValidator("2006-01-02", time.Time{}, time.Now()))
	pipe.RegisterValidator("website", validation.NewURLValidator(true, []string{"example.com", "goblin.dev"}))
	pipe.RegisterValidator("avatar", validation.NewFileValidator(5*1024*1024, []string{"image/jpeg", "image/png"}, true))

	return &UserController{
		validationPipe: pipe,
	}
}

// CreateUser handles user creation
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user data
	_, err := c.validationPipe.Transform(user, validation.ValidationMetadata{
		Field:      "user",
		Rules:      make(map[string]interface{}),
		Messages:   make(map[string]string),
		StopOnFail: true,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set creation time
	user.CreatedAt = time.Now()

	// Return success response
	ctx.JSON(http.StatusCreated, user)
}

func main() {
	// Create Gin router
	router := gin.Default()

	// Create user controller
	userController := NewUserController()

	// Define routes
	router.POST("/users", userController.CreateUser)

	// Start server
	log.Println("Server starting on http://localhost:8080")
	router.Run(":8080")
}
