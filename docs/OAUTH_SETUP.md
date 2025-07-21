# OAuth 2.0 Setup Guide

This guide explains how to set up OAuth 2.0 authentication with Google in your e-commerce application.

## Architecture Overview

The OAuth implementation follows hexagonal architecture principles with dependency injection:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  Application     │    │    Domain       │
│  (Handlers &    │───▶│    Services      │───▶│   Interfaces    │
│   Middleware)   │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Infrastructure │    │    Adapters      │    │  Repositories   │
│   (Container)   │    │ (OAuth, JWT)     │    │   (Storage)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Components

### 1. Domain Layer (`domain/auth/`)
- **Interfaces**: `OAuthProvider`, `AuthService`, `TokenRepository`
- **Models**: `User`, `TokenInfo`

### 2. Application Layer (`application/auth/`)
- **Service**: Implements business logic for authentication

### 3. Infrastructure Layer (`internal/`)
- **Adapters**: OAuth providers (Google), JWT service
- **Middleware**: Authentication middleware with DI
- **Handlers**: HTTP handlers for auth endpoints
- **Container**: Dependency injection container

## Setup Instructions

### 1. Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client IDs"
5. Set authorized redirect URIs: `http://localhost:8080/auth/callback`
6. Copy Client ID and Client Secret

### 2. Environment Configuration

Update your `.env` file:

```env
# OAuth Configuration
OAUTH_GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
OAUTH_GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
OAUTH_GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback

# JWT Configuration
OAUTH_JWT_SECRET_KEY=your-super-secret-jwt-key-change-this-in-production
OAUTH_JWT_EXPIRATION_TIME=3600
OAUTH_JWT_ISSUER=ecommerce-api
```

### 3. Run the Application

```bash
go run cmd/server/main.go
```

## API Endpoints

### Authentication Flow

1. **Initiate Login**
   ```bash
   POST /auth/login
   Content-Type: application/json
   
   {
     "provider": "google"
   }
   ```
   
   Response:
   ```json
   {
     "auth_url": "https://accounts.google.com/oauth/authorize?...",
     "state": "random-state-string"
   }
   ```

2. **Handle Callback**
   ```bash
   POST /auth/callback
   Content-Type: application/json
   
   {
     "code": "authorization-code-from-google",
     "state": "state-from-step-1",
     "provider": "google"
   }
   ```
   
   Response:
   ```json
   {
     "user": {
       "id": "user-id",
       "email": "user@example.com",
       "name": "User Name",
       "provider": "google"
     },
     "token": {
       "access_token": "jwt-token",
       "token_type": "Bearer",
       "expires_at": "2024-01-01T00:00:00Z"
     }
   }
   ```

### Protected Endpoints

Use the JWT token in Authorization header:

```bash
GET /api/profile
Authorization: Bearer your-jwt-token
```

### Public Endpoints (Optional Auth)

```bash
GET /products
# Works without auth, but shows personalized content if authenticated
Authorization: Bearer your-jwt-token  # Optional
```

## Middleware Configuration

The OAuth middleware supports three modes:

1. **Skip Paths**: No authentication required
   ```go
   SkipPaths: []string{"/health", "/auth/login", "/auth/callback"}
   ```

2. **Optional Paths**: Authentication optional, enhances response if present
   ```go
   OptionalPaths: []string{"/products"}
   ```

3. **Required Auth**: Strict authentication required
   ```go
   protected.Use(container.OAuthMiddleware.RequireAuth())
   ```

## Dependency Injection

The container (`internal/infrastructure/container.go`) wires all dependencies:

```go
// Initialize container
container := infrastructure.NewContainer(cfg)

// All dependencies are automatically wired:
// - OAuth providers
// - JWT service
// - Auth service
// - Token repository
// - Handlers
// - Middleware
```

## Security Features

- **JWT Token Validation**: Secure token validation with expiration
- **Token Blacklisting**: Logout functionality blacklists tokens
- **CSRF Protection**: State parameter validation (implement in production)
- **Secure Headers**: CORS and security headers configured
- **Token Storage**: In-memory storage (replace with Redis/DB for production)

## Production Considerations

1. **Replace In-Memory Storage**: Use Redis or database for token storage
2. **Implement State Validation**: Add proper CSRF protection
3. **Use HTTPS**: Always use HTTPS in production
4. **Rotate JWT Secrets**: Implement secret rotation
5. **Add Rate Limiting**: Protect against brute force attacks
6. **Implement Refresh Tokens**: Add proper refresh token flow
7. **Add Logging**: Comprehensive audit logging

## Testing

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test login initiation
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"provider": "google"}'

# Test protected endpoint
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer your-jwt-token"
```

## Adding New OAuth Providers

1. Implement `auth.OAuthProvider` interface
2. Add provider configuration to `config/oauth.go`
3. Register provider in container
4. Update environment variables

Example for GitHub:
```go
// Add to container.go
githubProvider := oauth.NewGitHubProvider(&c.Config.OAuth.GitHub)
c.OAuthProviders["github"] = githubProvider
```
