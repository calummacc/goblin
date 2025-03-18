package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"goblin/core"
	"goblin/pipe"

	"github.com/gin-gonic/gin"
)

// User là model dữ liệu người dùng
type User struct {
	ID        int    `json:"id" validate:"required,gt=0"`
	Name      string `json:"name" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Age       int    `json:"age" validate:"gte=18"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// UserController xử lý các request liên quan đến user
type UserController struct {
	core.BaseController
}

// NewUserController tạo một instance mới của UserController
func NewUserController() *UserController {
	controller := &UserController{}
	controller.SetMetadata(&core.ControllerMetadata{
		Prefix:     "/users",
		Routes:     make([]core.RouteMetadata, 0),
		Guards:     make([]core.Guard, 0),
		Middleware: make([]gin.HandlerFunc, 0),
	})
	return controller
}

// GetUser lấy thông tin một user theo ID
func (c *UserController) GetUser(ctx *gin.Context) {
	// Tạo validation pipe để validate ID
	idPipe := pipe.NewValidationPipe(pipe.ValidationOptions{
		Transform: true,
	})

	// Tạo parse pipe để chuyển đổi string sang int
	parsePipe := pipe.NewParseIntPipe(10)

	// Tạo composite pipe kết hợp validation và parse
	compositePipe := pipe.NewCompositePipe(
		pipe.DefaultPipeOptions(),
		idPipe,
		parsePipe,
	)

	// Lấy ID từ URL parameter
	id := ctx.Param("id")

	// Tạo context cho pipe
	pipeCtx := &pipe.TransformContext{
		Value: id,
		Type:  reflect.TypeOf(int(0)),
	}

	// Thực hiện transform
	result, err := compositePipe.Transform(pipeCtx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Chuyển đổi kết quả sang int
	userID := result.(int64)

	// TODO: Lấy thông tin user từ database
	user := User{
		ID:        int(userID),
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       25,
		IsActive:  true,
		CreatedAt: "2024-03-20",
	}

	ctx.JSON(http.StatusOK, user)
}

// CreateUser tạo một user mới
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user User

	// Tạo validation pipe với transform
	validationPipe := pipe.NewValidationPipe(pipe.ValidationOptions{
		Transform: true,
	})

	// Tạo trim pipe để xóa khoảng trắng thừa
	trimPipe := pipe.NewTrimPipe()

	// Tạo composite pipe kết hợp validation và trim
	compositePipe := pipe.NewCompositePipe(
		pipe.DefaultPipeOptions(),
		trimPipe,
		validationPipe,
	)

	// Lấy dữ liệu từ request body
	var input map[string]interface{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Tạo context cho pipe
	pipeCtx := &pipe.TransformContext{
		Value: input,
		Type:  reflect.TypeOf(user),
	}

	// Thực hiện transform
	result, err := compositePipe.Transform(pipeCtx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Chuyển đổi kết quả sang User
	user = *result.(*User)

	// TODO: Lưu user vào database
	user.CreatedAt = "2024-03-20"

	ctx.JSON(http.StatusCreated, user)
}

// UpdateUser cập nhật thông tin một user
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Tạo validation pipe để validate ID
	idPipe := pipe.NewValidationPipe(pipe.ValidationOptions{
		Transform: true,
	})

	// Tạo parse pipe để chuyển đổi string sang int
	parsePipe := pipe.NewParseIntPipe(10)

	// Tạo composite pipe kết hợp validation và parse
	compositePipe := pipe.NewCompositePipe(
		pipe.DefaultPipeOptions(),
		idPipe,
		parsePipe,
	)

	// Lấy ID từ URL parameter
	id := ctx.Param("id")

	// Tạo context cho pipe
	pipeCtx := &pipe.TransformContext{
		Value: id,
		Type:  reflect.TypeOf(int(0)),
	}

	// Thực hiện transform
	result, err := compositePipe.Transform(pipeCtx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Chuyển đổi kết quả sang int
	userID := result.(int64)

	var user User

	// Tạo validation pipe với transform
	validationPipe := pipe.NewValidationPipe(pipe.ValidationOptions{
		Transform: true,
	})

	// Tạo trim pipe để xóa khoảng trắng thừa
	trimPipe := pipe.NewTrimPipe()

	// Tạo composite pipe kết hợp validation và trim
	compositePipe = pipe.NewCompositePipe(
		pipe.DefaultPipeOptions(),
		trimPipe,
		validationPipe,
	)

	// Lấy dữ liệu từ request body
	var input map[string]interface{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Tạo context cho pipe
	pipeCtx = &pipe.TransformContext{
		Value: input,
		Type:  reflect.TypeOf(user),
	}

	// Thực hiện transform
	result, err = compositePipe.Transform(pipeCtx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Chuyển đổi kết quả sang User
	user = *result.(*User)

	// Đảm bảo ID khớp với ID trong URL
	user.ID = int(userID)

	// TODO: Cập nhật user trong database
	user.CreatedAt = "2024-03-20"

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser xóa một user
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Tạo validation pipe để validate ID
	idPipe := pipe.NewValidationPipe(pipe.ValidationOptions{
		Transform: true,
	})

	// Tạo parse pipe để chuyển đổi string sang int
	parsePipe := pipe.NewParseIntPipe(10)

	// Tạo composite pipe kết hợp validation và parse
	compositePipe := pipe.NewCompositePipe(
		pipe.DefaultPipeOptions(),
		idPipe,
		parsePipe,
	)

	// Lấy ID từ URL parameter
	id := ctx.Param("id")

	// Tạo context cho pipe
	pipeCtx := &pipe.TransformContext{
		Value: id,
		Type:  reflect.TypeOf(int(0)),
	}

	// Thực hiện transform
	result, err := compositePipe.Transform(pipeCtx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Chuyển đổi kết quả sang int
	userID := result.(int64)

	// TODO: Xóa user khỏi database

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %d deleted successfully", userID),
	})
}

// UserModule là module quản lý user
type UserModule struct {
	core.BaseModule
}

// NewUserModule tạo một instance mới của UserModule
func NewUserModule() *UserModule {
	module := &UserModule{}
	module.BaseModule = *core.NewBaseModule(core.ModuleMetadata{
		Imports:     make([]core.Module, 0),
		Exports:     make([]interface{}, 0),
		Providers:   make([]interface{}, 0),
		Controllers: make([]interface{}, 0),
	})
	return module
}

// RegisterControllers đăng ký các controller của module
func (m *UserModule) RegisterControllers(manager *core.ControllerManager) {
	userController := NewUserController()
	manager.RegisterController(userController)
}

func main() {
	// Tạo ứng dụng Goblin
	app := core.NewGoblinApp()

	// Đăng ký module
	app.RegisterModules(NewUserModule())

	// Khởi động ứng dụng
	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
