package auth

import (
	"context"
	"time"
)

// TokenInfo represents JWT token information
type TokenInfo struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// AuthService defines the authentication service interface
type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, *TokenInfo, error)
	Login(ctx context.Context, req *LoginRequest) (*User, *TokenInfo, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	RefreshToken(ctx context.Context, userID string) (*TokenInfo, error)
	Logout(ctx context.Context, userID string) error
}

// TokenRepository defines token storage interface
type TokenRepository interface {
	StoreToken(ctx context.Context, userID string, tokenInfo *TokenInfo) error
	GetToken(ctx context.Context, userID string) (*TokenInfo, error)
	DeleteToken(ctx context.Context, userID string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	BlacklistToken(ctx context.Context, token string, expiresAt time.Time) error
}
