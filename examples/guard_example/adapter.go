package main

import (
	"goblin/core"
	"goblin/guard"

	"github.com/gin-gonic/gin"
)

// GuardAdapter là một adapter chuyển đổi giữa guard.Guard và core.Guard
type GuardAdapter struct {
	guard guard.Guard
}

// NewGuardAdapter tạo một GuardAdapter mới
func NewGuardAdapter(g guard.Guard) *GuardAdapter {
	return &GuardAdapter{
		guard: g,
	}
}

// CanActivate triển khai core.Guard.CanActivate
func (a *GuardAdapter) CanActivate(ctx *core.Context) (bool, error) {
	// Chuyển đổi từ core.Context sang guard.GuardContext
	guardCtx := guard.NewGuardContext(ctx.GinContext, ctx.Handler, ctx.Class)

	// Gọi CanActivate của guard gốc
	return a.guard.CanActivate(guardCtx)
}

// CoreToGuardAdapter là một adapter chuyển đổi từ core.Guard sang guard.Guard
type CoreToGuardAdapter struct {
	guard core.Guard
}

// NewCoreToGuardAdapter tạo một CoreToGuardAdapter mới
func NewCoreToGuardAdapter(g core.Guard) *CoreToGuardAdapter {
	return &CoreToGuardAdapter{
		guard: g,
	}
}

// CanActivate triển khai guard.Guard.CanActivate
func (a *CoreToGuardAdapter) CanActivate(ctx *guard.GuardContext) (bool, error) {
	// Chuyển đổi từ guard.GuardContext sang core.Context
	coreCtx := &core.Context{
		GinContext: ctx.GinContext,
		Handler:    ctx.Handler,
		Class:      ctx.Controller,
	}

	// Gọi CanActivate của core.Guard gốc
	return a.guard.CanActivate(coreCtx)
}

// CreateGuardMiddleware tạo một gin.HandlerFunc từ guard.Guard
func CreateGuardMiddleware(g guard.Guard) gin.HandlerFunc {
	return func(c *gin.Context) {
		guardCtx := guard.NewGuardContext(c, nil, nil)
		allowed, err := g.CanActivate(guardCtx)

		if !allowed {
			if err != nil {
				c.AbortWithError(403, err)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}

		// Nếu guard đã thiết lập user, truyền vào context
		if guardCtx.User != nil {
			c.Set("user", guardCtx.User)
		}

		c.Next()
	}
}
