// goblin/examples/user_module/user_controller.go
package user_module

import (
	"net/http"

	"goblin/http"
	"goblin/middleware"
	"goblin/router"

	"github.com/gin-gonic/gin"
)

// User represents a user in the system
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserController manages user-related endpoints
type UserController struct {
	service *UserService
}

// NewUserController creates a new user controller
func NewUserController(service *UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// Routes returns the routes for this controller
func (c *UserController) Routes() []router.Route {
	builder := router.NewControllerBuilder("/users")

	builder.Get("/", c.GetUsers)
	builder.Get("/:id", c.GetUser)
	builder.Post("/", c.CreateUser)
	builder.Put("/:id", c.UpdateUser, middleware.Auth())
	builder.Delete("/:id", c.DeleteUser, middleware.Auth())

	routes := builder.Build().Routes
	return routes
}

// BasePath returns the base path for this controller
func (c *UserController) BasePath() string {
	return "/users"
}

// GetUsers returns all users
func (c *UserController) GetUsers(ctx *gin.Context) {
	response := http.NewResponse(ctx)
	users, err := c.service.GetUsers()
	if err != nil {
		response.Error(http.StatusInternalServerError, "Failed to get users")
		return
	}
	response.Success(users)
}

// GetUser returns a user by ID
func (c *UserController) GetUser(ctx *gin.Context) {
	response := http.NewResponse(ctx)
	request := http.NewRequest(ctx)

	id := request.GetParam("id")
	user, err := c.service.GetUser(id)
	if err != nil {
		response.NotFound("User not found")
		return
	}
	response.Success(user)
}

// CreateUser creates a new user
func (c *UserController) CreateUser(ctx *gin.Context) {
	response := http.NewResponse(ctx)
	request := http.NewRequest(ctx)

	var user User
	if err := request.GetBody(&user); err != nil {
		response.BadRequest("Invalid user data")
		return
	}

	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		response.Error(http.StatusInternalServerError, "Failed to create user")
		return
	}
	response.Created(createdUser)
}

// UpdateUser updates a user
func (c *UserController) UpdateUser(ctx *gin.Context) {
	response := http.NewResponse(ctx)
	request := http.NewRequest(ctx)

	id := request.GetParam("id")
	var user User
	if err := request.GetBody(&user); err != nil {
		response.BadRequest("Invalid user data")
		return
	}

	updatedUser, err := c.service.UpdateUser(id, user)
	if err != nil {
		response.Error(http.StatusInternalServerError, "Failed to update user")
		return
	}
	response.Success(updatedUser)
}

// DeleteUser deletes a user
func (c *UserController) DeleteUser(ctx *gin.Context) {
	response := http.NewResponse(ctx)
	request := http.NewRequest(ctx)

	id := request.GetParam("id")
	if err := c.service.DeleteUser(id); err != nil {
		response.Error(http.StatusInternalServerError, "Failed to delete user")
		return
	}
	response.NoContent()
}
