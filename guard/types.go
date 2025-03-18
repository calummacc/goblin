package guard

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Một số error thường gặp
var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidToken   = errors.New("invalid token")
	ErrTokenExpired   = errors.New("token expired")
	ErrMissingToken   = errors.New("missing authentication token")
	ErrInvalidRequest = errors.New("invalid request")
)

// GuardContext chứa thông tin về request context cho guards
type GuardContext struct {
	// Context Gin gốc
	GinContext *gin.Context
	// Handler là route handler gốc
	Handler interface{}
	// Controller là controller chứa handler
	Controller interface{}
	// Path là đường dẫn của route
	Path string
	// Method là HTTP method
	Method string
	// Dữ liệu tùy chỉnh có thể được chia sẻ giữa các guards
	Data map[string]interface{}
	// Request là HTTP request gốc
	Request *http.Request
	// User chứa thông tin về người dùng đã xác thực (nếu có)
	User interface{}
}

// NewGuardContext tạo một GuardContext mới
func NewGuardContext(c *gin.Context, handler interface{}, controller interface{}) *GuardContext {
	return &GuardContext{
		GinContext: c,
		Handler:    handler,
		Controller: controller,
		Path:       c.FullPath(),
		Method:     c.Request.Method,
		Data:       make(map[string]interface{}),
		Request:    c.Request,
	}
}

// Guard định nghĩa interface cho guards
type Guard interface {
	// CanActivate kiểm tra xem request có được phép tiếp tục không
	// Trả về true nếu request được phép tiếp tục, false nếu bị từ chối
	// Có thể trả về error để cung cấp thêm thông tin về lý do từ chối
	CanActivate(ctx *GuardContext) (bool, error)
}

// LifecycleGuard mở rộng Guard với các hook cho vòng đời
type LifecycleGuard interface {
	Guard
	// OnRegister được gọi khi guard được đăng ký
	OnRegister(ctx context.Context) error
	// OnShutdown được gọi khi ứng dụng đang shutdown
	OnShutdown(ctx context.Context) error
}

// BaseGuard cung cấp implementation mặc định cho Guard
type BaseGuard struct{}

// CanActivate implementation mặc định cho phép tất cả requests
func (g *BaseGuard) CanActivate(ctx *GuardContext) (bool, error) {
	return true, nil
}

// ConfigurableGuard mở rộng BaseGuard với các tùy chọn cấu hình
type ConfigurableGuard struct {
	BaseGuard
	// Options chứa các tùy chọn cấu hình cho guard
	Options GuardOptions
}

// GuardOptions chứa các tùy chọn cấu hình cho guard
type GuardOptions struct {
	// ErrorStatus là HTTP status code khi guard từ chối request
	ErrorStatus int
	// ErrorMessage là thông báo lỗi khi guard từ chối request
	ErrorMessage string
	// SkipPaths là danh sách các đường dẫn được bỏ qua kiểm tra
	SkipPaths []string
	// OnlyPaths là danh sách các đường dẫn được áp dụng kiểm tra
	OnlyPaths []string
}

// DefaultGuardOptions trả về các tùy chọn mặc định cho guard
func DefaultGuardOptions() GuardOptions {
	return GuardOptions{
		ErrorStatus:  http.StatusForbidden,
		ErrorMessage: "Forbidden",
		SkipPaths:    make([]string, 0),
		OnlyPaths:    make([]string, 0),
	}
}

// NewConfigurableGuard tạo một ConfigurableGuard mới với các tùy chọn
func NewConfigurableGuard(options GuardOptions) *ConfigurableGuard {
	return &ConfigurableGuard{
		Options: options,
	}
}

// ShouldSkip kiểm tra xem path có nên bỏ qua hay không
func (g *ConfigurableGuard) ShouldSkip(path string) bool {
	// Nếu có OnlyPaths và path không nằm trong đó, thì skip
	if len(g.Options.OnlyPaths) > 0 {
		for _, p := range g.Options.OnlyPaths {
			if p == path {
				return false
			}
		}
		return true
	}

	// Nếu path nằm trong SkipPaths, thì skip
	for _, p := range g.Options.SkipPaths {
		if p == path {
			return true
		}
	}

	return false
}

// GuardManager quản lý các guards trong ứng dụng
type GuardManager struct {
	// guards là danh sách các guards đã đăng ký
	guards []Guard
}

// NewGuardManager tạo một GuardManager mới
func NewGuardManager() *GuardManager {
	return &GuardManager{
		guards: make([]Guard, 0),
	}
}

// RegisterGuard đăng ký một guard với manager
func (m *GuardManager) RegisterGuard(guard Guard) {
	m.guards = append(m.guards, guard)
}

// CreateGuardMiddleware tạo middleware từ một guard
func CreateGuardMiddleware(guard Guard) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewGuardContext(c, nil, nil)

		if ok, err := guard.CanActivate(ctx); !ok {
			if err != nil {
				c.AbortWithError(http.StatusForbidden, err)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		// Nếu guard đã thiết lập user, truyền vào context
		if ctx.User != nil {
			c.Set("user", ctx.User)
		}

		c.Next()
	}
}

// IsAuthenticated là helper để kiểm tra xem request có được xác thực không
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user")
	return exists
}

// GetCurrentUser là helper để lấy thông tin user từ context
func GetCurrentUser(c *gin.Context) interface{} {
	user, _ := c.Get("user")
	return user
}
