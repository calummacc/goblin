package user

import (
	"github.com/calummacc/goblin/internal/core"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type UserModule struct {
	core.BaseModule
	controller *Controller
	service    Service
	repository Repository
}

func NewUserModule() *UserModule {
	return &UserModule{}
}

func (m *UserModule) Configure(container *core.Container) {
	m.repository = NewRepository()
	m.service = NewService(m.repository)
	m.controller = NewController(m.service)
}

func (m *UserModule) ProvideDependencies() fx.Option {
	return fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewController,
		),
	)
}

func (m *UserModule) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("", m.controller.GetUsers)
		users.GET("/:id", m.controller.GetUser)
		users.POST("", m.controller.CreateUser)
		users.PUT("/:id", m.controller.UpdateUser)
		users.DELETE("/:id", m.controller.DeleteUser)
	}
}
