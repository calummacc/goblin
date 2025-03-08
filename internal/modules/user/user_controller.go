package user

import "github.com/gin-gonic/gin"

type UserController struct{}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{}
}

// Example endpoint
func (uc UserController) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User endpoint working!"})
}
