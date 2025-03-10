package goblin

import (
	"github.com/gin-gonic/gin"
)

// Controller interface
type Controller interface {
	Routes() []Route // Trả về danh sách các route mà controller này quản lý.
}

// BaseController struct
type BaseController struct{}

// SendJSON is a helper function to send JSON responses
func (c BaseController) SendJSON(ctx *gin.Context, status int, data interface{}) {
	ctx.JSON(status, data)
}

// SendError is a helper function to send error responses
func (c BaseController) SendError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}
