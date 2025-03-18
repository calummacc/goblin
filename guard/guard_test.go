package guard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockGuard là một guard giả sử dụng cho kiểm thử
type MockGuard struct {
	ConfigurableGuard
	shouldAllow bool
	err         error
}

// CanActivate implementation cho MockGuard
func (g *MockGuard) CanActivate(ctx *GuardContext) (bool, error) {
	return g.shouldAllow, g.err
}

func TestGuardContext(t *testing.T) {
	// Tạo Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Tạo GuardContext
	ctx := NewGuardContext(c, nil, nil)

	// Kiểm tra các trường của GuardContext
	assert.NotNil(t, ctx.GinContext)
	assert.Equal(t, "/test", ctx.Path)
	assert.Equal(t, "GET", ctx.Method)
	assert.NotNil(t, ctx.Request)
	assert.NotNil(t, ctx.Data)
}

func TestBaseGuard(t *testing.T) {
	guard := &BaseGuard{}
	ctx := &GuardContext{}

	// Kiểm tra BaseGuard luôn cho phép request
	allowed, err := guard.CanActivate(ctx)
	assert.True(t, allowed)
	assert.NoError(t, err)
}

func TestConfigurableGuard(t *testing.T) {
	// Tạo guard với tùy chọn mặc định
	guard := NewConfigurableGuard(DefaultGuardOptions())

	// Kiểm tra ShouldSkip với path không nằm trong SkipPaths
	assert.False(t, guard.ShouldSkip("/test"))

	// Tạo guard với SkipPaths
	guard = NewConfigurableGuard(GuardOptions{
		SkipPaths: []string{"/public", "/health"},
	})
	assert.True(t, guard.ShouldSkip("/public"))
	assert.True(t, guard.ShouldSkip("/health"))
	assert.False(t, guard.ShouldSkip("/private"))

	// Tạo guard với OnlyPaths
	guard = NewConfigurableGuard(GuardOptions{
		OnlyPaths: []string{"/admin", "/api"},
	})
	assert.True(t, guard.ShouldSkip("/public"))
	assert.False(t, guard.ShouldSkip("/admin"))
	assert.False(t, guard.ShouldSkip("/api"))
}

func TestGuardManager(t *testing.T) {
	manager := NewGuardManager()

	// Tạo mock guards
	guard1 := &MockGuard{shouldAllow: true}
	guard2 := &MockGuard{shouldAllow: false}

	// Đăng ký guards
	manager.RegisterGuard(guard1)
	manager.RegisterGuard(guard2)

	// Tạo middleware từ guard
	middleware := CreateGuardMiddleware(guard1)

	// Tạo test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Chạy middleware
	middleware(c)

	// Kiểm tra kết quả
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTGuard(t *testing.T) {
	// Tạo JWTGuard
	guard := NewJWTGuard("test-secret", 1*time.Hour)

	// Test GenerateToken
	token, err := guard.GenerateToken("user1", "testuser", "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test ValidateToken
	claims, err := guard.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user1", claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "admin", claims.Role)

	// Test CanActivate với token hợp lệ
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	ctx := NewGuardContext(c, nil, nil)
	allowed, err := guard.CanActivate(ctx)
	assert.True(t, allowed)
	assert.NoError(t, err)

	// Test CanActivate với token không hợp lệ
	c.Request.Header.Set("Authorization", "Bearer invalid-token")
	allowed, err = guard.CanActivate(ctx)
	assert.False(t, allowed)
	assert.Error(t, err)

	// Test CanActivate không có token
	c.Request.Header.Del("Authorization")
	allowed, err = guard.CanActivate(ctx)
	assert.False(t, allowed)
	assert.Equal(t, ErrMissingToken, err)
}

func TestJWTGuardExpiration(t *testing.T) {
	// Tạo JWTGuard với token hết hạn sau 1 giây
	guard := NewJWTGuard("test-secret", 1*time.Second)

	// Tạo token
	token, err := guard.GenerateToken("user1", "testuser", "admin")
	assert.NoError(t, err)

	// Đợi token hết hạn
	time.Sleep(2 * time.Second)

	// Kiểm tra token đã hết hạn
	claims, err := guard.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTGuardSkipPaths(t *testing.T) {
	guard := NewJWTGuard("test-secret", 1*time.Hour)
	guard.Options.SkipPaths = []string{"/public", "/health"}

	// Test path được bỏ qua
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/public", nil)
	ctx := NewGuardContext(c, nil, nil)
	allowed, err := guard.CanActivate(ctx)
	assert.True(t, allowed)
	assert.NoError(t, err)

	// Test path không được bỏ qua
	c.Request = httptest.NewRequest("GET", "/private", nil)
	ctx = NewGuardContext(c, nil, nil)
	allowed, err = guard.CanActivate(ctx)
	assert.False(t, allowed)
	assert.Equal(t, ErrMissingToken, err)
}

func TestJWTGuardLifecycle(t *testing.T) {
	guard := NewJWTGuard("test-secret", 1*time.Hour)

	// Test OnRegister
	err := guard.OnRegister(context.Background())
	assert.NoError(t, err)

	// Test OnShutdown
	err = guard.OnShutdown(context.Background())
	assert.NoError(t, err)
}
