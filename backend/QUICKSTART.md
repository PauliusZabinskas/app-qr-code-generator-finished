# Quick Start Guide

## Start the Backend

From the project root directory:

```bash
docker-compose up -d postgres backend
```

The backend will be available at `http://localhost:8080`

## Test the API

### 1. Health Check
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

### 2. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Response will include a JWT token and user info.

### 3. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Copy the `token` from the response for authenticated requests.

### 4. Create WiFi Credential
Replace `YOUR_TOKEN_HERE` with the token from step 2 or 3:

```bash
curl -X POST http://localhost:8080/api/wifi \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "ssid": "MyHomeWiFi",
    "password": "mypassword123",
    "security_type": "WPA2",
    "is_hidden": false
  }'
```

Response includes the WiFi credential with a base64-encoded QR code in `qr_code_data`.

### 5. Get All WiFi Credentials
```bash
curl -X GET http://localhost:8080/api/wifi \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 6. Get Specific WiFi Credential
Replace `CREDENTIAL_ID` with an actual ID from step 4 or 5:

```bash
curl -X GET http://localhost:8080/api/wifi/CREDENTIAL_ID \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 7. Delete WiFi Credential
```bash
curl -X DELETE http://localhost:8080/api/wifi/CREDENTIAL_ID \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Create an Admin User

To test admin endpoints, you need to manually set a user's role to 'admin' in the database:

1. Connect to PostgreSQL:
```bash
docker exec -it wifiqr-postgres psql -U wifiqr -d wifiqr_db
```

2. Update user role:
```sql
UPDATE users SET role = 'admin' WHERE email = 'test@example.com';
\q
```

3. Login again to get a new token with admin role.

### Admin Endpoints

#### Get All Users
```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE"
```

#### Get All Credentials
```bash
curl -X GET http://localhost:8080/api/admin/credentials \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE"
```

#### Get System Stats
```bash
curl -X GET http://localhost:8080/api/admin/stats \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE"
```

## View Logs

```bash
docker-compose logs -f backend
```

## Stop the Backend

```bash
docker-compose down
```

## Troubleshooting

### Backend won't start
1. Check if PostgreSQL is running:
   ```bash
   docker-compose ps
   ```

2. View backend logs:
   ```bash
   docker-compose logs backend
   ```

3. Restart the backend:
   ```bash
   docker-compose restart backend
   ```

### Database connection errors
1. Ensure PostgreSQL is healthy:
   ```bash
   docker-compose ps postgres
   ```

2. Check environment variables in docker-compose.yml

### JWT Secret errors
- Ensure JWT_SECRET is at least 32 characters
- Ensure ENCRYPTION_KEY is exactly 32 characters
- Check docker-compose.yml environment section

## Development Mode

For hot reload during development:

```bash
docker-compose up backend
```

Any changes to `.go` files will automatically rebuild and restart the server.

## Next Steps

1. Integrate with the Angular frontend
2. Test all endpoints with Postman or similar tool
3. Implement additional features as needed
4. Deploy to production environment
