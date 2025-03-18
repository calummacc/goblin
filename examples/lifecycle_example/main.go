package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"goblin/core"

	"github.com/gin-gonic/gin"
)

// DatabaseService là một service sử dụng lifecycle hooks
type DatabaseService struct {
	isConnected bool
}

// NewDatabaseService tạo một DatabaseService mới
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// OnModuleInit được gọi khi module khởi tạo
func (s *DatabaseService) OnModuleInit(ctx context.Context) error {
	log.Println("DatabaseService: OnModuleInit - Connecting to database...")
	time.Sleep(500 * time.Millisecond) // Giả lập kết nối
	s.isConnected = true
	log.Println("DatabaseService: Database connected successfully")
	return nil
}

// OnApplicationBootstrap được gọi khi ứng dụng khởi động
func (s *DatabaseService) OnApplicationBootstrap(ctx context.Context) error {
	log.Println("DatabaseService: OnApplicationBootstrap - Running migrations...")
	time.Sleep(200 * time.Millisecond) // Giả lập migrations
	log.Println("DatabaseService: Migrations completed successfully")
	return nil
}

// OnApplicationShutdown được gọi khi ứng dụng shutdown
func (s *DatabaseService) OnApplicationShutdown(ctx context.Context) error {
	log.Println("DatabaseService: OnApplicationShutdown - Closing all active connections...")
	time.Sleep(300 * time.Millisecond) // Giả lập đóng kết nối
	log.Println("DatabaseService: All active connections closed")
	return nil
}

// OnModuleDestroy được gọi khi module bị hủy
func (s *DatabaseService) OnModuleDestroy(ctx context.Context) error {
	log.Println("DatabaseService: OnModuleDestroy - Releasing database resources...")
	time.Sleep(100 * time.Millisecond) // Giả lập giải phóng tài nguyên
	s.isConnected = false
	log.Println("DatabaseService: Database resources released")
	return nil
}

// IsConnected kiểm tra trạng thái kết nối
func (s *DatabaseService) IsConnected() bool {
	return s.isConnected
}

// CacheService là một service khác sử dụng lifecycle hooks
type CacheService struct {
	isInitialized bool
}

// NewCacheService tạo một CacheService mới
func NewCacheService() *CacheService {
	return &CacheService{}
}

// OnModuleInit được gọi khi module khởi tạo
func (s *CacheService) OnModuleInit(ctx context.Context) error {
	log.Println("CacheService: OnModuleInit - Initializing cache...")
	time.Sleep(200 * time.Millisecond) // Giả lập khởi tạo
	s.isInitialized = true
	log.Println("CacheService: Cache initialized")
	return nil
}

// OnApplicationShutdown được gọi khi ứng dụng shutdown
func (s *CacheService) OnApplicationShutdown(ctx context.Context) error {
	log.Println("CacheService: OnApplicationShutdown - Flushing cache...")
	time.Sleep(100 * time.Millisecond) // Giả lập flush cache
	log.Println("CacheService: Cache flushed")
	return nil
}

// AppController là controller của ứng dụng
type AppController struct {
	db    *DatabaseService
	cache *CacheService
}

// NewAppController tạo một AppController mới
func NewAppController(db *DatabaseService, cache *CacheService) *AppController {
	return &AppController{
		db:    db,
		cache: cache,
	}
}

// RegisterRoutes đăng ký các routes
func (c *AppController) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/status", c.getStatus)
}

// getStatus trả về trạng thái của services
func (c *AppController) getStatus(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"database_connected": c.db.IsConnected(),
		"cache_initialized":  c.cache.isInitialized,
		"timestamp":          time.Now().Format(time.RFC3339),
	})
}

// LifecycleModule là module sử dụng lifecycle hooks
type LifecycleModule struct {
	*core.BaseModule
}

// NewLifecycleModule tạo một LifecycleModule mới
func NewLifecycleModule() *LifecycleModule {
	module := &LifecycleModule{}

	module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
		Providers: []interface{}{
			NewDatabaseService,
			NewCacheService,
			NewAppController,
		},
		Controllers: []interface{}{
			func(controller *AppController, engine *gin.Engine) {
				controller.RegisterRoutes(engine)
			},
		},
	})

	return module
}

// OnModuleInit được gọi khi module khởi tạo
func (m *LifecycleModule) OnModuleInit(ctx context.Context) error {
	log.Println("LifecycleModule: OnModuleInit - Initializing module...")
	return nil
}

// OnApplicationBootstrap được gọi khi ứng dụng khởi động
func (m *LifecycleModule) OnApplicationBootstrap(ctx context.Context) error {
	log.Println("LifecycleModule: OnApplicationBootstrap - Module fully initialized")
	return nil
}

// OnApplicationShutdown được gọi khi ứng dụng shutdown
func (m *LifecycleModule) OnApplicationShutdown(ctx context.Context) error {
	log.Println("LifecycleModule: OnApplicationShutdown - Beginning module shutdown")
	return nil
}

// OnModuleDestroy được gọi khi module bị hủy
func (m *LifecycleModule) OnModuleDestroy(ctx context.Context) error {
	log.Println("LifecycleModule: OnModuleDestroy - Module destroyed")
	return nil
}

func main() {
	// Tạo Goblin app
	app := core.NewGoblinApp(core.GoblinAppOptions{
		Debug: true,
		Modules: []core.Module{
			NewLifecycleModule(),
		},
	})

	// Đăng ký shutdown hook
	app.RegisterShutdownHook(func(ctx context.Context) error {
		log.Println("Custom shutdown hook: Performing final cleanup...")
		time.Sleep(100 * time.Millisecond) // Giả lập cleanup
		log.Println("Custom shutdown hook: Cleanup completed")
		return nil
	})

	// Khởi động ứng dụng (port 8080)
	log.Println("Starting Goblin app with lifecycle hooks example...")

	fmt.Println("\n==== LIFECYCLE SEQUENCE ====")
	fmt.Println("1. ModuleInit: LifecycleModule → DatabaseService → CacheService")
	fmt.Println("2. ApplicationBootstrap: LifecycleModule → DatabaseService")
	fmt.Println("3. Application Running")
	fmt.Println("4. ApplicationShutdown: LifecycleModule → DatabaseService → CacheService → Custom hooks")
	fmt.Println("5. ModuleDestroy: LifecycleModule → DatabaseService")
	fmt.Println("============================\n")

	fmt.Println("Press Ctrl+C to trigger graceful shutdown\n")

	if err := app.Start(":8080"); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
