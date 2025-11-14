# System Architecture Documentation

## Table of Contents

- [Overview](#overview)
- [Architectural Patterns](#architectural-patterns)
- [System Components](#system-components)
- [Communication Patterns](#communication-patterns)
- [Data Flow](#data-flow)
- [Authentication & Authorization Flow](#authentication--authorization-flow)
- [Event-Driven Architecture](#event-driven-architecture)
- [Caching Strategy](#caching-strategy)
- [Database Design](#database-design)
- [Deployment Architecture](#deployment-architecture)
- [Security Considerations](#security-considerations)
- [Scalability & Performance](#scalability--performance)

## Overview

The Sports Activities Platform is a microservices-based application built for managing sports activities, user profiles, reservations, and providing search capabilities. The system is designed with modern architectural patterns including microservices, event-driven communication, CQRS, and distributed caching.

### Technology Stack

**Backend Services:**
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Databases**: MySQL 8.0 (relational), MongoDB 6.0 (document)
- **Search Engine**: Apache Solr 9.4
- **Message Broker**: RabbitMQ 3 with Management Plugin
- **Cache**: Memcached 1.6

**Frontend:**
- **Framework**: React 18+ with TypeScript
- **UI Library**: Material-UI (MUI)
- **Build Tool**: Create React App
- **Web Server**: Nginx (production)

**Infrastructure:**
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Docker Compose (development), Kubernetes-ready
- **Networking**: Docker bridge network

## Architectural Patterns

### 1. Microservices Architecture

The system is decomposed into four independent microservices, each with a single responsibility:

```
┌──────────────────────────────────────────────────────────────────┐
│                        Sports Activities Platform                 │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────┐│
│  │   Users     │  │ Activities  │  │   Search    │  │Reservat.││
│  │   API       │  │     API     │  │    API      │  │   API   ││
│  │             │  │             │  │             │  │         ││
│  │  Port 8081  │  │  Port 8082  │  │  Port 8083  │  │Port 8080││
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────┘│
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

**Benefits:**
- **Independence**: Services can be developed, deployed, and scaled independently
- **Technology Diversity**: Each service can use the most appropriate technology
- **Fault Isolation**: Failure in one service doesn't bring down the entire system
- **Team Autonomy**: Different teams can own different services

### 2. Event-Driven Architecture (EDA)

Services communicate asynchronously through RabbitMQ events:

```
Activities API ──> RabbitMQ ──> Search API
                     │
                     └──────> (Future: Notifications, Analytics)
```

**Event Flow:**
1. **Activities API** publishes events when activities change
2. **RabbitMQ** routes events to interested consumers
3. **Search API** consumes events and updates Solr index
4. System remains decoupled and responsive

**Event Types:**
- `activity.created` - New activity added
- `activity.updated` - Activity modified
- `activity.deleted` - Activity removed

### 3. CQRS (Command Query Responsibility Segregation)

The system separates read and write operations:

**Write Side (Command):**
- **Activities API** handles all activity mutations (CREATE, UPDATE, DELETE)
- Uses MongoDB for flexible document storage
- Publishes events for each state change

**Read Side (Query):**
- **Search API** handles all search queries
- Uses Apache Solr optimized for full-text search
- Maintains read-optimized index synchronized via events

**Benefits:**
- **Performance**: Each side optimized for its purpose
- **Scalability**: Can scale reads and writes independently
- **Flexibility**: Different data models for reading vs writing

### 4. DAO (Data Access Object) Pattern

Services use DAO pattern to separate database concerns:

```
Controller ──> Service ──> Repository ──> DAO ──> Database
                                  │
                                  └──> Domain Models
```

**Example (Activities API):**

```go
// Domain Model (business logic)
type Activity struct {
    ID   string  // String representation
    Name string
    ...
}

// DAO Model (database)
type ActivityDAO struct {
    ID   primitive.ObjectID  // MongoDB ObjectID
    Name string
    ...
}

// Conversion
func (dao *ActivityDAO) ToDomain() Activity { ... }
func FromDomain(domain Activity) ActivityDAO { ... }
```

**Benefits:**
- **Separation of Concerns**: Business logic separate from data persistence
- **Testability**: Easy to mock database layer
- **Flexibility**: Can change database without affecting business logic

### 5. API Gateway Pattern (Frontend)

The React frontend acts as an API gateway for browser clients:

```
Browser ──> React Frontend ──> Backend Services
                 │
                 ├──> Users API (auth, profiles)
                 ├──> Activities API (CRUD operations)
                 ├──> Search API (search queries)
                 └──> Reservations API (bookings)
```

## System Components

### Frontend Service (React + Nginx)

**Port**: 3000

**Responsibilities:**
- User interface and user experience
- Client-side routing
- State management (React Context)
- API communication
- JWT token storage and management

**Key Features:**
- Material-UI components for consistent design
- TypeScript for type safety
- Responsive design for mobile and desktop
- Protected routes based on authentication state
- Role-based UI rendering (admin features)

---

### Users API

**Port**: 8081  
**Database**: MySQL  
**Documentation**: [Users API Docs](./api/users-api.md)

**Responsibilities:**
- User registration and authentication
- JWT token generation and validation
- User profile management (extended profiles)
- Role-based access control (RBAC)
- Avatar upload and storage

**Key Features:**
- bcrypt password hashing
- JWT with role claims
- Hierarchical role system (user, moderator, admin, super_admin, root)
- Extended user profiles (bio, sports interests, fitness level)
- Public vs private profile data

**Database Schema:**
- Single `users` table with relational integrity
- Indexed on email, username, role
- JSON fields for social_links and sports_interests

---

### Activities API

**Port**: 8082  
**Database**: MongoDB  
**Message Queue**: RabbitMQ (publisher)  
**Documentation**: [Activities API Docs](./api/activities-api.md)

**Responsibilities:**
- Sports activities CRUD operations
- Activity categorization and metadata
- Publishing events to RabbitMQ
- Admin-only write operations

**Key Features:**
- Document-based storage (MongoDB)
- Flexible activity schema
- Event publishing on state changes
- Soft delete (is_active flag)
- Admin middleware for protected endpoints

**Data Model:**
- Activities with metadata (category, difficulty, location)
- Embedded arrays (schedule, equipment)
- User reference (created_by)
- Timestamps for auditing

---

### Search API

**Port**: 8083  
**Search Engine**: Apache Solr  
**Cache**: Memcached  
**Message Queue**: RabbitMQ (consumer)  
**Documentation**: [Search API Docs](./api/search-api.md)

**Responsibilities:**
- Full-text search across activities
- Multi-faceted filtering (category, difficulty, price, location)
- Caching search results
- Consuming activity events and updating search index

**Key Features:**
- Apache Solr for powerful search capabilities
- Memcached for sub-10ms response times on cached queries
- Event-driven index synchronization
- Support for complex queries with filters
- Pagination support

**Architecture Highlights:**
- **CQRS Read Side**: Optimized for queries, not writes
- **Two-Level Caching**: Memcached + Solr internal caches
- **Automatic Sync**: Listens to RabbitMQ events
- **Search Fields**: name, description, location, instructor

---

### Reservations API

**Port**: 8080  
**Database**: MongoDB  
**Documentation**: [Reservations API Docs](./api/reservations-api.md)

**Responsibilities:**
- Activity reservation management
- Capacity tracking and validation
- Reservation status workflow
- User and admin booking operations

**Key Features:**
- Multi-user reservations (groups)
- Status workflow (pendiente → confirmada → cancelada)
- Capacity validation against activities
- Integration with Activities API
- Admin override capabilities

**Data Model:**
- Reservations with user arrays
- Activity reference (ObjectID)
- Status tracking
- Date and capacity management

---

### Infrastructure Services

#### MySQL Database

**Port**: 3307  
**Purpose**: Relational data for Users API

**Features:**
- Persistent volume for data storage
- Initialization scripts (`init.sql`)
- Seed data for development
- Indexed columns for performance

#### MongoDB Database

**Port**: 27017  
**Purpose**: Document storage for Activities and Reservations

**Features:**
- Shared by Activities API and Reservations API
- Separate databases per service
- Flexible schema
- Initialization scripts

#### Apache Solr

**Port**: 8983  
**Purpose**: Full-text search engine

**Features:**
- Custom `activities` core
- Configurable schema (schema.xml)
- Faceted search support
- Web-based admin UI (http://localhost:8983/solr)

#### RabbitMQ

**Port**: 5672 (AMQP), 15672 (Management UI)  
**Purpose**: Message broker for event-driven communication

**Features:**
- Topic exchange (`activities_exchange`)
- Durable queues
- Message persistence
- Management UI for monitoring

#### Memcached

**Port**: 11211  
**Purpose**: Distributed in-memory cache

**Features:**
- Key-value storage
- TTL-based expiration
- LRU eviction policy
- Shared cache across services

## Communication Patterns

### Synchronous Communication (HTTP/REST)

Used for:
- Frontend ↔ Backend APIs
- Cross-service communication (e.g., Reservations → Activities)

**Characteristics:**
- Request-response pattern
- Immediate feedback
- Strong coupling between caller and callee

**Example:**
```
GET http://activities-api:8082/activities/673590beeb2a7e80f9ff0c62
```

### Asynchronous Communication (RabbitMQ Events)

Used for:
- Activities API → Search API (index synchronization)
- Future: Notifications, Analytics

**Characteristics:**
- Fire-and-forget pattern
- Loose coupling
- Scalable and resilient

**Example Event:**
```json
{
  "event": "activity.created",
  "payload": {
    "id": "673590beeb2a7e80f9ff0c62",
    "name": "Morning Yoga",
    ...
  }
}
```

### Caching (Memcached)

Used for:
- Search API result caching
- Future: User session caching, activity caching

**Characteristics:**
- Key-value access
- Sub-millisecond latency
- TTL-based expiration

## Data Flow

### Activity Creation Flow

```
┌──────────┐
│  Admin   │
│ (Browser)│
└────┬─────┘
     │ 1. POST /activities
     │    Authorization: Bearer <JWT>
     │    { name, description, ... }
     v
┌────────────┐
│  Frontend  │
│  (React)   │
└────┬───────┘
     │ 2. HTTP POST with JWT
     v
┌──────────────┐         ┌─────────┐
│ Activities   │   3.    │  Users  │
│    API       │──────>  │  API    │
│              │  Verify │         │
│  (Port 8082) │   JWT   │(Port    │
└──────┬───────┘         │ 8081)   │
       │                 └─────────┘
       │ 4. Validate admin role
       │    (from JWT claims)
       │
       │ 5. Insert activity
       v
   ┌─────────┐
   │ MongoDB │
   └────┬────┘
        │ 6. Activity created
        │
        v
   ┌────────────────┐
   │ Activities API │
   └────────┬───────┘
            │
            │ 7. Publish event
            │    "activity.created"
            v
       ┌─────────┐
       │RabbitMQ │
       │Exchange │
       └────┬────┘
            │
            │ 8. Route event
            v
       ┌─────────┐
       │  Queue  │
       └────┬────┘
            │
            │ 9. Consume event
            v
    ┌──────────────┐
    │  Search API  │
    │ (Port 8083)  │
    └──────┬───────┘
           │
           │ 10. Index activity
           v
       ┌──────┐
       │ Solr │
       └──────┘
```

### Search Query Flow

```
┌──────────┐
│   User   │
│(Browser) │
└────┬─────┘
     │ 1. GET /search?q=yoga
     v
┌────────────┐
│  Frontend  │
└────┬───────┘
     │ 2. HTTP GET
     v
┌────────────┐
│ Search API │
└────┬───────┘
     │
     │ 3. Check cache
     v
┌────────────┐     Cache HIT ─────────┐
│ Memcached  │                         │
└────┬───────┘                         │
     │ Cache MISS                      │
     │                                 │
     │ 4. Query Solr                  │
     v                                 │
 ┌──────┐                              │
 │ Solr │                              │
 └───┬──┘                              │
     │ 5. Search results               │
     │                                 │
     v                                 │
┌────────────┐                         │
│ Search API │                         │
└────┬───────┘                         │
     │ 6. Store in cache               │
     │                                 │
     v                                 │
┌────────────┐                         │
│ Memcached  │                         │
└────────────┘                         │
     │                                 │
     │ 7. Return results               │
     └────────────────> ───────────────┘
                        │
                        v
                   ┌──────────┐
                   │  User    │
                   │(Browser) │
                   └──────────┘
```

## Authentication & Authorization Flow

### JWT-Based Authentication

```
┌──────────┐
│   User   │
└────┬─────┘
     │ 1. POST /auth/login
     │    { email, password }
     v
┌────────────┐
│ Users API  │
└────┬───────┘
     │ 2. Validate credentials
     │    (bcrypt password check)
     │
     │ 3. Generate JWT
     │    {
     │      user_id: 1,
     │      email: "user@example.com",
     │      role: "admin",  ← Important!
     │      exp: timestamp
     │    }
     v
┌────────────┐
│   User     │  Stores JWT in localStorage
└────┬───────┘  or memory
     │
     │ 4. Subsequent requests
     │    Authorization: Bearer <JWT>
     v
┌────────────┐
│ Backend    │  Validates JWT signature
│    API     │  Extracts role from claims
└────┬───────┘  Checks permissions
     │
     └──> Authorized ──> Process request
     │
     └──> Unauthorized ──> 401/403 Error
```

### Role-Based Access Control (RBAC)

**Role Hierarchy:**

```
┌─────────┐
│  root   │  All permissions (user deletion)
└────┬────┘
     │
┌────▼────────┐
│super_admin  │  Advanced system configuration
└────┬────────┘
     │
┌────▼────────┐
│   admin     │  User management, activity management
└────┬────────┘
     │
┌────▼────────┐
│ moderator   │  Content moderation
└────┬────────┘
     │
┌────▼────────┐
│    user     │  Basic operations
└─────────────┘
```

**Permission Check Example:**

```go
// In middleware
func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Validate JWT
        claims, err := ValidateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        
        // 2. Check role (from JWT claims, not database!)
        if claims.Role != "admin" && 
           claims.Role != "super_admin" && 
           claims.Role != "root" {
            c.JSON(403, gin.H{"error": "Admin role required"})
            c.Abort()
            return
        }
        
        // 3. Set context
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}
```

**Key Design Decision:** Role is stored in JWT claims, eliminating the need for database lookups on every request. This improves performance but requires token refresh when roles change.

## Event-Driven Architecture

### RabbitMQ Configuration

**Exchange:**
- Name: `activities_exchange`
- Type: topic
- Durable: true

**Queues:**
- `activities_search_queue` (bound to Search API)

**Routing Keys:**
- `activity.created`
- `activity.updated`
- `activity.deleted`

### Event Publishing (Activities API)

```go
// Publish event after activity creation
func PublishActivityCreated(activity Activity) error {
    event := Event{
        Type:      "activity.created",
        Payload:   activity,
        Timestamp: time.Now(),
    }
    
    body, _ := json.Marshal(event)
    
    return rabbitChannel.Publish(
        "activities_exchange",  // exchange
        "activity.created",      // routing key
        false,                   // mandatory
        false,                   // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
```

### Event Consumption (Search API)

```go
// Consumer goroutine
func ConsumeActivityEvents() {
    msgs, _ := rabbitChannel.Consume(
        "activities_search_queue",  // queue
        "",                          // consumer tag
        false,                       // auto-ack (manual ack)
        false,                       // exclusive
        false,                       // no-local
        false,                       // no-wait
        nil,                         // args
    )
    
    for msg := range msgs {
        var event Event
        json.Unmarshal(msg.Body, &event)
        
        switch event.Type {
        case "activity.created":
            solrClient.Index(event.Payload)
        case "activity.updated":
            solrClient.Update(event.Payload)
            cache.InvalidateRelated(event.Payload.ID)
        case "activity.deleted":
            solrClient.Delete(event.Payload.ID)
            cache.InvalidateRelated(event.Payload.ID)
        }
        
        msg.Ack(false)  // Manual acknowledgment
    }
}
```

## Caching Strategy

### Two-Level Caching (Search API)

**Level 1: Memcached (Application Cache)**

```
Cache Key Format: search:query:<hash_of_query_params>
TTL: 5 minutes (300 seconds)
Eviction Policy: LRU
```

**Level 2: Solr Internal Caches**

- **Filter Cache**: Caches filter query results
- **Query Result Cache**: Caches top N results
- **Document Cache**: Caches document fields

### Cache Invalidation Strategy

**Event-Based Invalidation:**

```go
// When activity is updated or deleted
func InvalidateCache(activityID string) {
    // Option 1: Invalidate all search caches (aggressive)
    cache.DeletePattern("search:query:*")
    
    // Option 2: Selective invalidation (better performance)
    // Invalidate only caches that might contain this activity
    // Implementation depends on tracking what's in cache
}
```

**TTL-Based Expiration:**
- All cache entries expire after 5 minutes
- Ensures eventually consistent data even if invalidation fails

### Cache Performance Metrics

Expected performance:
- **Cache Hit Rate**: 60-80% for common searches
- **Cache Hit Response Time**: 1-5ms
- **Cache Miss (Solr) Response Time**: 10-50ms
- **Cache Miss (Solr + Indexing) Response Time**: 50-200ms

## Database Design

### MySQL (Users API)

**Schema: users_db**

```sql
TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(100) UNIQUE,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    
    -- Profile fields
    first_name, last_name, avatar_url,
    bio TEXT, phone, birth_date, location, gender,
    
    -- Sports fields
    height, weight, sports_interests JSON,
    fitness_level ENUM,
    social_links JSON,
    
    -- System
    role VARCHAR(50) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN,
    
    -- Timestamps
    created_at, updated_at,
    
    INDEX (email),
    INDEX (username),
    INDEX (role)
)
```

**Design Decisions:**
- **Relational Model**: Users have well-defined structure
- **JSON Fields**: Flexible for social_links and sports_interests
- **Indexes**: On email, username for fast lookups
- **ENUM for Roles**: Type safety at database level

### MongoDB (Activities & Reservations)

**Database: activities_db**  
**Collection: activities**

```javascript
{
  _id: ObjectId,
  name: String,
  description: String,
  category: String,
  difficulty: String,
  location: String,
  price: Number,
  duration: Number,
  max_capacity: Number,
  instructor: String,
  schedule: [String],
  equipment: [String],
  image_url: String,
  is_active: Boolean,
  created_by: Number,  // User ID from Users API
  created_at: ISODate,
  updated_at: ISODate
}
```

**Database: reservations_db**  
**Collection: reservas**

```javascript
{
  _id: ObjectId,
  users_id: [Number],  // Array of user IDs
  cupo: Number,        // Number of spots
  actividad: String,   // Activity ObjectID as string
  date: ISODate,
  status: String,      // pendiente | confirmada | cancelada
  created_at: ISODate,
  updated_at: ISODate
}
```

**Design Decisions:**
- **Document Model**: Flexible schema for activities
- **Embedded Arrays**: schedule and equipment as arrays
- **ObjectID**: MongoDB native IDs for performance
- **Separate Databases**: Logical separation per service

### Apache Solr (Search Index)

**Core: activities**

```xml
<field name="id" type="string" />
<field name="name" type="text_general" />
<field name="description" type="text_general" />
<field name="category" type="string" />
<field name="difficulty" type="string" />
<field name="location" type="text_general" />
<field name="instructor" type="text_general" />
<field name="price" type="pfloat" />
<field name="is_active" type="boolean" />
```

**Design Decisions:**
- **Read-Optimized**: Denormalized for fast searches
- **Text Fields**: Full-text indexing on name, description, location
- **Facet Fields**: Category and difficulty for filtering
- **No Relationships**: Flat structure, no joins

## Deployment Architecture

### Docker Compose (Development)

```yaml
services:
  frontend:         # Port 3000
  users-api:        # Port 8081
  activities-api:   # Port 8082
  search-api:       # Port 8083
  reservations-api: # Port 8080
  mysql:            # Port 3307
  mongodb:          # Port 27017
  solr:             # Port 8983
  rabbitmq:         # Ports 5672, 15672
  memcached:        # Port 11211

networks:
  app-network:
    driver: bridge

volumes:
  mysql_data:
  mongo_data:
  solr_data:
  rabbitmq_data:
```

### Service Dependencies

```
┌─────────────┐
│   mysql     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  users-api  │
└──────┬──────┘
       │
       └──────────────────┐
                          │
┌─────────────┐           │
│   mongodb   │           │
└──────┬──────┘           │
       │                  │
       ├─────────────┐    │
       │             │    │
       ▼             ▼    │
┌─────────────┐ ┌──────────────┐
│activities-  │ │reservations- │
│    api      │ │     api      │
└──────┬──────┘ └──────────────┘
       │
       │ publishes events
       ▼
┌─────────────┐
│  rabbitmq   │
└──────┬──────┘
       │ consumes events
       ▼
┌─────────────┐      ┌─────────────┐
│  search-api │─────>│    solr     │
└──────┬──────┘      └─────────────┘
       │
       └─────────────>┌─────────────┐
                      │  memcached  │
                      └─────────────┘
```

### Health Check Strategy

Each service implements health check endpoints:

```
GET /health   or   GET /healthz
```

Docker Compose health checks:

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8081/api/v1/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

## Security Considerations

### 1. Authentication Security

- **Password Storage**: bcrypt with cost factor 10
- **JWT Signing**: HS256 with secret key (should use RS256 with key rotation in production)
- **Token Expiration**: 24 hours
- **Token Refresh**: Explicit refresh endpoint

### 2. Authorization Security

- **Role-Based Access**: Middleware enforces role requirements
- **JWT Claims**: Role stored in token for fast checks
- **Principle of Least Privilege**: Users get minimum necessary permissions

### 3. Network Security

- **Internal Network**: Services communicate via Docker network
- **Port Exposure**: Only frontend and APIs exposed to host
- **Database Isolation**: Databases not exposed to host in production

### 4. Data Security

- **SQL Injection**: Protected by GORM parameterized queries
- **NoSQL Injection**: Protected by MongoDB driver
- **XSS**: React automatically escapes output
- **CORS**: Configurable allowed origins

### 5. Input Validation

- **Backend Validation**: Gin binding validators on all inputs
- **Frontend Validation**: React form validation
- **File Upload**: Size and type restrictions on avatar uploads

## Scalability & Performance

### Horizontal Scaling

Each service can be scaled independently:

```yaml
services:
  users-api:
    deploy:
      replicas: 3  # Run 3 instances
```

**Stateless Services**: All backend services are stateless, enabling easy horizontal scaling.

### Vertical Scaling

Resource limits can be adjusted per service:

```yaml
services:
  search-api:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

### Database Scaling

**MySQL:**
- **Read Replicas**: For read-heavy operations
- **Connection Pooling**: GORM manages connection pools

**MongoDB:**
- **Replica Sets**: For high availability
- **Sharding**: For horizontal partitioning of large datasets

**Solr:**
- **SolrCloud**: Distributed search with leader-follower architecture
- **Sharding**: Partition index across multiple nodes

### Caching Benefits

- **Reduced Database Load**: Memcached absorbs repeated queries
- **Improved Response Time**: 10-50x faster than database queries
- **Increased Throughput**: More requests per second

### Performance Targets

| Metric | Target | Actual (Development) |
|--------|--------|----------------------|
| API Response Time (p95) | < 200ms | ~100ms |
| Search Query (cached) | < 10ms | ~5ms |
| Search Query (uncached) | < 100ms | ~50ms |
| JWT Validation | < 5ms | ~2ms |
| Event Processing Latency | < 1s | ~500ms |

## Future Enhancements

### 1. Service Mesh (Istio)
- Service discovery
- Load balancing
- Circuit breaking
- Distributed tracing

### 2. API Gateway (Kong/Traefik)
- Centralized routing
- Rate limiting
- Request authentication
- API versioning

### 3. Observability
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Metrics**: Prometheus + Grafana
- **Tracing**: Jaeger for distributed tracing

### 4. CI/CD Pipeline
- Automated testing
- Docker image building
- Kubernetes deployment
- Blue-green deployments

### 5. Additional Services
- **Notifications Service**: Email/SMS notifications
- **Analytics Service**: Usage analytics and reporting
- **Payment Service**: Payment processing integration
- **Recommendations Service**: ML-based activity recommendations

## Related Documentation

- [Users API](./api/users-api.md)
- [Activities API](./api/activities-api.md)
- [Search API](./api/search-api.md)
- [Reservations API](./api/reservations-api.md)
- [Setup Guide](../SETUP.md)
- [Testing Guide](../TESTING.md)

