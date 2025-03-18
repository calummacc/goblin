package main

import (
	"fmt"
	"log"

	"goblin/core"

	"github.com/gin-gonic/gin"
)

// AuthGuard là một guard kiểm tra authentication
type AuthGuard struct {
	core.BaseGuard
}

// CanActivate kiểm tra token trong header
func (g AuthGuard) CanActivate(ctx *core.Context) (bool, error) {
	token := ctx.GinContext.GetHeader("Authorization")
	if token == "" {
		return false, fmt.Errorf("missing authorization token")
	}
	return true, nil
}

// LoggerMiddleware là một middleware ghi log
func LoggerMiddleware(c *gin.Context) {
	log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
	c.Next()
}

// User đại diện cho một user
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserController xử lý các request liên quan đến user
type UserController struct {
	*core.BaseController
}

// NewUserController tạo một UserController mới
func NewUserController() *UserController {
	ctrl := &UserController{}
	core.SetController("/users")(ctrl)
	ctrl.UseMiddleware(LoggerMiddleware)
	ctrl.UseGuards(&AuthGuard{})
	return ctrl
}

// GetUsers trả về danh sách users
func (c *UserController) GetUsers(ctx *gin.Context) {
	users := []User{
		{ID: 1, Name: "User 1"},
		{ID: 2, Name: "User 2"},
	}
	ctx.JSON(200, users)
}

// GetUser trả về thông tin một user
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user := User{ID: 1, Name: fmt.Sprintf("User %s", id)}
	ctx.JSON(200, user)
}

// CreateUser tạo một user mới
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, user)
}

// UpdateUser cập nhật thông tin user
func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user.ID = 1 // Giả lập ID
	ctx.JSON(200, user)
}

// DeleteUser xóa một user
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("User %s deleted successfully", id),
	})
}

func init() {
	// Đăng ký các routes
	controller := &UserController{}

	core.Get("/")(controller, "GetUsers")
	core.Get("/:id")(controller, "GetUser")
	core.Post("/")(controller, "CreateUser")
	core.Put("/:id")(controller, "UpdateUser")
	core.Delete("/:id")(controller, "DeleteUser")
}

func main() {
	// Tạo Goblin app
	app := core.NewGoblinApp(core.GoblinAppOptions{
		Debug: true,
		Modules: []core.Module{
			NewUserModule(),
		},
	})

	// Start app
	if err := app.Start(":8080"); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

// UserModule là module chứa UserController
type UserModule struct {
	*core.BaseModule
}

// NewUserModule tạo một UserModule mới
func NewUserModule() *UserModule {
	module := &UserModule{}
	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			NewUserController,
		},
		Controllers: []interface{}{
			func(cm *core.ControllerManager) error {
				return cm.RegisterController(NewUserController())
			},
		},
	})
	return module
}
