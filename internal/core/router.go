package core

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// NewGinEngine provides a *gin.Engine to the fx container.
func NewGinEngine() *gin.Engine {
	// gin.New() creates a new engine without any default middleware
	// gin.Default() creates an engine with some default middlewares
	router := gin.New()

	// For this example, let's add gin.Recovery() to handle panics
	router.Use(gin.Recovery())

	return router
}

// RouterModule wraps router-related providers or invocations
var RouterModule = fx.Options(
	fx.Provide(NewGinEngine),
	// You can provide more router-related things, e.g. middlewares, later
)
