# Reservations API Documentation

## Overview

The Reservations API is a Go-based microservice that manages activity reservations and bookings. It provides endpoints for creating, viewing, updating, and managing reservations for sports activities.

**Technology Stack:**
- Language: Go 1.21+
- Framework: Gin
- Database: MongoDB 6.0
- Database Driver: Official MongoDB Go Driver
- Message Queue: RabbitMQ (for event publishing)
- Authentication: JWT validation

**Port:** 8080

**Base URL:** `http://localhost:8080`

## Architecture

### Database Schema (MongoDB)

The service uses MongoDB to store reservations. Each reservation document follows this structure:

```javascript
{
  "_id": ObjectId("673590beeb2a7e80f9ff0c63"),
  "users_id": [1, 5, 12],  // Array of user IDs from Users API
  "cupo": 3,  // Number of spots reserved
  "actividad": "673590beeb2a7e80f9ff0c62",  // Activity ID (MongoDB ObjectID as string)
  "date": ISODate("2025-11-20T10:00:00Z"),
  "status": "confirmada",  // Status: pendiente, confirmada, cancelada
  "created_at": ISODate("2025-11-14T10:00:00Z"),
  "updated_at": ISODate("2025-11-14T10:00:00Z")
}
```

### Reservation Status Flow

```
    ┌──────────┐
    │          │
    │ Creación │
    │          │
    └────┬─────┘
         │
         v
   ┌──────────┐
   │          │
   │Pendiente │  <───────┐
   │          │          │
   └────┬─────┘          │
         │               │ Admin puede
         │               │ cambiar estado
         │ Confirmación  │
         │ automática o  │
         │ por admin     │
         v               │
   ┌──────────┐          │
   │          │          │
   │Confirmada│ ─────────┘
   │          │
   └────┬─────┘
         │
         │ Cancelación
         │ por usuario
         │ o admin
         v
   ┌──────────┐
   │          │
   │Cancelada │
   │          │
   └──────────┘
```

### Reservation Statuses

| Status | Description | User Can Cancel | Admin Can Modify |
|--------|-------------|-----------------|------------------|
| `pendiente` | Awaiting confirmation | ✅ Yes | ✅ Yes |
| `confirmada` | Confirmed reservation | ✅ Yes | ✅ Yes |
| `cancelada` | Cancelled reservation | ❌ No | ✅ Yes |

## API Endpoints

### Public Endpoints

#### GET /reservas

List reservations.

**Authentication:** Required (JWT token in Authorization header)

**Authorization:**
- Regular users see only their own reservations
- Admins see all reservations

**Response (200 OK):**
```json
{
  "reservas": [
    {
      "id": "673590beeb2a7e80f9ff0c63",
      "users_id": [1, 5, 12],
      "cupo": 3,
      "actividad": "673590beeb2a7e80f9ff0c62",
      "date": "2025-11-20T10:00:00Z",
      "status": "confirmada",
      "created_at": "2025-11-14T10:00:00Z",
      "updated_at": "2025-11-14T10:00:00Z"
    },
    {
      "id": "673590beeb2a7e80f9ff0c64",
      "users_id": [3],
      "cupo": 1,
      "actividad": "673590beeb2a7e80f9ff0c65",
      "date": "2025-11-21T15:00:00Z",
      "status": "pendiente",
      "created_at": "2025-11-14T11:00:00Z",
      "updated_at": "2025-11-14T11:00:00Z"
    }
  ]
}
```

---

#### GET /reservas/:id

Get a specific reservation by ID.

**Authentication:** Not required (public endpoint)

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Response (200 OK):**
```json
{
  "reserva": {
    "id": "673590beeb2a7e80f9ff0c63",
    "users_id": [1, 5, 12],
    "cupo": 3,
    "actividad": "673590beeb2a7e80f9ff0c62",
    "date": "2025-11-20T10:00:00Z",
    "status": "confirmada",
    "created_at": "2025-11-14T10:00:00Z",
    "updated_at": "2025-11-14T10:00:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Reserva no encontrada"
}
```

---

### Protected Endpoints (Requires JWT)

#### POST /reservas

Create a new reservation.

**Authentication:** Required (JWT token in Authorization header)

**Request Body:**
```json
{
  "users_id": [1, 5, 12],
  "cupo": 3,
  "actividad": "673590beeb2a7e80f9ff0c62",
  "date": "2025-11-20T10:00:00Z",
  "status": "pendiente"
}
```

**Field Descriptions:**
- `users_id` (required, array of integers) - IDs of users making the reservation
- `cupo` (required, integer, min: 1) - Number of spots to reserve
- `actividad` (required, string) - MongoDB ObjectID of the activity
- `date` (required, ISO 8601 date) - Reservation date and time
- `status` (optional, default: "pendiente") - Initial status

**Validation Rules:**
- `cupo` must be greater than 0
- `cupo` must not exceed activity's `max_capacity`
- `users_id` array length should match `cupo` (one user per spot)
- `actividad` must reference an existing, active activity
- `date` should be in the future

**Response (201 Created):**
```json
{
  "message": "Reserva creada exitosamente",
  "reserva": {
    "id": "673590beeb2a7e80f9ff0c66",
    "users_id": [1, 5, 12],
    "cupo": 3,
    "actividad": "673590beeb2a7e80f9ff0c62",
    "date": "2025-11-20T10:00:00Z",
    "status": "pendiente",
    "created_at": "2025-11-14T15:30:00Z",
    "updated_at": "2025-11-14T15:30:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Cupo must be greater than 0"
}
```

**Error Response (409 Conflict):**
```json
{
  "error": "No hay cupos disponibles para esta actividad"
}
```

---

### Admin Endpoints (Admin Role Required)

All admin endpoints require a valid JWT token with `admin`, `super_admin`, or `root` role.

**Authentication:** Include JWT in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

#### PUT /reservas/:id

Update an existing reservation (admin only).

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Request Body (all fields optional):**
```json
{
  "users_id": [1, 5, 12, 8],
  "cupo": 4,
  "date": "2025-11-20T11:00:00Z",
  "status": "confirmada"
}
```

**Response (200 OK):**
```json
{
  "message": "Reserva actualizada exitosamente",
  "reserva": {
    "id": "673590beeb2a7e80f9ff0c63",
    "users_id": [1, 5, 12, 8],
    "cupo": 4,
    "actividad": "673590beeb2a7e80f9ff0c62",
    "date": "2025-11-20T11:00:00Z",
    "status": "confirmada",
    "created_at": "2025-11-14T10:00:00Z",
    "updated_at": "2025-11-14T16:00:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Reserva no encontrada"
}
```

---

#### DELETE /reservas/:id

Delete/cancel a reservation.

**Authentication:** Required (JWT token in Authorization header)

**Authorization:**
- Users can delete their own reservations
- Admins can delete any reservation

**Parameters:**
- `id` (path parameter) - MongoDB ObjectID

**Response (200 OK):**
```json
{
  "message": "Reserva eliminada exitosamente"
}
```

**Error Response (403 Forbidden):**
```json
{
  "error": "You can only delete your own reservations"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Reserva not found"
}
```

---

### Health Check

#### GET /healthz

Check if the service is running and healthy.

**Response (200 OK):**
```json
{
  "status": "ok",
  "service": "reservations-api",
  "database": "connected",
  "rabbitmq": "connected"
}
```

## Authentication & Authorization

### JWT Validation

The Reservations API validates JWT tokens issued by the Users API. The JWT must contain:

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "role": "user"
}
```

### Role Requirements

| Endpoint | Public | User | Admin |
|----------|--------|------|-------|
| GET /reservas | ❌ | ✅ (own only) | ✅ (all) |
| GET /reservas/:id | ✅ | ✅ | ✅ |
| POST /reservas | ❌ | ✅ | ✅ |
| PUT /reservas/:id | ❌ | ❌ | ✅ |
| DELETE /reservas/:id | ❌ | ✅ (own) | ✅ (any) |

**Notes:**
- `GET /reservas`: Regular users see only their own reservations. Admins see all reservations.
- `DELETE /reservas/:id`: Users can only delete their own reservations. Admins can delete any reservation.

### Admin Middleware

Admin endpoints are protected by middleware that:
1. Validates the JWT token
2. Extracts the `role` claim
3. Checks if role is `admin`, `super_admin`, or `root`
4. Rejects requests with `403 Forbidden` if role is insufficient

## Event Publishing

### RabbitMQ Integration

The Reservations API can publish events to RabbitMQ for integration with other services (currently minimal implementation).

**Potential Events:**
- `reservation.created` - When a new reservation is created
- `reservation.updated` - When a reservation is modified
- `reservation.cancelled` - When a reservation is cancelled

**Future Integration:** These events could be consumed by:
- **Notifications Service** - Send email/SMS confirmations
- **Analytics Service** - Track booking patterns
- **Activities API** - Update available capacity

## Business Logic

### Capacity Management

When creating a reservation:
1. Fetch the activity from Activities API
2. Check `max_capacity`
3. Count existing confirmed reservations for the activity
4. Verify: `existing_reservations + new_cupo <= max_capacity`
5. Reject if capacity exceeded

### Reservation Rules

1. **Capacity Limit**: Total `cupo` across all confirmed reservations cannot exceed activity's `max_capacity`
2. **User Validation**: Users in `users_id` should exist in Users API (optional validation)
3. **Date Validation**: Reservation `date` should align with activity's `schedule`
4. **Status Transitions**:
   - `pendiente` → `confirmada` (user/admin)
   - `pendiente` → `cancelada` (user/admin)
   - `confirmada` → `cancelada` (user/admin)
   - `cancelada` → (no transitions, final state)

### Data Access Object (DAO) Pattern

Similar to Activities API, Reservations API uses the DAO pattern:

#### Domain Model (domain/reserva.go)

```go
type Reserva struct {
    ID        string    `json:"id"`
    UsersID   []int     `json:"users_id"`
    Cupo      int       `json:"cupo"`
    Actividad string    `json:"actividad"`
    Date      time.Time `json:"date"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### DAO Model (dao/reserva.go)

```go
type Reserva struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    UsersID   []int              `bson:"users_id"`
    Cupo      int                `bson:"cupo"`
    Actividad string             `bson:"actividad"`
    Date      time.Time          `bson:"date"`
    Status    string             `bson:"status"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}
```

## Integration with Other Services

### Activities API Integration

The Reservations API integrates with Activities API to:
1. **Validate Activity Existence**: Check that `actividad` ID exists
2. **Check Capacity**: Fetch `max_capacity` from the activity
3. **Verify Status**: Ensure activity is active

**Example Request:**
```go
activityResponse := GET("http://activities-api:8082/activities/" + activityID)
if activityResponse.StatusCode == 404 {
    return Error("Activity not found")
}
activity := activityResponse.Data
if !activity.IsActive {
    return Error("Activity is not active")
}
```

### Users API Integration

Optional integration to:
1. Verify users exist
2. Fetch user details for notifications
3. Check user permissions

## Error Handling

### Common HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, PUT |
| 201 | Created | Successful POST |
| 400 | Bad Request | Invalid request data, validation errors |
| 401 | Unauthorized | Missing/invalid JWT |
| 403 | Forbidden | Insufficient role permissions |
| 404 | Not Found | Reservation or activity doesn't exist |
| 409 | Conflict | Capacity exceeded, constraint violation |
| 500 | Internal Server Error | Database or service error |

### Error Format

```json
{
  "error": "Error message description"
}
```

## Environment Variables

```bash
# Server
PORT=8080

# MongoDB
MONGO_URI=mongodb://mongodb:27017
MONGO_DATABASE=reservations_db
MONGO_COLLECTION=reservas

# RabbitMQ (optional, for event publishing)
RABBITMQ_URL=amqp://admin:admin123@rabbitmq:5672/
RABBITMQ_EXCHANGE=reservations_exchange

# JWT
JWT_SECRET=your-super-secret-jwt-key-here

# External Services
ACTIVITIES_API_URL=http://activities-api:8082
USERS_API_URL=http://users-api:8081

# CORS
CORS_ALLOWED_ORIGINS=*
```

## Service Dependencies

### Required Services

1. **MongoDB** - Primary database for storing reservations
   - Host: `mongodb:27017`
   - Database: `reservations_db`
   - Collection: `reservas`

### Optional Services

1. **RabbitMQ** - For event publishing (future feature)
   - Host: `rabbitmq:5672`

2. **Activities API** - For activity validation and capacity checks
   - URL: `http://activities-api:8082`

3. **Users API** - For user validation (optional)
   - URL: `http://users-api:8081`

## Development

### Running Locally

```bash
cd backend/reservations-api

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

Sample reservations can be seeded using:

```bash
# From project root
./scripts/seed-data.sh
```

## Example Usage

### Create a Reservation (Authenticated User)

```bash
# Login as user first
TOKEN=$(curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.data.token')

# Create reservation
curl -X POST http://localhost:8080/reservas \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "users_id": [1],
    "cupo": 1,
    "actividad": "673590beeb2a7e80f9ff0c62",
    "date": "2025-11-20T10:00:00Z",
    "status": "pendiente"
  }'
```

### List All Reservations (Public)

```bash
curl http://localhost:8080/reservas
```

### Get Specific Reservation

```bash
curl http://localhost:8080/reservas/673590beeb2a7e80f9ff0c63
```

### Update Reservation Status (Admin)

```bash
# Login as admin
ADMIN_TOKEN=$(curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  | jq -r '.data.token')

# Update reservation
curl -X PUT http://localhost:8080/reservas/673590beeb2a7e80f9ff0c63 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "status": "confirmada"
  }'
```

### Cancel Reservation (Admin)

```bash
curl -X DELETE http://localhost:8080/reservas/673590beeb2a7e80f9ff0c63 \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Future Enhancements

1. **Email Notifications**: Send confirmation emails when reservations are created/updated
2. **SMS Notifications**: Send SMS reminders before activity date
3. **Waiting List**: Implement waiting list when activity is full
4. **Recurring Reservations**: Support for booking multiple sessions
5. **Payment Integration**: Integrate payment gateway for paid activities
6. **Cancellation Policy**: Implement cancellation deadlines and refund rules
7. **User Dashboard**: Endpoint to get user's own reservations only
8. **Capacity Real-time Updates**: Emit events when capacity changes
9. **Reservation Expiry**: Auto-cancel pending reservations after timeout
10. **Check-in System**: Track user attendance for reservations

## Related Documentation

- [Users API](./users-api.md)
- [Activities API](./activities-api.md)
- [Search API](./search-api.md)
- [Architecture Overview](../architecture.md)

