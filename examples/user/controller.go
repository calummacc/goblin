package user

import (
	"net/http"

	goblin "github.com/calummacc/goblin/internal/core"
	"github.com/calummacc/goblin/examples/middlewares"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	goblin.BaseController
	UserService *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{UserService: service}
}

// Routes define all route of user controller
func (ctrl *UserController) Routes() []goblin.Route {
	return []goblin.Route{
		{
			Prefix: "/users",
			RouteConfig: []goblin.RouteConfig{
				{
					Method:  http.MethodGet,
					Path:    "",
					Handler: ctrl.GetUsers,
				},
				{
					Method:      http.MethodGet,
					Path:        "/:id",
					Handler:     ctrl.GetUserByID,
					Middlewares: []gin.HandlerFunc{middlewares.LoggerMiddleware()},
				},
				{
					Method:  http.MethodPost,
					Path:    "",
					Handler: ctrl.CreateUser,
				},
				{
					Method:  http.MethodPut,
					Path:    "/:id",
					Handler: ctrl.CreateUser,
				},
				{
					Method:  http.MethodDelete,
					Path:    "/test",
					Handler: ctrl.Test,
				},
			},
		},
	}
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users := ctrl.UserService.FindAll()
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := ctrl.UserService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var newUser struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ctrl.UserService.Create(newUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ctrl *UserController) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "test propfind"})
}

