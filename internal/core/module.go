package core

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Module interface {
	Configure(container *Container)
}

type RouteModule interface {
	Module
	RegisterRoutes(router *gin.RouterGroup)
}

type FxModule interface {
	Module
	ProvideDependencies() fx.Option
}

type LifecycleModule interface {
	Module
	OnInit() error
	OnDestroy() error
}

type BaseModule struct{}

func (b *BaseModule) Configure(container *Container) {}
