# WiFi QR Code Generator - API Specification

## Overview

**Version**: 1.0
**Base URL**: `https://api.wifiqr.com/api` (Production)
**Base URL**: `http://localhost:8080/api` (Development)
**Protocol**: HTTPS (Production), HTTP (Development)
**Authentication**: JWT Bearer Token
**Content-Type**: `application/json`

---

## Table of Contents

1. [Authentication](#authentication)
2. [Error Handling](#error-handling)
3. [Pagination](#pagination)
4. [Endpoints](#endpoints)
   - [Auth Endpoints](#auth-endpoints)
   - [QR Code Endpoints](#qr-code-endpoints)
   - [Admin Endpoints](#admin-endpoints)
5. [Data Models](#data-models)
6. [Example Requests](#example-requests)

---

## Authentication

All protected endpoints require a valid JWT token in the Authorization header.

### Header Format
```http
Authorization: Bearer <jwt_token>
```

### Token Acquisition
Tokens are obtained through the `/auth/login` or `/auth/register` endpoints.

### Token Expiration
- **Duration**: 1 hour
- **Refresh**: Not implemented (user must re-authenticate)
- **Future Enhancement**: Implement refresh token mechanism

### Token Claims
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440001",
  "email": "user@example.com",
  "role": "user",
  "exp": 1705324200,
  "iat": 1705320600
}
```

---

## Error Handling

All error responses follow a consistent format.

### Error Response Structure
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "specific error detail"
    }
  }
}
```

### HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, DELETE |
| 201 | Created | Successful POST |
| 400 | Bad Request | Validation errors, malformed request |
| 401 | Unauthorized | Missing or invalid JWT token |
| 403 | Forbidden | Valid token but insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate resource (e.g., email exists) |
| 422 | Unprocessable Entity | Business logic validation failed |
| 500 | Internal Server Error | Unexpected server error |

### Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Input validation failed |
| `INVALID_CREDENTIALS` | Email or password incorrect |
| `UNAUTHORIZED` | Missing or invalid token |
| `FORBIDDEN` | Insufficient permissions |
| `RESOURCE_NOT_FOUND` | Requested resource doesn't exist |
| `DUPLICATE_EMAIL` | Email already registered |
| `INTERNAL_ERROR` | Unexpected server error |

### Example Error Responses

**Validation Error (400)**:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "ssid": "SSID is required",
      "password": "Password is required for secured networks"
    }
  }
}
```

**Unauthorized (401)**:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid authentication token"
  }
}
```

**Forbidden (403)**:
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}
```

---

## Pagination

List endpoints support pagination using query parameters.

### Parameters
- `page` (integer): Page number, starting from 1. Default: 1
- `pageSize` (integer): Items per page. Default: 10, Max: 100

### Response Format
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "total": 150,
      "page": 1,
      "pageSize": 10,
      "totalPages": 15
    }
  }
}
```

### Example Request
```http
GET /api/qr-codes?page=2&pageSize=20
```

---

## Endpoints

## Auth Endpoints

### Register User

Create a new user account.

**Endpoint**: `POST /api/auth/register`
**Access**: Public
**Rate Limit**: 5 requests per minute per IP

#### Request Body
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Validation Rules
| Field | Type | Required | Rules |
|-------|------|----------|-------|
| email | string | Yes | Valid email format, unique, max 255 chars |
| password | string | Yes | Min 8 chars, 1 uppercase, 1 lowercase, 1 number |

#### Success Response (201 Created)
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "email": "user@example.com",
      "role": "user",
      "createdAt": "2025-01-15T10:30:00Z"
    }
  },
  "message": "User registered successfully"
}
```

#### Error Responses

**Email already exists (409)**:
```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE_EMAIL",
    "message": "Email address already registered"
  }
}
```

**Validation error (400)**:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "email": "Invalid email format",
      "password": "Password must be at least 8 characters"
    }
  }
}
```

---

### Login User

Authenticate user and receive JWT token.

**Endpoint**: `POST /api/auth/login`
**Access**: Public
**Rate Limit**: 5 requests per minute per IP

#### Request Body
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Validation Rules
| Field | Type | Required | Rules |
|-------|------|----------|-------|
| email | string | Yes | Valid email format |
| password | string | Yes | Non-empty |

#### Success Response (200 OK)
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "email": "user@example.com",
      "role": "user",
      "createdAt": "2025-01-15T10:30:00Z"
    }
  },
  "message": "Login successful"
}
```

#### Error Responses

**Invalid credentials (401)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

---

## QR Code Endpoints

### Create QR Code

Create a new WiFi QR code.

**Endpoint**: `POST /api/qr-codes`
**Access**: Protected (requires valid JWT)
**Rate Limit**: 20 requests per minute per user

#### Request Headers
```http
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body
```json
{
  "ssid": "MyHomeWiFi",
  "password": "SecurePassword123",
  "securityType": "WPA2",
  "isHidden": false
}
```

#### Validation Rules
| Field | Type | Required | Rules |
|-------|------|----------|-------|
| ssid | string | Yes | Max 32 chars (IEEE 802.11 spec) |
| password | string | Conditional | Required for WPA/WPA2/WEP, max 63 chars |
| securityType | string | Yes | One of: "WPA", "WPA2", "WEP", "nopass" |
| isHidden | boolean | No | Default: false |

#### Success Response (201 Created)
```json
{
  "success": true,
  "data": {
    "id": "650e8400-e29b-41d4-a716-446655440010",
    "userId": "550e8400-e29b-41d4-a716-446655440001",
    "ssid": "MyHomeWiFi",
    "password": "SecurePassword123",
    "securityType": "WPA2",
    "isHidden": false,
    "qrCodeData": "WIFI:T:WPA2;S:MyHomeWiFi;P:SecurePassword123;H:false;;",
    "qrCodeImageUrl": null,
    "createdAt": "2025-01-15T11:00:00Z",
    "updatedAt": "2025-01-15T11:00:00Z"
  },
  "message": "QR code created successfully"
}
```

#### WiFi QR Code Format Specification
```
WIFI:T:<security_type>;S:<ssid>;P:<password>;H:<hidden>;;
```

**Components**:
- `T`: Security type (WPA, WPA2, WEP, or nopass)
- `S`: Network SSID (escaped special characters)
- `P`: Password (escaped special characters)
- `H`: Hidden network (true or false)

**Special Character Escaping**:
- Backslash (`\`), semicolon (`;`), colon (`:`), comma (`,`) must be escaped with backslash
- Example: `My;Network` becomes `My\;Network`

#### Error Responses

**Validation error (400)**:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "password": "Password is required for WPA2 networks"
    }
  }
}
```

**Unauthorized (401)**:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid authentication token"
  }
}
```

---

### Get User's QR Codes

Retrieve all QR codes created by the authenticated user.

**Endpoint**: `GET /api/qr-codes`
**Access**: Protected (requires valid JWT)

#### Request Headers
```http
Authorization: Bearer <token>
```

#### Query Parameters
- `page` (integer, optional): Page number, default 1
- `pageSize` (integer, optional): Items per page, default 10, max 100

#### Example Request
```http
GET /api/qr-codes?page=1&pageSize=10
```

#### Success Response (200 OK)
```json
{
  "success": true,
  "data": {
    "qrCodes": [
      {
        "id": "650e8400-e29b-41d4-a716-446655440010",
        "userId": "550e8400-e29b-41d4-a716-446655440001",
        "ssid": "MyHomeWiFi",
        "password": "SecurePassword123",
        "securityType": "WPA2",
        "isHidden": false,
        "qrCodeData": "WIFI:T:WPA2;S:MyHomeWiFi;P:SecurePassword123;H:false;;",
        "qrCodeImageUrl": null,
        "createdAt": "2025-01-15T11:00:00Z",
        "updatedAt": "2025-01-15T11:00:00Z"
      },
      {
        "id": "650e8400-e29b-41d4-a716-446655440011",
        "userId": "550e8400-e29b-41d4-a716-446655440001",
        "ssid": "GuestNetwork",
        "password": null,
        "securityType": "nopass",
        "isHidden": false,
        "qrCodeData": "WIFI:T:nopass;S:GuestNetwork;P:;H:false;;",
        "qrCodeImageUrl": null,
        "createdAt": "2025-01-15T10:00:00Z",
        "updatedAt": "2025-01-15T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 25,
      "page": 1,
      "pageSize": 10,
      "totalPages": 3
    }
  },
  "message": "QR codes retrieved successfully"
}
```

---

### Get QR Code by ID

Retrieve a specific QR code by its ID.

**Endpoint**: `GET /api/qr-codes/:id`
**Access**: Protected (user must own the QR code)

#### Request Headers
```http
Authorization: Bearer <token>
```

#### URL Parameters
- `id` (UUID): QR code identifier

#### Example Request
```http
GET /api/qr-codes/650e8400-e29b-41d4-a716-446655440010
```

#### Success Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": "650e8400-e29b-41d4-a716-446655440010",
    "userId": "550e8400-e29b-41d4-a716-446655440001",
    "ssid": "MyHomeWiFi",
    "password": "SecurePassword123",
    "securityType": "WPA2",
    "isHidden": false,
    "qrCodeData": "WIFI:T:WPA2;S:MyHomeWiFi;P:SecurePassword123;H:false;;",
    "qrCodeImageUrl": null,
    "createdAt": "2025-01-15T11:00:00Z",
    "updatedAt": "2025-01-15T11:00:00Z"
  },
  "message": "QR code retrieved successfully"
}
```

#### Error Responses

**Not found (404)**:
```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "QR code not found"
  }
}
```

**Forbidden (403)**:
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this QR code"
  }
}
```

---

### Delete QR Code

Delete a specific QR code (soft delete).

**Endpoint**: `DELETE /api/qr-codes/:id`
**Access**: Protected (user must own the QR code)

#### Request Headers
```http
Authorization: Bearer <token>
```

#### URL Parameters
- `id` (UUID): QR code identifier

#### Example Request
```http
DELETE /api/qr-codes/650e8400-e29b-41d4-a716-446655440010
```

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "QR code deleted successfully"
}
```

#### Error Responses

**Not found (404)**:
```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "QR code not found"
  }
}
```

**Forbidden (403)**:
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to delete this QR code"
  }
}
```

---

## Admin Endpoints

### Get All Credentials

Retrieve all WiFi credentials across all users (admin only).

**Endpoint**: `GET /api/admin/credentials`
**Access**: Protected (requires admin role)
**Rate Limit**: 100 requests per minute

#### Request Headers
```http
Authorization: Bearer <token>
```

#### Query Parameters
- `page` (integer, optional): Page number, default 1
- `pageSize` (integer, optional): Items per page, default 20, max 100
- `userId` (UUID, optional): Filter by specific user
- `search` (string, optional): Search by SSID (partial match)

#### Example Requests
```http
GET /api/admin/credentials?page=1&pageSize=20
GET /api/admin/credentials?userId=550e8400-e29b-41d4-a716-446655440001
GET /api/admin/credentials?search=office
GET /api/admin/credentials?page=2&pageSize=50&search=wifi
```

#### Success Response (200 OK)
```json
{
  "success": true,
  "data": {
    "credentials": [
      {
        "id": "650e8400-e29b-41d4-a716-446655440010",
        "userId": "550e8400-e29b-41d4-a716-446655440001",
        "userEmail": "user@example.com",
        "userRole": "user",
        "ssid": "MyHomeWiFi",
        "password": "SecurePassword123",
        "securityType": "WPA2",
        "isHidden": false,
        "qrCodeData": "WIFI:T:WPA2;S:MyHomeWiFi;P:SecurePassword123;H:false;;",
        "createdAt": "2025-01-15T11:00:00Z",
        "updatedAt": "2025-01-15T11:00:00Z"
      },
      {
        "id": "650e8400-e29b-41d4-a716-446655440015",
        "userId": "550e8400-e29b-41d4-a716-446655440002",
        "userEmail": "another@example.com",
        "userRole": "user",
        "ssid": "Office-WiFi",
        "password": "OfficeSecure456",
        "securityType": "WPA2",
        "isHidden": true,
        "qrCodeData": "WIFI:T:WPA2;S:Office-WiFi;P:OfficeSecure456;H:true;;",
        "createdAt": "2025-01-15T09:30:00Z",
        "updatedAt": "2025-01-15T09:30:00Z"
      }
    ],
    "pagination": {
      "total": 150,
      "page": 1,
      "pageSize": 20,
      "totalPages": 8
    }
  },
  "message": "Credentials retrieved successfully"
}
```

#### Error Responses

**Forbidden (403)**:
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Admin role required to access this endpoint"
  }
}
```

---

## Data Models

### User Model

```typescript
interface User {
  id: string;              // UUID
  email: string;           // Unique email address
  role: 'user' | 'admin';  // User role
  createdAt: string;       // ISO 8601 timestamp
  updatedAt: string;       // ISO 8601 timestamp
}
```

**Note**: `password_hash` is never exposed in API responses.

---

### WiFi QR Code Model

```typescript
interface WiFiQRCode {
  id: string;                      // UUID
  userId: string;                  // UUID of owner
  ssid: string;                    // Network name (max 32 chars)
  password: string | null;         // Plain password (decrypted), null for open networks
  securityType: 'WPA' | 'WPA2' | 'WEP' | 'nopass';
  isHidden: boolean;               // Hidden network flag
  qrCodeData: string;              // WiFi QR code string format
  qrCodeImageUrl: string | null;   // Optional image URL
  createdAt: string;               // ISO 8601 timestamp
  updatedAt: string;               // ISO 8601 timestamp
}
```

---

### Admin Credential Model

```typescript
interface AdminCredential extends WiFiQRCode {
  userEmail: string;    // Email of QR code owner
  userRole: string;     // Role of QR code owner
}
```

---

## Example Requests

### Using cURL

#### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@12345"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@12345"
  }'
```

#### Create QR Code
```bash
curl -X POST http://localhost:8080/api/qr-codes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "ssid": "MyWiFi",
    "password": "WiFiPass123",
    "securityType": "WPA2",
    "isHidden": false
  }'
```

#### Get My QR Codes
```bash
curl -X GET "http://localhost:8080/api/qr-codes?page=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Delete QR Code
```bash
curl -X DELETE http://localhost:8080/api/qr-codes/650e8400-e29b-41d4-a716-446655440010 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Admin: Get All Credentials
```bash
curl -X GET "http://localhost:8080/api/admin/credentials?page=1&pageSize=20" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### Using JavaScript Fetch

#### Login
```javascript
const response = await fetch('http://localhost:8080/api/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'test@example.com',
    password: 'Test@12345'
  })
});

const data = await response.json();
const token = data.data.token;
```

#### Create QR Code
```javascript
const response = await fetch('http://localhost:8080/api/qr-codes', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    ssid: 'MyWiFi',
    password: 'WiFiPass123',
    securityType: 'WPA2',
    isHidden: false
  })
});

const qrCode = await response.json();
```

---

### Using Angular HttpClient

#### Auth Service
```typescript
login(credentials: { email: string; password: string }): Observable<AuthResponse> {
  return this.http.post<AuthResponse>(
    `${this.baseUrl}/auth/login`,
    credentials
  );
}
```

#### QR Code Service
```typescript
createQRCode(request: CreateQRCodeRequest): Observable<WiFiQRCode> {
  return this.http.post<WiFiQRCode>(
    `${this.baseUrl}/qr-codes`,
    request
  );
}

getMyQRCodes(page = 1, pageSize = 10): Observable<PaginatedResponse<WiFiQRCode>> {
  const params = new HttpParams()
    .set('page', page.toString())
    .set('pageSize', pageSize.toString());

  return this.http.get<PaginatedResponse<WiFiQRCode>>(
    `${this.baseUrl}/qr-codes`,
    { params }
  );
}
```

---

## Rate Limiting

### Current Limits (Recommended)

| Endpoint | Limit | Window |
|----------|-------|--------|
| `/auth/register` | 5 requests | 1 minute |
| `/auth/login` | 5 requests | 1 minute |
| `/qr-codes` (POST) | 20 requests | 1 minute |
| `/qr-codes` (GET) | 100 requests | 1 minute |
| `/admin/*` | 100 requests | 1 minute |

### Rate Limit Headers
```http
X-RateLimit-Limit: 20
X-RateLimit-Remaining: 15
X-RateLimit-Reset: 1705324200
```

### Rate Limit Error (429)
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please try again later.",
    "retryAfter": 45
  }
}
```

---

## Versioning

**Current Version**: v1 (implicit in `/api` prefix)

**Future Versioning Strategy**:
- Breaking changes: New version path (`/api/v2`)
- Non-breaking changes: Same version with documentation updates
- Deprecation: 6-month notice period

---

## Security Best Practices

1. **Always use HTTPS** in production
2. **Never log sensitive data** (passwords, tokens)
3. **Validate all input** on both client and server
4. **Sanitize error messages** (don't expose internal details)
5. **Implement rate limiting** to prevent abuse
6. **Rotate JWT secrets** periodically
7. **Use secure password requirements**
8. **Encrypt sensitive data** at rest (WiFi passwords)

---

**Document Version**: 1.0
**Last Updated**: 2025-01-15
**Maintainer**: Backend Team
