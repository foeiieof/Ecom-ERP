# E-Commerce Hexagonal Architecture in Go

A modern e-commerce API built with Go using hexagonal architecture principles, featuring OAuth authentication, JWT tokens, and containerized deployment.

## Features

- **Hexagonal Architecture**: Clean separation of concerns with domain, application, and infrastructure layers
- **OAuth Authentication**: Google OAuth integration with JWT tokens
- **RESTful API**: Well-structured REST endpoints with proper HTTP status codes
- **Middleware**: Authentication, CORS, logging, and recovery middleware
- **Container Ready**: Docker and Docker Compose support
- **Hot Reload**: Development setup with Air for hot reloading
- **Health Checks**: Built-in health check endpoints
- **Graceful Shutdown**: Proper server shutdown handling

## Architecture

```
├── application/          # Application services (use cases)
│   └── auth/            # Authentication service
├── cmd/                 # Application entry points
│   └── server/          # Main server application
├── domain/              # Domain entities and interfaces
│   └── auth/            # Authentication domain
├── internal/            # Internal packages
│   ├── adapter/         # External adapters
│   │   ├── jwt/         # JWT service
│   │   ├── oauth/       # OAuth providers
│   │   └── repository/  # Data repositories
│   ├── config/          # Configuration management
│   ├── delivery/        # Delivery mechanisms
│   │   └── http/        # HTTP handlers and middleware
│   └── infrastructure/  # Infrastructure setup
└── docs/                # Documentation
```

## Prerequisites

- Go 1.23 or later
- Docker (optional, for containerized deployment)
- Google OAuth credentials (for authentication)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd ecommerce
```

### 2. Install Dependencies

```bash
make deps
# or
go mod download
```

### 3. Configure Environment

Copy and update the environment variables:

```bash
cp .env.example .env
```

Update the following required variables in `.env`:

```env
# JWT Configuration (REQUIRED)
OAUTH_JWT_SECRET_KEY=your-super-secret-jwt-key-minimum-32-characters

# Google OAuth (REQUIRED for authentication)
OAUTH_GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
OAUTH_GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
OAUTH_GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback
```

### 4. Run the Application

#### Development Mode (with hot reload)
```bash
make dev
# or
air
```

#### Production Mode
```bash
make run
# or
go run ./cmd/server/main.go
```

#### Using Docker
```bash
make docker-compose-up
# or
docker-compose up --build
```

## API Endpoints

### Public Endpoints

- `GET /` - Welcome message (requires authentication)
- `GET /health` - Health check
- `GET /products` - Public products (optional authentication for personalization)

### Authentication Endpoints

- `POST /auth/login` - Initiate OAuth login
- `POST /auth/callback` - OAuth callback handler
- `POST /auth/refresh` - Refresh access token

### Protected Endpoints (Require Authentication)

- `GET /api/profile` - Get user profile
- `POST /api/logout` - Logout user
- `GET /api/dashboard` - User dashboard

## Authentication Flow

### 1. Initiate Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"provider": "google"}'
```

Response:
```json
{
  "auth_url": "https://accounts.google.com/oauth/authorize?...",
  "state": "random-state-string"
}
```

### 2. Handle Callback

After user authorizes, your application receives a callback with code and state:

```bash
curl -X POST http://localhost:8080/auth/callback \
  -H "Content-Type: application/json" \
  -d '{
    "code": "authorization-code-from-google",
    "state": "state-from-step-1",
    "provider": "google"
  }'
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

### 3. Use Protected Endpoints

```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer your-jwt-token"
```

## Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/callback`
6. Copy Client ID and Client Secret to your `.env` file

## Development

### Available Make Commands

```bash
make help                 # Show all available commands
make build               # Build the application
make run                 # Run the application
make dev                 # Run with hot reload
make test                # Run tests
make clean               # Clean build artifacts
make docker-build        # Build Docker image
make docker-run          # Run Docker container
make docker-compose-up   # Start with docker-compose
make fmt                 # Format code
make lint                # Run linter
```

### Project Structure

The project follows hexagonal architecture principles:

- **Domain Layer**: Contains business entities and interfaces
- **Application Layer**: Contains use cases and business logic
- **Infrastructure Layer**: Contains external dependencies and adapters

### Adding New Features

1. Define domain entities in `domain/`
2. Create application services in `application/`
3. Implement adapters in `internal/adapter/`
4. Add HTTP handlers in `internal/delivery/http/handler/`
5. Wire dependencies in `internal/infrastructure/container.go`

## Deployment

### Docker

```bash
# Build image
docker build -t ecommerce-api .

# Run container
docker run -p 8080:8080 --env-file .env ecommerce-api
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Considerations

1. **Environment Variables**: Use proper secret management
2. **Database**: Replace in-memory storage with persistent database
3. **Redis**: Add Redis for token storage and caching
4. **Load Balancer**: Use nginx or similar for load balancing
5. **Monitoring**: Add metrics and logging
6. **Security**: Implement rate limiting and security headers

## Configuration

All configuration is done through environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | Server port | `8080` | No |
| `SERVER_HOST` | Server host | `0.0.0.0` | No |
| `OAUTH_JWT_SECRET_KEY` | JWT signing key | - | Yes |
| `OAUTH_JWT_EXPIRATION_TIME` | Token expiration (seconds) | `3600` | No |
| `OAUTH_GOOGLE_OAUTH_CLIENT_ID` | Google OAuth client ID | - | Yes |
| `OAUTH_GOOGLE_OAUTH_CLIENT_SECRET` | Google OAuth client secret | - | Yes |
| `OAUTH_GOOGLE_OAUTH_REDIRECT_URL` | OAuth redirect URL | - | Yes |

## Testing

```bash
# Run all tests
make test

# Run with coverage
go test -v -cover ./...

# Test specific package
go test -v ./application/auth
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make fmt` and `make lint`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please open an issue in the repository.
