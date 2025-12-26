# WiFi QR Code Generator - Backend API

A production-ready Go backend API built with Gin framework and PostgreSQL for generating and managing WiFi QR codes.

## Features

- **User Authentication**: JWT-based authentication with bcrypt password hashing
- **WiFi Credential Management**: Create, read, and delete WiFi credentials
- **QR Code Generation**: Automatic QR code generation for WiFi credentials
- **Security**: AES-256 encryption for WiFi passwords
- **Role-Based Access Control**: Admin and user roles with appropriate permissions
- **RESTful API**: Clean, consistent API design

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 16 with GORM ORM
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **QR Code**: skip2/go-qrcode
- **Encryption**: AES-256-GCM
- **Containerization**: Docker & Docker Compose
- **Hot Reload**: Air

## Project Structure

```
backend/
├── cmd/api/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── models/
│   │   ├── user.go             # User model
│   │   └── wifi_credential.go  # WiFi credential model
│   ├── repositories/
│   │   ├── user.go             # User data access layer
│   │   └── wifi.go             # WiFi credential data access layer
│   ├── services/
│   │   ├── auth.go             # Authentication business logic
│   │   ├── wifi.go             # WiFi credential business logic
│   │   └── qrcode.go           # QR code generation
│   ├── handlers/
│   │   ├── auth.go             # Authentication HTTP handlers
│   │   ├── wifi.go             # WiFi CRUD HTTP handlers
│   │   └── admin.go            # Admin HTTP handlers
│   ├── middleware/
│   │   ├── auth.go             # JWT authentication middleware
│   │   └── cors.go             # CORS middleware
│   └── routes/
│       └── routes.go           # Route definitions
├── Dockerfile                  # Multi-stage Docker build
├── .air.toml                   # Hot reload configuration
├── go.mod                      # Go module dependencies
└── README.md                   # This file
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and receive JWT token

### WiFi Credentials (Protected)
- `GET /api/wifi` - Get all WiFi credentials for current user
- `POST /api/wifi` - Create new WiFi credential with QR code
- `GET /api/wifi/:id` - Get specific WiFi credential
- `DELETE /api/wifi/:id` - Delete WiFi credential

### Admin (Protected, Admin Only)
- `GET /api/admin/users` - Get all users
- `GET /api/admin/credentials` - Get all WiFi credentials
- `GET /api/admin/stats` - Get system statistics

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development without Docker)

### Running with Docker Compose (Recommended)

1. Navigate to the project root:
   ```bash
   cd /path/to/app-wifi-qr-code
   ```

2. Start the services:
   ```bash
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`

4. Health check:
   ```bash
   curl http://localhost:8080/health
   ```

### Running Locally

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Set environment variables:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=wifiqr
   export DB_PASSWORD=dev_password_123
   export DB_NAME=wifiqr_db
   export DB_SSL_MODE=disable
   export JWT_SECRET=dev_jwt_secret_key_change_in_production_minimum_32_chars
   export ENCRYPTION_KEY=dev_encryption_key_change_in_production_exactly_32_chars
   export PORT=8080
   ```

3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

### Development with Hot Reload

Using Air for hot reload:

```bash
air -c .air.toml
```

## Configuration

Configuration is managed through environment variables:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| DB_HOST | PostgreSQL host | Yes | localhost |
| DB_PORT | PostgreSQL port | Yes | 5432 |
| DB_USER | Database user | Yes | wifiqr |
| DB_PASSWORD | Database password | Yes | - |
| DB_NAME | Database name | Yes | wifiqr_db |
| DB_SSL_MODE | SSL mode | No | disable |
| JWT_SECRET | JWT signing secret (min 32 chars) | Yes | - |
| ENCRYPTION_KEY | AES-256 key (exactly 32 chars) | Yes | - |
| PORT | Server port | No | 8080 |
| FRONTEND_URL | Frontend URL for CORS | No | http://localhost:4200 |
| GIN_MODE | Gin mode (debug/release) | No | debug |
| ALLOWED_ORIGINS | Comma-separated CORS origins | No | http://localhost:4200 |

## Security Features

1. **Password Hashing**: bcrypt with default cost factor
2. **JWT Tokens**: 24-hour expiration, HS256 signing
3. **WiFi Password Encryption**: AES-256-GCM encryption
4. **SQL Injection Protection**: Parameterized queries via GORM
5. **CORS Configuration**: Configurable allowed origins
6. **Role-Based Access**: Admin-only endpoints protected

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### WiFi Credentials Table
```sql
CREATE TABLE wifi_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ssid VARCHAR(255) NOT NULL,
    encrypted_password TEXT NOT NULL,
    security_type VARCHAR(20) NOT NULL,
    is_hidden BOOLEAN DEFAULT FALSE,
    qr_code_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Example Usage

### Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Create WiFi Credential
```bash
curl -X POST http://localhost:8080/api/wifi \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "ssid": "MyWiFi",
    "password": "wifipassword",
    "security_type": "WPA2",
    "is_hidden": false
  }'
```

## Production Deployment

For production deployment:

1. **Use strong secrets**:
   - Generate a secure JWT_SECRET (at least 32 characters)
   - Generate a random 32-character ENCRYPTION_KEY

2. **Enable SSL/TLS**:
   - Set DB_SSL_MODE to "require" or "verify-full"
   - Use HTTPS for the API

3. **Set Gin mode**:
   ```bash
   export GIN_MODE=release
   ```

4. **Configure CORS**:
   - Set ALLOWED_ORIGINS to your production frontend URL

5. **Use production Dockerfile target**:
   ```bash
   docker build --target production -t wifiqr-backend .
   ```

## Testing

Health check endpoint:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "database": "connected"
}
```

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running
- Verify database credentials
- Check network connectivity between containers

### JWT Token Issues
- Ensure JWT_SECRET is set and at least 32 characters
- Check token expiration (24 hours by default)

### QR Code Generation Issues
- Verify WiFi password meets length requirements (max 63 characters)
- Check security type is valid (WPA, WPA2, WEP, nopass)

## License

MIT License
