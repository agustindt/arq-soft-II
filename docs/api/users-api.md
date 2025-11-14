# Users API Documentation

## Overview

The Users API is a Go-based microservice built with Gin framework that handles user authentication, authorization, and profile management. It implements JWT-based authentication with role-based access control (RBAC).

**Technology Stack:**
- Language: Go 1.21+
- Framework: Gin
- ORM: GORM
- Database: MySQL 8.0
- Authentication: JWT (golang-jwt/jwt/v5)
- Password Hashing: bcrypt

**Port:** 8081

**Base URL:** `http://localhost:8081/api/v1`

## Architecture

### Database Schema

The Users API uses MySQL as its primary database with the following schema:

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    
    -- Profile fields
    avatar_url VARCHAR(500),
    bio TEXT,
    phone VARCHAR(20),
    birth_date DATE,
    location VARCHAR(200),
    gender ENUM('male','female','other','prefer_not_to_say'),
    
    -- Sports-specific fields
    height DECIMAL(5,2) COMMENT 'Height in cm',
    weight DECIMAL(5,2) COMMENT 'Weight in kg',
    sports_interests JSON COMMENT 'Array of sports as JSON',
    fitness_level ENUM('beginner','intermediate','advanced','professional'),
    social_links JSON,
    
    -- System fields
    role VARCHAR(50) DEFAULT 'user',
    email_verified BOOLEAN DEFAULT false,
    email_verified_at DATETIME(3),
    is_active BOOLEAN DEFAULT true,
    last_login_at DATETIME(3),
    
    -- Timestamps
    created_at DATETIME(3),
    updated_at DATETIME(3),
    
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_role (role)
);
```

### JWT Structure

The API generates JWT tokens with the following claims:

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "role": "admin",
  "iss": "sports-activities-api",
  "sub": "user-authentication",
  "exp": 1763237675,
  "nbf": 1763150125,
  "iat": 1763150125
}
```

**Important:** The `role` claim is included in the JWT to enable efficient authorization checks without database queries.

### Role-Based Access Control

The API implements a hierarchical role system:

| Role | Level | Permissions |
|------|-------|-------------|
| **user** | 1 | Basic user operations, profile management |
| **moderator** | 2 | User role + content moderation |
| **admin** | 3 | Moderator role + user management, system configuration |
| **super_admin** | 4 | Admin role + advanced system settings |
| **root** | 5 | Full system access, user deletion |

**Hierarchy Rule:** Higher-level roles inherit all permissions from lower levels.

## API Endpoints

### Authentication Endpoints

#### POST /auth/register

Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2025-11-14T10:00:00Z"
    }
  }
}
```

**Error Responses:**
- `409 Conflict` - Email or username already exists
- `400 Bad Request` - Invalid request data

---

#### POST /auth/login

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "role": "admin",
      "is_active": true
    }
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials
- `401 Unauthorized` - Account disabled

---

#### POST /auth/refresh

Refresh an existing JWT token.

**Headers:**
```
Authorization: Bearer <existing_token>
```

**Response (200 OK):**
```json
{
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### Profile Endpoints (Protected)

All profile endpoints require JWT authentication via `Authorization: Bearer <token>` header.

#### GET /profile

Get the authenticated user's full profile.

**Response (200 OK):**
```json
{
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "avatar_url": "/uploads/avatars/avatar_1_1699900000.jpg",
    "bio": "Passionate about sports and fitness",
    "phone": "+5491123456789",
    "birth_date": "1995-06-15",
    "location": "Buenos Aires, Argentina",
    "gender": "male",
    "height": 180.5,
    "weight": 75.0,
    "sports_interests": "[\"football\", \"tennis\", \"running\"]",
    "fitness_level": "intermediate",
    "social_links": {
      "instagram": "@johndoe",
      "twitter": "@johndoe_sports"
    },
    "role": "user",
    "email_verified": false,
    "is_active": true,
    "created_at": "2025-11-14T10:00:00Z",
    "updated_at": "2025-11-14T12:30:00Z"
  }
}
```

---

#### PUT /profile

Update the authenticated user's profile.

**Request Body (all fields optional):**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "bio": "Updated bio text",
  "phone": "+5491123456789",
  "birth_date": "1995-06-15",
  "location": "Buenos Aires, Argentina",
  "gender": "male",
  "height": 180.5,
  "weight": 75.0,
  "sports_interests": "[\"football\", \"tennis\"]",
  "fitness_level": "advanced",
  "social_links": {
    "instagram": "@johndoe",
    "twitter": "@johndoe_sports"
  }
}
```

**Response (200 OK):**
```json
{
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "bio": "Updated bio text",
    ...
  }
}
```

---

#### PUT /profile/password

Change the authenticated user's password.

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response (200 OK):**
```json
{
  "message": "Password changed successfully"
}
```

**Error Responses:**
- `401 Unauthorized` - Current password is incorrect

---

#### POST /profile/avatar

Upload a new avatar image.

**Request:**
- Content-Type: `multipart/form-data`
- Form field: `avatar` (file)

**Accepted formats:** JPEG, JPG, PNG, GIF, WebP
**Max size:** 5MB

**Response (200 OK):**
```json
{
  "message": "Avatar uploaded successfully",
  "avatar_url": "/uploads/avatars/avatar_1_1699900000.jpg",
  "data": {
    "id": 1,
    "avatar_url": "/uploads/avatars/avatar_1_1699900000.jpg",
    ...
  }
}
```

---

#### DELETE /profile/avatar

Delete the authenticated user's avatar.

**Response (200 OK):**
```json
{
  "message": "Avatar deleted successfully",
  "data": {
    "id": 1,
    "avatar_url": null,
    ...
  }
}
```

---

### Public User Endpoints

#### GET /users

List public user profiles with pagination.

**Query Parameters:**
- `page` (default: 1) - Page number
- `limit` (default: 10, max: 100) - Results per page

**Response (200 OK):**
```json
{
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "johndoe",
        "email": "john@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "avatar_url": "/uploads/avatars/avatar_1.jpg",
        "bio": "Sports enthusiast",
        "location": "Buenos Aires",
        "fitness_level": "intermediate",
        "role": "user",
        "created_at": "2025-11-14T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

---

#### GET /users/:id

Get public profile of a specific user.

**Response (200 OK):**
```json
{
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "avatar_url": "/uploads/avatars/avatar_1.jpg",
    "bio": "Sports enthusiast",
    "location": "Buenos Aires",
    "role": "user",
    "created_at": "2025-11-14T10:00:00Z"
  }
}
```

---

### Admin Endpoints (Admin Role Required)

All admin endpoints require JWT with `admin`, `super_admin`, or `root` role.

#### GET /admin/users

List all users with full profiles and filtering options.

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20)
- `role` (filter by role: user, moderator, admin, etc.)
- `status` (filter by status: active, inactive)
- `search` (search by email or username)

**Response (200 OK):**
```json
{
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": 1,
        "email": "admin@example.com",
        "username": "admin",
        "role": "admin",
        "is_active": true,
        "email_verified": true,
        ...
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    },
    "filters": {
      "role": "admin",
      "status": "active",
      "search": ""
    }
  }
}
```

---

#### POST /admin/users

Create a new user (admin function).

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "username": "newuser",
  "password": "securepassword123",
  "first_name": "New",
  "last_name": "User",
  "role": "user",
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "message": "User created successfully",
  "data": {
    "id": 10,
    "email": "newuser@example.com",
    "username": "newuser",
    "role": "user",
    "is_active": true
  }
}
```

---

#### PUT /admin/users/:id/role

Update a user's role.

**Request Body:**
```json
{
  "role": "moderator"
}
```

**Valid roles:** user, moderator, admin, super_admin, root

**Response (200 OK):**
```json
{
  "message": "User role updated successfully",
  "data": {
    "id": 10,
    "email": "user@example.com",
    "role": "moderator"
  }
}
```

---

#### PUT /admin/users/:id/status

Activate or deactivate a user account.

**Request Body:**
```json
{
  "is_active": false
}
```

**Response (200 OK):**
```json
{
  "message": "User deactivated successfully",
  "data": {
    "id": 10,
    "email": "user@example.com",
    "is_active": false
  }
}
```

---

#### GET /admin/stats

Get system statistics.

**Response (200 OK):**
```json
{
  "message": "Statistics retrieved successfully",
  "data": {
    "TotalUsers": 150,
    "ActiveUsers": 145,
    "InactiveUsers": 5,
    "RootUsers": 1,
    "AdminUsers": 3,
    "ModeratorUsers": 10,
    "RegularUsers": 136
  }
}
```

---

### Root Endpoints (Root Role Required)

#### POST /admin/create-root

Create a root user (special endpoint, requires secret key).

**Request Body:**
```json
{
  "email": "root@example.com",
  "username": "root",
  "password": "supersecurepassword",
  "secret_key": "your-root-creation-secret"
}
```

**Note:** The `secret_key` must match the `ROOT_SECRET_KEY` environment variable.

**Response (201 Created):**
```json
{
  "message": "Root user created successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "root@example.com",
      "username": "root",
      "role": "root"
    }
  }
}
```

---

#### DELETE /admin/users/:id

Permanently delete a user (root only).

**Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

**Note:** This is a hard delete and cannot be undone.

---

## Authentication Flow

### Registration & Login Flow

```
┌──────────┐          ┌──────────────┐          ┌────────┐
│  Client  │          │  Users API   │          │ MySQL  │
└────┬─────┘          └──────┬───────┘          └───┬────┘
     │                       │                      │
     │ POST /auth/register   │                      │
     ├──────────────────────>│                      │
     │                       │ Hash password        │
     │                       │ (bcrypt)             │
     │                       ├─────────────────────>│
     │                       │ Create user record   │
     │                       │<─────────────────────┤
     │                       │ Generate JWT         │
     │                       │ (with role claim)    │
     │<──────────────────────┤                      │
     │ JWT Token + User data │                      │
     │                       │                      │
     │ POST /auth/login      │                      │
     ├──────────────────────>│                      │
     │                       ├─────────────────────>│
     │                       │ Verify credentials   │
     │                       │<─────────────────────┤
     │                       │ Generate JWT         │
     │<──────────────────────┤                      │
     │ JWT Token + User data │                      │
     │                       │                      │
```

### Protected Endpoint Access Flow

```
┌──────────┐          ┌──────────────┐
│  Client  │          │  Users API   │
└────┬─────┘          └──────┬───────┘
     │                       │
     │ GET /profile          │
     │ Authorization: Bearer │
     ├──────────────────────>│
     │                       │ Validate JWT
     │                       │ Extract user_id & role
     │                       │ Check role permissions
     │<──────────────────────┤
     │ User profile data     │
     │                       │
```

## Error Responses

The API uses consistent error response format:

```json
{
  "error": "Error type",
  "message": "Human-readable error message"
}
```

### Common HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, PUT requests |
| 201 | Created | Successful POST (resource created) |
| 400 | Bad Request | Invalid request data/validation errors |
| 401 | Unauthorized | Missing/invalid JWT or wrong credentials |
| 403 | Forbidden | Insufficient role permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate email/username |
| 500 | Internal Server Error | Unexpected server error |

## Environment Variables

Required environment variables for the Users API:

```bash
# Server
PORT=8081

# Database
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=rootpassword
DB_NAME=users_db

# JWT
JWT_SECRET=your-super-secret-jwt-key-here

# Root user creation
ROOT_SECRET_KEY=your-root-creation-secret

# CORS (optional)
CORS_ALLOWED_ORIGINS=*
```

## Service Dependencies

- **MySQL**: Required for user data storage
- **No external service dependencies**: Users API is fully self-contained

## Health Check

**Endpoint:** `GET /api/v1/health`

**Response:**
```json
{
  "status": "ok",
  "message": "Users API is running",
  "service": "users-api"
}
```

## Security Considerations

1. **Password Storage**: Passwords are hashed using bcrypt with default cost (10)
2. **JWT Expiration**: Tokens expire after 24 hours
3. **Role in JWT**: Role is included in JWT claims for efficient authorization
4. **HTTPS Required**: Always use HTTPS in production
5. **Rate Limiting**: Should be implemented at API Gateway level
6. **Input Validation**: All inputs are validated using Gin's binding validators

## Development

### Running Locally

```bash
cd backend/users-api

# Install dependencies
go mod tidy

# Run the service
go run main.go
```

### Running Tests

```bash
go test ./...
```

### Database Migrations

The service uses GORM Auto-Migrate, which automatically creates/updates tables on startup.

## Related Documentation

- [Activities API](./activities-api.md)
- [Search API](./search-api.md)
- [Reservations API](./reservations-api.md)
- [Architecture Overview](../architecture.md)

