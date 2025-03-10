package user

import (
	"time"

	goblin "github.com/calummacc/goblin/internal/core"
	"github.com/calummacc/goblin/examples/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module struct
type Module struct {
	Middlewares []gin.HandlerFunc
}

// NewModule function to create a new Module
func NewModule() *Module {
	return &Module{
		Middlewares: []gin.HandlerFunc{
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(5, time.Minute),
		},
	}
}

func (m *Module) Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			NewUserController,
			NewUserService,
			func() *Module { return m}, //Provide *Module to Fx
		),
		fx.Provide(func(userController *UserController) goblin.Controller {
			return userController
		}),
		fx.Invoke(func(userController *UserController, module *Module) { //inject *Module
			for _, r := range userController.Routes() {
				r.Middlewares = append(module.Middlewares, r.Middlewares...)
			}
		}),
	)
}
