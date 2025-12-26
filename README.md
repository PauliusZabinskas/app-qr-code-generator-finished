# WiFi QR Code Generator

A modern full-stack web application that enables users to create scannable QR codes containing WiFi credentials. Built with Angular 21+, Go Gin framework, and PostgreSQL.

## Features

- User registration and authentication with JWT
- Role-based access control (User and Admin roles)
- WiFi QR code generation and management
- Secure credential storage with AES-256-GCM encryption
- Admin dashboard for system-wide credential viewing
- Responsive design with Tailwind CSS
- RESTful API with comprehensive validation

## Technology Stack

### Frontend
- **Framework**: Angular 21+ (Standalone Components)
- **Language**: TypeScript 5.3+
- **Styling**: Tailwind CSS 3.4+
- **State Management**: Angular Signals
- **QR Code Library**: angularx-qrcode
- **Build Tool**: esbuild

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **ORM**: GORM with PostgreSQL driver
- **Authentication**: JWT (golang-jwt/jwt)
- **Encryption**: AES-256-GCM
- **Validation**: go-playground/validator
- **Logging**: zap (structured logging)

### Database
- **RDBMS**: PostgreSQL 16+
- **Containerization**: Docker with docker-compose
- **Migrations**: GORM AutoMigrate (dev) / golang-migrate (prod)

## Project Structure

```
app-wifi-qr-code/
├── docs/                           # Documentation
│   ├── ARCHITECTURE.md             # System architecture and design
│   ├── API_SPECIFICATION.md        # API endpoint documentation
│   ├── FRONTEND_SPECIFICATION.md   # Frontend component specs
│   └── DATABASE_SCHEMA.sql         # Database schema and migrations
├── backend/                        # Go backend application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── repositories/
│   │   └── utils/
│   ├── pkg/
│   ├── go.mod
│   └── Dockerfile
├── frontend/                       # Angular frontend application
│   ├── src/
│   │   ├── app/
│   │   │   ├── core/
│   │   │   ├── features/
│   │   │   └── shared/
│   │   ├── assets/
│   │   └── environments/
│   ├── angular.json
│   ├── package.json
│   └── tailwind.config.js
├── docker-compose.yml              # Docker composition for all services
└── README.md
```

## Documentation

Comprehensive documentation is available in the `/docs` directory:

1. **[ARCHITECTURE.md](./docs/ARCHITECTURE.md)** - Complete system architecture including:
   - System overview and architecture diagrams
   - Technology stack details
   - Database schema design with ERD
   - Backend architecture (layers, middleware, services)
   - Frontend architecture (components, routing, state management)
   - API specifications
   - Security architecture
   - Deployment architecture
   - Data flow diagrams
   - Architectural Decision Records (ADRs)

2. **[API_SPECIFICATION.md](./docs/API_SPECIFICATION.md)** - Detailed API documentation:
   - Authentication endpoints
   - QR code CRUD endpoints
   - Admin endpoints
   - Request/response schemas
   - Error handling
   - Validation rules
   - Example requests (cURL, JavaScript, Angular)

3. **[FRONTEND_SPECIFICATION.md](./docs/FRONTEND_SPECIFICATION.md)** - Frontend implementation guide:
   - Project structure
   - Component specifications
   - Routing configuration
   - State management with Signals
   - Services and guards
   - Forms and validation
   - UI/UX guidelines
   - Performance optimization

4. **[DATABASE_SCHEMA.sql](./docs/DATABASE_SCHEMA.sql)** - Database implementation:
   - Complete schema definitions
   - Indexes and constraints
   - Triggers and functions
   - Views for analytics
   - Seed data for development
   - Security configuration
   - Maintenance queries

## Key Features

### Authentication & Authorization
- JWT-based authentication with 1-hour token expiration
- bcrypt password hashing (cost factor 10)
- Role-based access control with two roles:
  - **User**: Create and manage own QR codes
  - **Admin**: View all credentials system-wide

### WiFi QR Code Generation
- Standard WiFi QR code format: `WIFI:T:<type>;S:<ssid>;P:<password>;H:<hidden>;;`
- Support for multiple security types (WPA, WPA2, WEP, Open)
- Hidden network support
- Credential storage with AES-256-GCM encryption

### Security Features
- AES-256-GCM encryption for WiFi passwords
- CORS configuration with whitelisted origins
- Input validation on frontend and backend
- SQL injection prevention via prepared statements
- XSS protection via Angular's default sanitization
- Secure environment variable management

### User Experience
- Responsive design with Tailwind CSS
- Lazy-loaded routes for optimal performance
- Real-time form validation
- Error handling with user-friendly messages
- Loading states and spinners
- Clean, modern UI

## Getting Started

### Prerequisites
- Node.js 18+ and npm
- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 16+ (or use Docker)

### Quick Start with Docker Compose

1. Clone the repository:
```bash
git clone <repository-url>
cd app-wifi-qr-code
```

2. Create environment files:
```bash
# Backend .env
cp backend/.env.example backend/.env
# Edit backend/.env with your configuration
```

3. Start all services:
```bash
docker-compose up -d
```

4. Access the application:
   - Frontend: http://localhost:4200
   - Backend API: http://localhost:8080
   - PostgreSQL: localhost:5432

### Manual Setup

#### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run database migrations (automatic with GORM)
# Start the server
go run cmd/server/main.go
```

#### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm start

# Or build for production
npm run build
```

#### Database Setup

```bash
# Using Docker
docker run -d \
  --name wifiqr-postgres \
  -e POSTGRES_USER=wifiqr \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=wifiqr_db \
  -p 5432:5432 \
  postgres:16-alpine

# Load schema
psql -h localhost -U wifiqr -d wifiqr_db -f docs/DATABASE_SCHEMA.sql
```

## Default Credentials (Development)

### Admin Account
- **Email**: admin@wifiqr.com
- **Password**: Admin@123
- **Role**: admin

### User Account
- **Email**: user@example.com
- **Password**: User@123
- **Role**: user

**IMPORTANT**: Change these credentials in production!

## API Endpoints

### Authentication (Public)
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login

### QR Codes (Protected)
- `POST /api/qr-codes` - Create QR code
- `GET /api/qr-codes` - Get user's QR codes (paginated)
- `GET /api/qr-codes/:id` - Get specific QR code
- `DELETE /api/qr-codes/:id` - Delete QR code

### Admin (Admin Role Required)
- `GET /api/admin/credentials` - View all credentials

For detailed API documentation, see [API_SPECIFICATION.md](./docs/API_SPECIFICATION.md).

## Environment Variables

### Backend

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=wifiqr
DB_PASSWORD=your_password
DB_NAME=wifiqr_db

# Security
JWT_SECRET=your_256_bit_secret_key
ENCRYPTION_KEY=your_256_bit_encryption_key

# Server
PORT=8080
FRONTEND_URL=http://localhost:4200
GIN_MODE=release
```

### Frontend

```typescript
// environments/environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api'
};
```

## Development Workflow

### Backend Development
```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Frontend Development
```bash
# Development server with hot reload
npm start

# Run tests
npm test

# Lint
npm run lint
```

### Database Management
```bash
# Create migration (using golang-migrate)
migrate create -ext sql -dir migrations -seq migration_name

# Run migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/db" up

# Rollback
migrate -path migrations -database "postgres://user:pass@localhost:5432/db" down 1
```

## Testing

### Backend Tests
```bash
cd backend
go test ./... -v
```

### Frontend Tests
```bash
cd frontend
npm test
```

### End-to-End Tests
```bash
cd frontend
npm run e2e
```

## Deployment

### Production Build

#### Frontend
```bash
cd frontend
npm run build
# Output in dist/ directory
```

#### Backend
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -o wifiqr-server cmd/server/main.go
```

### Docker Deployment

```bash
# Build images
docker-compose -f docker-compose.prod.yml build

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

### Deployment Checklist
- [ ] Change default credentials
- [ ] Generate new JWT_SECRET and ENCRYPTION_KEY
- [ ] Configure CORS with production frontend URL
- [ ] Enable HTTPS with SSL certificate
- [ ] Set up database backups
- [ ] Configure monitoring and logging
- [ ] Set GIN_MODE=release
- [ ] Review and apply rate limiting
- [ ] Set up CI/CD pipeline
- [ ] Configure environment-specific settings

## Security Considerations

1. **Secrets Management**: Use environment variables or secret management services (AWS Secrets Manager, HashiCorp Vault)
2. **HTTPS**: Always use HTTPS in production
3. **Database**: Use connection pooling and prepared statements
4. **CORS**: Whitelist only trusted origins
5. **Rate Limiting**: Implement rate limiting on all endpoints
6. **Logging**: Never log sensitive data (passwords, tokens)
7. **Updates**: Keep dependencies up to date
8. **Backups**: Implement automated database backups

## Performance Optimization

### Frontend
- Lazy loading of routes
- OnPush change detection strategy
- Angular Signals for fine-grained reactivity
- Code splitting per route
- Asset optimization (images, fonts)

### Backend
- Database connection pooling
- Query optimization with proper indexes
- Caching strategies (Redis for future enhancement)
- Gzip compression for responses
- Structured logging for performance monitoring

### Database
- Proper indexing strategy
- Query optimization with EXPLAIN ANALYZE
- Regular VACUUM and ANALYZE
- Connection pooling configuration

## Troubleshooting

### Common Issues

**Issue**: Cannot connect to database
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection string
psql -h localhost -U wifiqr -d wifiqr_db
```

**Issue**: CORS errors
- Verify FRONTEND_URL in backend .env matches frontend origin
- Check CORS middleware configuration

**Issue**: JWT token invalid
- Check JWT_SECRET matches between requests
- Verify token hasn't expired (1 hour default)

**Issue**: Frontend build errors
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Contributing

1. Read the architecture documentation
2. Follow the project structure and conventions
3. Write tests for new features
4. Update documentation for API changes
5. Submit pull requests with clear descriptions

## License

[Your License Here]

## Support

For issues and questions:
- Check the [documentation](./docs/)
- Review the [API specification](./docs/API_SPECIFICATION.md)
- Create an issue in the repository

---

**Version**: 1.0.0
**Last Updated**: 2025-01-15
**Status**: Architecture Design Complete
