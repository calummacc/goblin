package user

import (
	"net/http"

	"github.com/calummacc/goblin/internal/common/enums"
	goblin "github.com/calummacc/goblin/internal/core"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

// Routes define all route of user controller
func (ctrl *UserController) Routes() []goblin.Route {
	return []goblin.Route{
		{
			Prefix: "/users",
			RouteConfig: []goblin.RouteConfig{
				{
					Method:  enums.GET,
					Path:    "",
					Handler: ctrl.GetUsers,
				},
				{
					Method:  enums.GET,
					Path:    "/:id",
					Handler: ctrl.GetUserByID,
				},
				{
					Method:  enums.POST,
					Path:    "",
					Handler: ctrl.CreateUser,
				},
				{
					Method:  enums.PROPFIND,
					Path:    "/test",
					Handler: ctrl.Test,
				},
			},
		},
	}
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users := ctrl.service.FindAll()
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	// Add error handling for invalid IDs
	user, err := ctrl.service.FindOne(id) // You'll need to implement FindOne in service.go
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // Improved error handling
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	//Implement CreateUser in service.go
	// Example:
	// var newUser struct {
	// 	Name string `json:"name"`
	// }
	// if err := c.ShouldBindJSON(&newUser); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// err := ctrl.service.Create(newUser.Name) //Implement service function
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

	c.JSON(http.StatusOK, gin.H{"data": "create user"}) //Temporary
}

func (ctrl *UserController) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "test propfind"})
}
