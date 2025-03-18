// goblin/examples/user_module/user_controller.go
package user_module

import (
	"goblin/core"
	ghttp "goblin/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

// User represents a user in the system
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserController handles user-related endpoints
// @SetController("/users")
type UserController struct {
	core.BaseController
	service *UserService
}

// NewUserController creates a new user controller
func NewUserController(service *UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// GetUsers returns a list of users
// @Get("/")
func (c *UserController) GetUsers(ctx *gin.Context) {
	response := ghttp.NewResponse(ctx)
	users, err := c.service.GetUsers()
	if err != nil {
		response.Error(http.StatusInternalServerError, "Failed to get users")
		return
	}
	response.Success(users)
}

// GetUser returns a single user by ID
// @Get("/:id")
func (c *UserController) GetUser(ctx *gin.Context) {
	response := ghttp.NewResponse(ctx)
	request := ghttp.NewRequest(ctx)

	id := request.GetParam("id")
	user, err := c.service.GetUser(id)
	if err != nil {
		response.NotFound("User not found")
		return
	}
	response.Success(user)
}

// CreateUser creates a new user
// @Post("/")
func (c *UserController) CreateUser(ctx *gin.Context) {
	response := ghttp.NewResponse(ctx)
	request := ghttp.NewRequest(ctx)

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

// UpdateUser updates an existing user
// @Put("/:id")
// @UseMiddleware(middleware.Auth())
func (c *UserController) UpdateUser(ctx *gin.Context) {
	response := ghttp.NewResponse(ctx)
	request := ghttp.NewRequest(ctx)

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
// @Delete("/:id")
// @UseMiddleware(middleware.Auth())
func (c *UserController) DeleteUser(ctx *gin.Context) {
	response := ghttp.NewResponse(ctx)
	request := ghttp.NewRequest(ctx)

	id := request.GetParam("id")
	if err := c.service.DeleteUser(id); err != nil {
		response.Error(http.StatusInternalServerError, "Failed to delete user")
		return
	}
	response.NoContent()
}
