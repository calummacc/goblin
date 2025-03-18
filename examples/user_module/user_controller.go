// goblin/examples/user_module/user_controller.go
package user_module

import (
	"goblin/core"
	"goblin/errors"
	ghttp "goblin/http"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User represents a user in the system
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserController handles HTTP requests for user operations
type UserController struct {
	*core.BaseController
	service *UserService
}

// NewUserController creates a new user controller
func NewUserController(service *UserService) *UserController {
	return &UserController{
		BaseController: core.NewBaseController(),
		service:        service,
	}
}

// RegisterRoutes registers all user routes
func (c *UserController) RegisterRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.GET("", c.GetUsers)
		users.GET("/:id", c.GetUser)
		users.POST("", c.CreateUser)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
	}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, users)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	response := ghttp.NewResponse(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx.JSON(200, gin.H{"message": "GetUser" + strconv.FormatUint(uint64(id), 10)})
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	ctx.JSON(201, gin.H{"message": "CreateUser"})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.Error(errors.NewValidationError("Invalid user ID"))
		return
	}

	ctx.JSON(200, gin.H{"message": "UpdateUser" + strconv.FormatUint(uint64(id), 10)})

}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.Error(errors.NewValidationError("Invalid user ID"))
		return
	}

	ctx.JSON(200, gin.H{"message": "DeleteUser" + strconv.FormatUint(uint64(id), 10)})
}
