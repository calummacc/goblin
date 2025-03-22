package main

import (
	"github.com/calummacc/goblin/examples/basic/modules/user"
	"github.com/calummacc/goblin/internal/core"
	"github.com/calummacc/goblin/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AppModule struct {
	core.BaseModule
	userModule *user.UserModule
}

func NewAppModule() *AppModule {
	return &AppModule{
		userModule: user.NewUserModule(),
	}
}

func (m *AppModule) Configure(container *core.Container) {
	m.userModule.Configure(container)
}

func (m *AppModule) ProvideDependencies() fx.Option {
	return fx.Options(
		m.userModule.ProvideDependencies(),
	)
}

func (m *AppModule) RegisterRoutes(router *gin.RouterGroup) {
	// Add middleware
	router.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.ErrorHandler(),
	)

	// Register API routes
	api := router.Group("/api/v1")
	{
		m.userModule.RegisterRoutes(api)
	}
}
