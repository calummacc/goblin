package goblin

import (
	"github.com/gin-gonic/gin"
)

// Middleware interface
type Middleware interface {
	Handler() gin.HandlerFunc
}

// globalMiddleware struct for storing global middleware
type globalMiddleware struct {
	middlewares []gin.HandlerFunc
}

// newGlobalMiddleware creates a new globalMiddleware instance
func newGlobalMiddleware(middlewares ...gin.HandlerFunc) *globalMiddleware {
	return &globalMiddleware{middlewares: middlewares}
}

// getMiddlewares returns the slice of global middlewares
func (g *globalMiddleware) getMiddlewares() []gin.HandlerFunc {
	return g.middlewares
}

// addMiddlewares appends middlewares to the global middleware slice.
func (g *globalMiddleware) addMiddlewares(middlewares ...gin.HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
