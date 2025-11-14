# Activities API Documentation

## Overview

The Activities API is a Go-based microservice that manages sports activities. It provides CRUD operations for activities and integrates with RabbitMQ for event-driven communication with other services.

**Technology Stack:**
- Language: Go 1.21+
- Framework: Gin
- Database: MongoDB 6.0
- Database Driver: Official MongoDB Go Driver
- Message Queue: RabbitMQ (for event publishing)
- Authentication: JWT validation

**Port:** 8082

**Base URL:** `http://localhost:8082`

## Architecture

### Database Schema (MongoDB)

The service uses MongoDB to store activities. Each activity document follows this structure:

```javascript
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
  "name": "Morning Yoga Session",
  "description": "Relaxing yoga session for all levels",
  "category": "yoga",
  "difficulty": "beginner",
  "location": "Palermo Park, Buenos Aires",
  "price": 1500.00,
  "duration": 90,  // in minutes
  "max_capacity": 15,
  "instructor": "Maria Rodriguez",
  "schedule": [
    "Monday 08:00-09:30",
    "Wednesday 08:00-09:30",
    "Friday 08:00-09:30"
  ],
  "equipment": ["yoga mat", "towel", "water bottle"],
  "image_url": "https://example.com/images/yoga.jpg",
  "is_active": true,
  "created_by": 1,  // User ID from Users API
  "created_at": ISODate("2025-11-14T10:00:00Z"),
  "updated_at": ISODate("2025-11-14T10:00:00Z")
}
```

### Activity Categories

The following categories are supported:
- `football` - Soccer/football
- `basketball` - Basketball
- `tennis` - Tennis
- `running` - Running activities
- `yoga` - Yoga classes
- `swimming` - Swimming
- `cycling` - Cycling
- `gym` - Gym workouts
- `martial_arts` - Martial arts
- `dance` - Dance classes
- `other` - Other activities

### Difficulty Levels

- `beginner` - Suitable for beginners
- `intermediate` - Requires some experience
- `advanced` - For experienced participants
- `professional` - Professional/competitive level

### Event-Driven Architecture

The Activities API publishes events to RabbitMQ when activities are created, updated, or deleted. This enables other services (like Search API) to react to changes.

**RabbitMQ Exchange:** `activities_exchange` (type: topic)

**Events Published:**

1. **activity.created**
   - Routing key: `activity.created`
   - Payload: Full activity object
   - Purpose: Notify other services that a new activity was created

2. **activity.updated**
   - Routing key: `activity.updated`
   - Payload: Full activity object
   - Purpose: Notify other services about activity updates

3. **activity.deleted**
   - Routing key: `activity.deleted`
   - Payload: `{ "id": "507f1f77bcf86cd799439011" }`
   - Purpose: Notify other services that an activity was deleted

## API Endpoints

### Public Endpoints

#### GET /activities

List all active activities.

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "673590beeb2a7e80f9ff0c63",
      "name": "Morning Yoga Session",
      "description": "Relaxing yoga session for all levels",
      "category": "yoga",
      "difficulty": "beginner",
      "location": "Palermo Park, Buenos Aires",
      "price": 1500.00,
      "duration": 90,
      "max_capacity": 15,
      "instructor": "Maria Rodriguez",
      "schedule": ["Monday 08:00-09:30", "Wednesday 08:00-09:30"],
      "equipment": ["yoga mat", "towel"],
      "image_url": "https://example.com/images/yoga.jpg",
      "is_active": true,
      "created_by": 1,
      "created_at": "2025-11-14T10:00:00Z",
      "updated_at": "2025-11-14T10:00:00Z"
    }
  ]
}
```

---

#### GET /activities/:id

Get a specific activity by ID.

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Response (200 OK):**
```json
{
  "data": {
    "id": "673590beeb2a7e80f9ff0c63",
    "name": "Morning Yoga Session",
    "description": "Relaxing yoga session for all levels",
    "category": "yoga",
    "difficulty": "beginner",
    "location": "Palermo Park, Buenos Aires",
    "price": 1500.00,
    "duration": 90,
    "max_capacity": 15,
    "instructor": "Maria Rodriguez",
    "schedule": ["Monday 08:00-09:30"],
    "equipment": ["yoga mat", "towel"],
    "image_url": "https://example.com/images/yoga.jpg",
    "is_active": true,
    "created_by": 1,
    "created_at": "2025-11-14T10:00:00Z",
    "updated_at": "2025-11-14T10:00:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Activity not found"
}
```

---

#### GET /activities/category/:category

Get all active activities in a specific category.

**Parameters:**
- `category` (path parameter) - One of: football, basketball, tennis, yoga, running, etc.

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "673590beeb2a7e80f9ff0c63",
      "name": "Morning Yoga Session",
      "category": "yoga",
      ...
    },
    {
      "id": "673590beeb2a7e80f9ff0c64",
      "name": "Evening Yoga Flow",
      "category": "yoga",
      ...
    }
  ]
}
```

---

### Admin Endpoints (Admin Role Required)

All admin endpoints require a valid JWT token with `admin`, `super_admin`, or `root` role.

**Authentication:** Include JWT in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

#### GET /activities/all

List all activities (including inactive).

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "673590beeb2a7e80f9ff0c63",
      "name": "Morning Yoga Session",
      "is_active": true,
      ...
    },
    {
      "id": "673590beeb2a7e80f9ff0c64",
      "name": "Cancelled Football Match",
      "is_active": false,
      ...
    }
  ]
}
```

---

#### POST /activities

Create a new activity.

**Request Body:**
```json
{
  "name": "Beach Volleyball Tournament",
  "description": "Competitive beach volleyball for intermediate players",
  "category": "volleyball",
  "difficulty": "intermediate",
  "location": "Playa Grande, Mar del Plata",
  "price": 2500.00,
  "duration": 120,
  "max_capacity": 20,
  "instructor": "Carlos Martinez",
  "schedule": [
    "Saturday 16:00-18:00",
    "Sunday 16:00-18:00"
  ],
  "equipment": ["volleyball", "net", "sunscreen"],
  "image_url": "https://example.com/images/beach-volleyball.jpg"
}
```

**Response (201 Created):**
```json
{
  "message": "Activity created successfully",
  "data": {
    "id": "673590beeb2a7e80f9ff0c65",
    "name": "Beach Volleyball Tournament",
    "description": "Competitive beach volleyball for intermediate players",
    "category": "volleyball",
    "difficulty": "intermediate",
    "location": "Playa Grande, Mar del Plata",
    "price": 2500.00,
    "duration": 120,
    "max_capacity": 20,
    "instructor": "Carlos Martinez",
    "schedule": ["Saturday 16:00-18:00", "Sunday 16:00-18:00"],
    "equipment": ["volleyball", "net", "sunscreen"],
    "image_url": "https://example.com/images/beach-volleyball.jpg",
    "is_active": true,
    "created_by": 1,
    "created_at": "2025-11-14T15:30:00Z",
    "updated_at": "2025-11-14T15:30:00Z"
  }
}
```

**Validation:**
- `name` (required, min: 3, max: 200)
- `description` (required, min: 10)
- `category` (required)
- `difficulty` (required)
- `location` (required)
- `price` (required, >= 0)
- `duration` (required, > 0)
- `max_capacity` (required, > 0)
- `instructor` (required)

**Event:** Publishes `activity.created` event to RabbitMQ.

---

#### PUT /activities/:id

Update an existing activity.

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Request Body (all fields optional):**
```json
{
  "name": "Updated Activity Name",
  "description": "Updated description",
  "price": 1800.00,
  "max_capacity": 20,
  "schedule": ["Monday 10:00-11:30", "Friday 10:00-11:30"]
}
```

**Response (200 OK):**
```json
{
  "message": "Activity updated successfully",
  "data": {
    "id": "673590beeb2a7e80f9ff0c63",
    "name": "Updated Activity Name",
    "description": "Updated description",
    ...
  }
}
```

**Event:** Publishes `activity.updated` event to RabbitMQ.

---

#### DELETE /activities/:id

Delete an activity (soft delete - sets `is_active` to false).

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Response (200 OK):**
```json
{
  "message": "Activity deleted successfully"
}
```

**Event:** Publishes `activity.deleted` event to RabbitMQ.

---

#### PATCH /activities/:id/toggle

Toggle activity active status.

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Response (200 OK):**
```json
{
  "message": "Activity status toggled successfully",
  "data": {
    "id": "673590beeb2a7e80f9ff0c63",
    "is_active": false,
    ...
  }
}
```

**Event:** Publishes `activity.updated` event to RabbitMQ.

---

### Health Check

#### GET /healthz

Check if the service is running and healthy.

**Response (200 OK):**
```json
{
  "status": "ok",
  "service": "activities-api",
  "database": "connected",
  "rabbitmq": "connected"
}
```

## Authentication & Authorization

### JWT Validation

The Activities API validates JWT tokens issued by the Users API. The JWT must contain:

```json
{
  "user_id": 1,
  "email": "admin@example.com",
  "username": "admin",
  "role": "admin"
}
```

### Role Requirements

| Endpoint | Public | User | Admin |
|----------|--------|------|-------|
| GET /activities | ✅ | ✅ | ✅ |
| GET /activities/:id | ✅ | ✅ | ✅ |
| GET /activities/category/:category | ✅ | ✅ | ✅ |
| GET /activities/all | ❌ | ❌ | ✅ |
| POST /activities | ❌ | ❌ | ✅ |
| PUT /activities/:id | ❌ | ❌ | ✅ |
| DELETE /activities/:id | ❌ | ❌ | ✅ |
| PATCH /activities/:id/toggle | ❌ | ❌ | ✅ |

### Admin Middleware

Admin endpoints are protected by middleware that:
1. Validates the JWT token
2. Extracts the `role` claim
3. Checks if role is `admin`, `super_admin`, or `root`
4. Rejects requests with `403 Forbidden` if role is insufficient

## Event Publishing Flow

### Activity Creation Flow

```
┌──────────┐       ┌─────────────────┐       ┌─────────┐       ┌────────────┐
│  Admin   │       │ Activities API  │       │ MongoDB │       │  RabbitMQ  │
└────┬─────┘       └────────┬────────┘       └────┬────┘       └──────┬─────┘
     │                      │                     │                   │
     │ POST /activities     │                     │                   │
     ├─────────────────────>│                     │                   │
     │                      │ Validate JWT        │                   │
     │                      │ Check admin role    │                   │
     │                      │                     │                   │
     │                      │ Insert document     │                   │
     │                      ├────────────────────>│                   │
     │                      │<────────────────────┤                   │
     │                      │ Activity created    │                   │
     │                      │                     │                   │
     │                      │ Publish event       │                   │
     │                      │ (activity.created)  │                   │
     │                      ├─────────────────────────────────────────>│
     │                      │                     │                   │
     │<─────────────────────┤                     │                   │
     │ 201 Created          │                     │                   │
     │                      │                     │                   │
```

### Event Consumer

The **Search API** consumes these events to keep its Solr index synchronized with the activities database.

## Data Access Object (DAO) Pattern

The Activities API uses the DAO pattern to separate database concerns from business logic:

### Domain Model (domain/activity.go)

Business logic representation with string ID:

```go
type Activity struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Difficulty  string    `json:"difficulty"`
    Location    string    `json:"location"`
    Price       float64   `json:"price"`
    Duration    int       `json:"duration"`
    MaxCapacity int       `json:"max_capacity"`
    Instructor  string    `json:"instructor"`
    Schedule    []string  `json:"schedule"`
    Equipment   []string  `json:"equipment"`
    ImageURL    string    `json:"image_url"`
    IsActive    bool      `json:"is_active"`
    CreatedBy   uint      `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### DAO Model (dao/activity.go)

Database representation with MongoDB ObjectID:

```go
type Activity struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`
    Description string             `bson:"description"`
    Category    string             `bson:"category"`
    Difficulty  string             `bson:"difficulty"`
    Location    string             `bson:"location"`
    Price       float64            `bson:"price"`
    Duration    int                `bson:"duration"`
    MaxCapacity int                `bson:"max_capacity"`
    Instructor  string             `bson:"instructor"`
    Schedule    []string           `bson:"schedule"`
    Equipment   []string           `bson:"equipment"`
    ImageURL    string             `bson:"image_url"`
    IsActive    bool               `bson:"is_active"`
    CreatedBy   uint               `bson:"created_by"`
    CreatedAt   time.Time          `bson:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at"`
}
```

**Conversion Functions:**
- `ToDomain()` - Converts DAO model to domain model
- `FromDomain()` - Converts domain model to DAO model

## Error Responses

### Common HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST |
| 400 | Bad Request | Invalid request data |
| 401 | Unauthorized | Missing/invalid JWT |
| 403 | Forbidden | Insufficient role permissions |
| 404 | Not Found | Activity doesn't exist |
| 500 | Internal Server Error | Database or RabbitMQ error |

### Error Format

```json
{
  "error": "Error message description"
}
```

## Environment Variables

```bash
# Server
PORT=8082

# MongoDB
MONGO_URI=mongodb://mongodb:27017
MONGO_DATABASE=activities_db
MONGO_COLLECTION=activities

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin123@rabbitmq:5672/
RABBITMQ_EXCHANGE=activities_exchange

# JWT
JWT_SECRET=your-super-secret-jwt-key-here

# CORS
CORS_ALLOWED_ORIGINS=*

# Users API (for user verification, currently unused in favor of JWT)
USERS_API_URL=http://users-api:8081
```

## Service Dependencies

### Required Services

1. **MongoDB** - Primary database for storing activities
   - Host: `mongodb:27017`
   - Database: `activities_db`
   - Collection: `activities`

2. **RabbitMQ** - Message broker for event publishing
   - Host: `rabbitmq:5672`
   - Exchange: `activities_exchange`
   - Type: topic

### Optional Services

1. **Users API** - For extended user information (not currently used)

## Integration with Other Services

### Search API Integration

The Search API subscribes to activity events:

1. **activity.created** → Indexes new activity in Solr
2. **activity.updated** → Re-indexes updated activity in Solr
3. **activity.deleted** → Removes activity from Solr index

### Reservations API Integration

The Reservations API references activities by their MongoDB ObjectID to create bookings.

## Development

### Running Locally

```bash
cd backend/activities-api

# Install dependencies
go mod tidy

# Run the service
go run main.go
```

### Running Tests

```bash
go test ./...
```

### Seeding Data

Sample activities can be seeded using:

```bash
# From project root
./scripts/seed-data.sh
```

## Example Usage

### Create an Activity (Admin)

```bash
# Login as admin first
TOKEN=$(curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  | jq -r '.data.token')

# Create activity
curl -X POST http://localhost:8082/activities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Advanced Tennis Training",
    "description": "High-intensity tennis training for advanced players",
    "category": "tennis",
    "difficulty": "advanced",
    "location": "Buenos Aires Lawn Tennis Club",
    "price": 3500.00,
    "duration": 120,
    "max_capacity": 8,
    "instructor": "Rafael Gonzalez",
    "schedule": ["Tuesday 18:00-20:00", "Thursday 18:00-20:00"],
    "equipment": ["tennis racket", "tennis balls", "sports shoes"],
    "image_url": "https://example.com/images/tennis.jpg"
  }'
```

### List All Activities (Public)

```bash
curl http://localhost:8082/activities
```

### Get Activity by Category

```bash
curl http://localhost:8082/activities/category/yoga
```

## Related Documentation

- [Users API](./users-api.md)
- [Search API](./search-api.md)
- [Reservations API](./reservations-api.md)
- [Architecture Overview](../architecture.md)

