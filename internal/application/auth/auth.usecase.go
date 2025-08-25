package auth

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"ecommerce/domain/auth"
// 	// jwtService "ecommerce/internal/infrastructure/jwt"

// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// // Service implements the auth.AuthService interface
// type Service struct {
// 	userRepo   auth.UserRepository
// 	tokenRepo  auth.TokenRepository
// 	// jwtService *jwtService.Service
// }

// // NewService creates a new authentication service
// func NewService(
// 	userRepo auth.UserRepository,
// 	tokenRepo auth.TokenRepository,
// 	// jwtService *jwtService.Service,
// ) *Service {
// 	return &Service{
// 		userRepo:   userRepo,
// 		tokenRepo:  tokenRepo,
// 		// jwtService: jwtService,
// 	}
// }

// // Register creates a new user account
// func (s *Service) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.User, *auth.TokenInfo, error) {
// 	// Check if email already exists
// 	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
// 		return nil, nil, fmt.Errorf("email already registered")
// 	}

// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	// Create user
// 	now := time.Now()
// 	user := &auth.User{
// 		ID:        uuid.New().String(),
// 		Email:     req.Email,
// 		Password:  string(hashedPassword),
// 		Name:      req.Name,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}

// 	if err := s.userRepo.Create(ctx, user); err != nil {
// 		return nil, nil, fmt.Errorf("failed to create user: %w", err)
// 	}

// 	// Generate JWT token
// 	// token, expiresAt, err := s.jwtService.GenerateToken(user)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	tokenInfo := &auth.TokenInfo{
// 		// AccessToken: token,
// 		TokenType:   "Bearer",
// 		// ExpiresAt:   expiresAt,
// 	}

// 	// Store token
// 	if err := s.tokenRepo.StoreToken(ctx, user.ID, tokenInfo); err != nil {
// 		return nil, nil, fmt.Errorf("failed to store token: %w", err)
// 	}

// 	return user, tokenInfo, nil
// }

// // Login authenticates a user
// func (s *Service) Login(ctx context.Context, req *auth.LoginRequest) (*auth.User, *auth.TokenInfo, error) {
// 	// Get user by email
// 	user, err := s.userRepo.GetByEmail(ctx, req.Email)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("invalid email or password")
// 	}

// 	// Verify password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
// 		return nil, nil, fmt.Errorf("invalid email or password")
// 	}

// 	// Generate JWT token
// 	token, expiresAt, err := s.jwtService.GenerateToken(user)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	tokenInfo := &auth.TokenInfo{
// 		AccessToken: token,
// 		TokenType:   "Bearer",
// 		ExpiresAt:   expiresAt,
// 	}

// 	// Store token
// 	if err := s.tokenRepo.StoreToken(ctx, user.ID, tokenInfo); err != nil {
// 		return nil, nil, fmt.Errorf("failed to store token: %w", err)
// 	}

// 	return user, tokenInfo, nil
// }

// // ValidateToken validates a JWT token
// func (s *Service) ValidateToken(ctx context.Context, token string) (*auth.User, error) {
// 	// Check if token is blacklisted
// 	blacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, token)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
// 	}
// 	if blacklisted {
// 		return nil, fmt.Errorf("token has been revoked")
// 	}

// 	// Validate JWT
// 	claims, err := s.jwtService.ValidateToken(token)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid token: %w", err)
// 	}

// 	// Get user
// 	user, err := s.userRepo.GetByID(ctx, claims.UserID)
// 	if err != nil {
// 		return nil, fmt.Errorf("user not found: %w", err)
// 	}

// 	return user, nil
// }

// // RefreshToken generates a new token for a user
// func (s *Service) RefreshToken(ctx context.Context, userID string) (*auth.TokenInfo, error) {
// 	// Get user
// 	user, err := s.userRepo.GetByID(ctx, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("user not found: %w", err)
// 	}

// 	// Generate new token
// 	token, expiresAt, err := s.jwtService.GenerateToken(user)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	tokenInfo := &auth.TokenInfo{
// 		AccessToken: token,
// 		TokenType:   "Bearer",
// 		ExpiresAt:   expiresAt,
// 	}

// 	// Store new token
// 	if err := s.tokenRepo.StoreToken(ctx, userID, tokenInfo); err != nil {
// 		return nil, fmt.Errorf("failed to store token: %w", err)
// 	}

// 	return tokenInfo, nil
// }

// // Logout logs out a user
// func (s *Service) Logout(ctx context.Context, userID string) error {
// 	// Get current token
// 	tokenInfo, err := s.tokenRepo.GetToken(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get token: %w", err)
// 	}

// 	// Blacklist current token
// 	if err := s.tokenRepo.BlacklistToken(ctx, tokenInfo.AccessToken, tokenInfo.ExpiresAt); err != nil {
// 		return fmt.Errorf("failed to blacklist token: %w", err)
// 	}

// 	// Delete stored token
// 	if err := s.tokenRepo.DeleteToken(ctx, userID); err != nil {
// 		return fmt.Errorf("failed to delete token: %w", err)
// 	}

// 	return nil
// }
