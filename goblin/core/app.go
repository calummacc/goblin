// goblin/core/app.go
package core

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// GoblinApp đại diện cho một ứng dụng Goblin và nhúng Gin Engine
type GoblinApp struct {
	*gin.Engine // Nhúng Gin Engine để kế thừa tất cả phương thức của nó
	app         *fx.App
}

// GoblinAppOptions cấu hình ứng dụng Goblin
type GoblinAppOptions struct {
	Modules []GoblinModule
	Debug   bool
}

// NewGoblinApp tạo một ứng dụng Goblin mới (tương tự NestFactory.create())
func NewGoblinApp(opts ...GoblinAppOptions) *GoblinApp {
	var options GoblinAppOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// Chế độ mặc định là production mode
	if !options.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	if options.Debug {
		engine.Use(gin.Logger())
	}

	// Thu thập tất cả các module options
	var fxOptions []fx.Option
	for _, module := range options.Modules {
		fxOptions = append(fxOptions, module.Options)
	}

	// Thêm Gin engine vào DI container
	fxOptions = append(fxOptions,
		fx.Provide(func() *gin.Engine {
			return engine
		}),
	)

	goblinApp := &GoblinApp{
		Engine: engine,
		app:    fx.New(fxOptions...),
	}

	return goblinApp
}

// Start khởi động ứng dụng Goblin
func (g *GoblinApp) Start(port string) error {
	// Khởi động Fx app
	startCtx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	if err := g.app.Start(startCtx); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Khởi động Gin server
	log.Printf("Goblin app started on port %s", port)
	return g.Engine.Run(port) // Sử dụng g.Engine thay vì g.engine
}

// Stop dừng ứng dụng Goblin
func (g *GoblinApp) Stop() error {
	stopCtx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	return g.app.Stop(stopCtx)
}

// RegisterModules đăng ký thêm modules vào ứng dụng
func (g *GoblinApp) RegisterModules(modules ...GoblinModule) {
	// Đây là một triển khai đơn giản
	// Trong triển khai thực tế, bạn cần xử lý việc thêm modules sau khi ứng dụng đã được tạo
	log.Println("RegisterModules được gọi sau khi tạo ứng dụng - chức năng này chưa được triển khai đầy đủ")
}

// GetApp trả về ứng dụng Fx
func (g *GoblinApp) GetApp() *fx.App {
	return g.app
}
