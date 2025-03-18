package guard

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGuard là một guard kiểm tra JWT token
type JWTGuard struct {
	ConfigurableGuard
	// SecretKey là khóa bí mật để ký và xác thực token
	SecretKey []byte
	// TokenExpiration là thời gian token hết hạn
	TokenExpiration time.Duration
}

// JWTClaims định nghĩa các claims trong JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTGuard tạo một JWTGuard mới
func NewJWTGuard(secretKey string, tokenExpiration time.Duration) *JWTGuard {
	return &JWTGuard{
		ConfigurableGuard: *NewConfigurableGuard(GuardOptions{
			ErrorStatus:  http.StatusUnauthorized,
			ErrorMessage: "Invalid or expired token",
		}),
		SecretKey:       []byte(secretKey),
		TokenExpiration: tokenExpiration,
	}
}

// CanActivate kiểm tra JWT token trong request
func (g *JWTGuard) CanActivate(ctx *GuardContext) (bool, error) {
	// Kiểm tra xem có nên bỏ qua path này không
	if g.ShouldSkip(ctx.Path) {
		return true, nil
	}

	// Lấy token từ header
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return false, ErrMissingToken
	}

	// Kiểm tra format của token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false, ErrInvalidToken
	}

	// Parse và validate token
	token, err := jwt.ParseWithClaims(parts[1], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return g.SecretKey, nil
	})

	if err != nil {
		if err.Error() == "token is expired" {
			return false, ErrTokenExpired
		}
		return false, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Lưu thông tin user vào context
		ctx.User = claims
		return true, nil
	}

	return false, ErrInvalidToken
}

// GenerateToken tạo một JWT token mới
func (g *JWTGuard) GenerateToken(userID, username, role string) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.SecretKey)
}

// ValidateToken kiểm tra tính hợp lệ của token
func (g *JWTGuard) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return g.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// OnRegister được gọi khi guard được đăng ký
func (g *JWTGuard) OnRegister(ctx context.Context) error {
	// Có thể thêm logic khởi tạo ở đây
	return nil
}

// OnShutdown được gọi khi ứng dụng đang shutdown
func (g *JWTGuard) OnShutdown(ctx context.Context) error {
	// Có thể thêm logic cleanup ở đây
	return nil
}
