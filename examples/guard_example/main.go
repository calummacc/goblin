package main

import (
	"log"
	"net/http"
	"time"

	"goblin/core"
	"goblin/guard"

	"github.com/gin-gonic/gin"
)

// User đại diện cho một người dùng trong hệ thống
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// UserDatabase là một mock database đơn giản
type UserDatabase struct {
	users map[string]*User
}

// NewUserDatabase tạo một UserDatabase mới
func NewUserDatabase() *UserDatabase {
	return &UserDatabase{
		users: map[string]*User{
			"1": {ID: "1", Username: "admin", Role: "admin"},
			"2": {ID: "2", Username: "user", Role: "user"},
		},
	}
}

// FindByID tìm user theo ID
func (db *UserDatabase) FindByID(id string) *User {
	return db.users[id]
}

// FindAll trả về tất cả users
func (db *UserDatabase) FindAll() []*User {
	result := make([]*User, 0, len(db.users))
	for _, user := range db.users {
		result = append(result, user)
	}
	return result
}

// AuthService quản lý xác thực người dùng
type AuthService struct {
	jwtGuard *guard.JWTGuard
	db       *UserDatabase
}

// NewAuthService tạo một AuthService mới
func NewAuthService(db *UserDatabase) *AuthService {
	// Tạo JWTGuard với secret key và thời gian token hết hạn
	jwtGuard := guard.NewJWTGuard("super-secret-key", 1*time.Hour)

	// Cấu hình paths không cần xác thực
	jwtGuard.Options.SkipPaths = []string{"/auth/login", "/public"}

	return &AuthService{
		jwtGuard: jwtGuard,
		db:       db,
	}
}

// GenerateToken tạo token cho user
func (s *AuthService) GenerateToken(userID, username, role string) (string, error) {
	return s.jwtGuard.GenerateToken(userID, username, role)
}

// GetJWTGuard trả về JWTGuard
func (s *AuthService) GetJWTGuard() *guard.JWTGuard {
	return s.jwtGuard
}

// AdminGuard chỉ cho phép admins truy cập
type AdminGuard struct {
	core.BaseGuard
}

// CanActivate kiểm tra quyền admin
func (g *AdminGuard) CanActivate(ctx *core.Context) (bool, error) {
	// Lấy user từ context (được thiết lập bởi JWTGuard)
	if user, exists := ctx.GinContext.Get("user"); exists {
		if claims, ok := user.(*guard.JWTClaims); ok {
			// Kiểm tra role
			if claims.Role == "admin" {
				return true, nil
			}
			return false, guard.ErrForbidden
		}
	}

	return false, guard.ErrUnauthorized
}

// LoggerMiddleware ghi log cho mỗi request
func LoggerMiddleware(c *gin.Context) {
	startTime := time.Now()

	// Xử lý request
	c.Next()

	// Ghi log sau khi xử lý xong
	duration := time.Since(startTime)
	log.Printf("[%s] %s %s | %d | %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		c.Writer.Status(),
		duration,
	)
}

// AuthController xử lý các request liên quan đến xác thực
type AuthController struct {
	*core.BaseController
	authService *AuthService
}

// NewAuthController tạo một AuthController mới
func NewAuthController(authService *AuthService) *AuthController {
	ctrl := &AuthController{
		BaseController: &core.BaseController{},
		authService:    authService,
	}
	core.SetController("/auth")(ctrl)
	return ctrl
}

// Login xử lý yêu cầu đăng nhập
func (c *AuthController) Login(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid login data"})
		return
	}

	// Đơn giản hóa: không kiểm tra password cho ví dụ này
	var userID, role string
	if loginData.Username == "admin" {
		userID = "1"
		role = "admin"
	} else {
		userID = "2"
		role = "user"
	}

	// Tạo token
	token, err := c.authService.GenerateToken(userID, loginData.Username, role)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":       userID,
			"username": loginData.Username,
			"role":     role,
		},
	})
}

// AdminController xử lý các request chỉ dành cho admin
type AdminController struct {
	*core.BaseController
	db *UserDatabase
}

// NewAdminController tạo một AdminController mới
func NewAdminController(db *UserDatabase) *AdminController {
	ctrl := &AdminController{
		BaseController: &core.BaseController{},
		db:             db,
	}
	core.SetController("/admin")(ctrl)

	// Thêm AdminGuard cho controller này
	ctrl.UseGuards(&AdminGuard{})

	return ctrl
}

// GetUsers trả về danh sách users (chỉ admin mới xem được)
func (c *AdminController) GetUsers(ctx *gin.Context) {
	users := c.db.FindAll()
	ctx.JSON(200, users)
}

// UserController xử lý các request liên quan đến user
type UserController struct {
	*core.BaseController
	db *UserDatabase
}

// NewUserController tạo một UserController mới
func NewUserController(db *UserDatabase) *UserController {
	ctrl := &UserController{
		BaseController: &core.BaseController{},
		db:             db,
	}
	core.SetController("/users")(ctrl)
	return ctrl
}

// GetProfile trả về thông tin người dùng hiện tại
func (c *UserController) GetProfile(ctx *gin.Context) {
	// Lấy user từ context đã được thiết lập bởi JWTGuard
	claims, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	jwtClaims, _ := claims.(*guard.JWTClaims)
	user := c.db.FindByID(jwtClaims.UserID)

	ctx.JSON(200, user)
}

// PublicController xử lý các request công khai
type PublicController struct {
	*core.BaseController
}

// NewPublicController tạo một PublicController mới
func NewPublicController() *PublicController {
	ctrl := &PublicController{
		BaseController: &core.BaseController{},
	}
	core.SetController("/public")(ctrl)
	return ctrl
}

// GetInfo trả về thông tin công khai
func (c *PublicController) GetInfo(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"name":    "Goblin Framework",
		"version": "1.0.0",
		"author":  "Goblin Team",
	})
}

// AppModule là module chính của ứng dụng
type AppModule struct {
	*core.BaseModule
}

// NewAppModule tạo một AppModule mới
func NewAppModule() *AppModule {
	// Khởi tạo dependencies
	db := NewUserDatabase()
	authService := NewAuthService(db)

	// Khởi tạo controllers
	authController := NewAuthController(authService)
	adminController := NewAdminController(db)
	userController := NewUserController(db)
	publicController := NewPublicController()

	// Tạo module với các providers và controllers
	builder := core.NewModuleBuilder()
	module := &AppModule{}
	module.BaseModule = core.NewBaseModule(
		builder.
			Provide(db).
			Provide(authService).
			Controller(authController).
			Controller(adminController).
			Controller(userController).
			Controller(publicController).
			Build().GetMetadata(),
	)

	return module
}

func main() {
	// Tạo Gin engine
	engine := gin.Default()

	// Thêm middleware ghi log
	engine.Use(LoggerMiddleware)

	// Tạo ứng dụng Goblin
	app := core.NewGoblinApp(core.GoblinAppOptions{
		Debug: true,
	})

	// Khởi tạo module
	appModule := NewAppModule()

	// Đăng ký module
	app.RegisterModules(appModule)

	// Lấy AuthService để đăng ký JWTGuard
	providers := appModule.GetMetadata().Providers
	var authService *AuthService
	for _, provider := range providers {
		if service, ok := provider.(*AuthService); ok {
			authService = service
			break
		}
	}

	// Đăng ký JWTGuard làm middleware toàn cục
	if authService != nil {
		// Sử dụng adapter để tạo middleware từ guard.Guard
		jwtMiddleware := CreateGuardMiddleware(authService.GetJWTGuard())
		engine.Use(jwtMiddleware)
	}

	// Thiết lập các route cho controllers
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	engine.POST("/auth/login", func(c *gin.Context) {
		// Tìm authController và gọi Login
		for _, controller := range appModule.GetMetadata().Controllers {
			if ac, ok := controller.(*AuthController); ok {
				ac.Login(c)
				return
			}
		}
	})

	engine.GET("/public", func(c *gin.Context) {
		// Tìm publicController và gọi GetInfo
		for _, controller := range appModule.GetMetadata().Controllers {
			if pc, ok := controller.(*PublicController); ok {
				pc.GetInfo(c)
				return
			}
		}
	})

	engine.GET("/users/profile", func(c *gin.Context) {
		// Tìm userController và gọi GetProfile
		for _, controller := range appModule.GetMetadata().Controllers {
			if uc, ok := controller.(*UserController); ok {
				uc.GetProfile(c)
				return
			}
		}
	})

	engine.GET("/admin/users", func(c *gin.Context) {
		// Tìm adminController và gọi GetUsers
		for _, controller := range appModule.GetMetadata().Controllers {
			if ac, ok := controller.(*AdminController); ok {
				ac.GetUsers(c)
				return
			}
		}
	})

	// Khởi động ứng dụng
	log.Println("Starting Goblin Guard Example...")
	log.Println("Server is running on :8080")

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
