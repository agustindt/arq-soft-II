# Search API Documentation

## Overview

The Search API is a Go-based microservice that provides full-text search capabilities for activities using Apache Solr. It implements a caching layer with Memcached and consumes events from RabbitMQ to keep the search index synchronized.

**Technology Stack:**
- Language: Go 1.21+
- Framework: Gin
- Search Engine: Apache Solr 9.4
- Cache: Memcached
- Message Queue: RabbitMQ (event consumer)

**Port:** 8083

**Base URL:** `http://localhost:8083`

## Architecture

### CQRS Pattern

The Search API implements the **Command Query Responsibility Segregation (CQRS)** pattern:

- **Write Side**: Activities API handles all write operations (create, update, delete)
- **Read Side**: Search API handles all search/query operations

This separation provides:
- **Performance**: Optimized search queries without affecting write operations
- **Scalability**: Search and write operations can scale independently
- **Flexibility**: Different data models for searching vs. storing

### System Components

```
┌─────────────────┐      ┌──────────────┐      ┌──────────────┐
│ Activities API  │      │   RabbitMQ   │      │  Search API  │
│  (Write Side)   │─────>│   Exchange   │─────>│  (Read Side) │
└─────────────────┘      └──────────────┘      └───────┬──────┘
                                                        │
                                                        v
                                    ┌──────────────────────────────┐
                                    │                              │
                                    │   ┌───────────┐              │
                                    │   │   Solr    │              │
                                    │   │  (Index)  │              │
                                    │   └───────────┘              │
                                    │                              │
                                    │   ┌───────────┐              │
                                    │   │ Memcached │              │
                                    │   │  (Cache)  │              │
                                    │   └───────────┘              │
                                    │                              │
                                    └──────────────────────────────┘
```

### Apache Solr Integration

#### Solr Core Configuration

- **Core Name**: `activities`
- **Configset**: Custom configuration with activity-specific schema
- **Location**: `/opt/solr/server/solr/activities`

#### Solr Schema (schema.xml)

```xml
<schema name="activities" version="1.6">
  <field name="id" type="string" indexed="true" stored="true" required="true" />
  <field name="name" type="text_general" indexed="true" stored="true" />
  <field name="description" type="text_general" indexed="true" stored="true" />
  <field name="category" type="string" indexed="true" stored="true" />
  <field name="difficulty" type="string" indexed="true" stored="true" />
  <field name="location" type="text_general" indexed="true" stored="true" />
  <field name="instructor" type="text_general" indexed="true" stored="true" />
  <field name="price" type="pfloat" indexed="true" stored="true" />
  <field name="duration" type="pint" indexed="true" stored="true" />
  <field name="max_capacity" type="pint" indexed="true" stored="true" />
  <field name="is_active" type="boolean" indexed="true" stored="true" />
  <field name="created_at" type="pdate" indexed="true" stored="true" />
  
  <!-- Default search field -->
  <field name="_text_" type="text_general" indexed="true" stored="false" multiValued="true"/>
  
  <!-- Copy fields for comprehensive search -->
  <copyField source="name" dest="_text_"/>
  <copyField source="description" dest="_text_"/>
  <copyField source="location" dest="_text_"/>
  <copyField source="instructor" dest="_text_"/>
</schema>
```

#### Indexed Fields

| Field | Type | Indexed | Stored | Description |
|-------|------|---------|--------|-------------|
| `id` | string | ✅ | ✅ | MongoDB ObjectID |
| `name` | text_general | ✅ | ✅ | Activity name (searchable) |
| `description` | text_general | ✅ | ✅ | Full description (searchable) |
| `category` | string | ✅ | ✅ | Category (exact match) |
| `difficulty` | string | ✅ | ✅ | Difficulty level (exact match) |
| `location` | text_general | ✅ | ✅ | Location (searchable) |
| `instructor` | text_general | ✅ | ✅ | Instructor name (searchable) |
| `price` | float | ✅ | ✅ | Price (range queries) |
| `duration` | int | ✅ | ✅ | Duration in minutes |
| `max_capacity` | int | ✅ | ✅ | Maximum capacity |
| `is_active` | boolean | ✅ | ✅ | Active status |
| `created_at` | date | ✅ | ✅ | Creation timestamp |

### Caching Strategy

The Search API implements a two-level caching strategy:

#### Level 1: Query Result Caching (Memcached)

- **TTL**: 5 minutes (300 seconds)
- **Key Format**: `search:query:<hash_of_query_params>`
- **Purpose**: Cache complete search results for identical queries
- **Invalidation**: Time-based (TTL) and event-based (activity updates)

#### Level 2: Solr Internal Cache

Solr maintains its own internal caches:
- Filter cache
- Query result cache
- Document cache

#### Cache Flow

```
┌────────┐
│ Client │
└───┬────┘
    │ GET /search?q=yoga
    v
┌───────────────┐
│  Search API   │
└───────┬───────┘
        │
        │ 1. Check Memcached
        v
    ┌────────────┐     Cache HIT
    │ Memcached  │────────────> Return cached result
    └────────────┘              (Skip Solr query)
        │
        │ Cache MISS
        v
    ┌────────────┐
    │    Solr    │
    └─────┬──────┘
          │
          │ 2. Query Solr
          │ 3. Store in Memcached
          │ 4. Return result
          v
    ┌────────────┐
    │   Client   │
    └────────────┘
```

## API Endpoints

### Search Endpoints

#### GET /search

Search for activities with optional filters.

**Query Parameters:**

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `q` | string | Yes | Search query | `q=yoga` |
| `category` | string | No | Filter by category | `category=yoga` |
| `difficulty` | string | No | Filter by difficulty | `difficulty=beginner` |
| `minPrice` | float | No | Minimum price | `minPrice=0` |
| `maxPrice` | float | No | Maximum price | `maxPrice=5000` |
| `location` | string | No | Filter by location | `location=Palermo` |
| `rows` | int | No | Results per page (default: 10) | `rows=20` |
| `start` | int | No | Offset for pagination (default: 0) | `start=10` |

**Example Requests:**

```bash
# Simple search
GET /search?q=yoga

# Search with category filter
GET /search?q=fitness&category=gym

# Search with price range
GET /search?q=tennis&minPrice=1000&maxPrice=3000

# Search with multiple filters
GET /search?q=football&difficulty=intermediate&location=Buenos Aires

# Paginated search
GET /search?q=running&rows=20&start=20
```

**Response (200 OK):**

```json
{
  "query": "yoga",
  "results": [
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
      "is_active": true,
      "created_at": "2025-11-14T10:00:00Z"
    },
    {
      "id": "673590beeb2a7e80f9ff0c64",
      "name": "Evening Yoga Flow",
      "description": "Dynamic yoga flow for intermediate practitioners",
      "category": "yoga",
      "difficulty": "intermediate",
      "location": "Recoleta Studio",
      "price": 2000.00,
      "duration": 75,
      "max_capacity": 12,
      "instructor": "Laura Martinez",
      "is_active": true,
      "created_at": "2025-11-14T11:00:00Z"
    }
  ],
  "total": 2,
  "page": 0,
  "rows": 10,
  "cached": false
}
```

**Response Fields:**

- `query` - Original search query
- `results` - Array of matching activities
- `total` - Total number of matching results
- `page` - Current page (calculated from `start` / `rows`)
- `rows` - Number of results per page
- `cached` - Whether result was served from cache

**Error Response (400 Bad Request):**

```json
{
  "error": "Query parameter 'q' is required"
}
```

---

### Health Check

#### GET /health

Check if the service and its dependencies are healthy.

**Response (200 OK):**

```json
{
  "status": "healthy",
  "service": "search-api",
  "solr": "connected",
  "memcached": "connected",
  "rabbitmq": "connected"
}
```

**Response (503 Service Unavailable):**

```json
{
  "status": "unhealthy",
  "service": "search-api",
  "solr": "disconnected",
  "memcached": "connected",
  "rabbitmq": "connected"
}
```

## Event-Driven Indexing

### RabbitMQ Consumer

The Search API consumes events from the `activities_exchange` to keep the Solr index synchronized.

#### Consumed Events

##### 1. activity.created

**Routing Key:** `activity.created`

**Payload:**
```json
{
  "id": "673590beeb2a7e80f9ff0c63",
  "name": "New Activity",
  "description": "Description",
  "category": "yoga",
  "difficulty": "beginner",
  "location": "Buenos Aires",
  "price": 1500.00,
  "duration": 90,
  "max_capacity": 15,
  "instructor": "John Doe",
  "is_active": true,
  "created_at": "2025-11-14T10:00:00Z"
}
```

**Action:** Index the new activity in Solr

---

##### 2. activity.updated

**Routing Key:** `activity.updated`

**Payload:**
```json
{
  "id": "673590beeb2a7e80f9ff0c63",
  "name": "Updated Activity Name",
  "description": "Updated description",
  ...
}
```

**Action:** 
1. Re-index the updated activity in Solr
2. Invalidate related cache entries

---

##### 3. activity.deleted

**Routing Key:** `activity.deleted`

**Payload:**
```json
{
  "id": "673590beeb2a7e80f9ff0c63"
}
```

**Action:**
1. Remove the activity from Solr index
2. Invalidate related cache entries

---

### Event Processing Flow

```
┌─────────────────┐
│ Activities API  │
│ publishes event │
└────────┬────────┘
         │
         │ activity.created
         │ activity.updated
         │ activity.deleted
         v
    ┌─────────┐
    │ RabbitMQ│
    │Exchange │
    └────┬────┘
         │
         │ Route to queue
         v
    ┌─────────┐
    │  Queue  │
    │activities│
    │_search  │
    └────┬────┘
         │
         │ Consume
         v
┌────────────────┐
│   Search API   │
│Consumer Service│
└────────┬───────┘
         │
         ├─> Index/Update/Delete in Solr
         │
         └─> Invalidate Memcached entries
```

### Consumer Configuration

- **Exchange**: `activities_exchange`
- **Queue**: `activities_search_queue`
- **Routing Keys**: `activity.created`, `activity.updated`, `activity.deleted`
- **Exchange Type**: topic
- **Durable**: Yes
- **Auto-Acknowledge**: No (manual acknowledgment after processing)

## Search Query Syntax

### Basic Search

Simple text search across all indexed text fields:

```
GET /search?q=yoga
```

Searches in: name, description, location, instructor

### Category Filter

Exact match on category:

```
GET /search?q=fitness&category=gym
```

### Difficulty Filter

Exact match on difficulty level:

```
GET /search?q=training&difficulty=advanced
```

### Price Range

Filter activities by price range:

```
GET /search?q=sports&minPrice=1000&maxPrice=3000
```

### Location Search

Search in location field:

```
GET /search?q=football&location=Palermo
```

### Pagination

Navigate through results:

```
GET /search?q=running&rows=10&start=0   # First page
GET /search?q=running&rows=10&start=10  # Second page
GET /search?q=running&rows=10&start=20  # Third page
```

### Combined Filters

Use multiple filters together:

```
GET /search?q=tennis&category=tennis&difficulty=intermediate&minPrice=2000&maxPrice=5000&location=Buenos%20Aires
```

## Cache Invalidation

### Automatic Invalidation

The cache is automatically invalidated when:

1. **Activity Created**: New activities don't affect existing cache
2. **Activity Updated**: Invalidates cache entries that might include this activity
3. **Activity Deleted**: Invalidates cache entries that included this activity

### Manual Cache Clear

For development/debugging, you can clear specific cache entries by connecting to Memcached:

```bash
# Connect to Memcached
docker-compose exec memcached nc localhost 11211

# Flush all cache
flush_all

# Or delete specific key
delete search:query:<hash>
```

### Cache Key Generation

Cache keys are generated based on:
- Search query (`q`)
- All filter parameters (category, difficulty, price range, location)
- Pagination parameters (rows, start)

Example cache key:
```
search:query:sha256(q=yoga&category=yoga&difficulty=beginner&rows=10&start=0)
```

## Performance Considerations

### Indexing Performance

- **Batch Indexing**: When reindexing, activities are indexed in batches
- **Async Indexing**: Event consumption is asynchronous
- **Commit Strategy**: Soft commits every 5 seconds, hard commits every 60 seconds

### Search Performance

- **Caching**: Memcached reduces Solr load for repeated queries
- **Cache Hit Rate**: Typically 60-80% for common searches
- **Solr Response Time**: ~10-50ms for indexed queries
- **Cached Response Time**: ~1-5ms for cache hits

### Scalability

- **Horizontal Scaling**: Multiple Search API instances can share the same Solr and Memcached
- **Solr Cloud**: For high availability, Solr can be deployed in SolrCloud mode
- **Sharding**: Large datasets can be sharded across multiple Solr cores

## Error Handling

### Common HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful search query |
| 400 | Bad Request | Missing query parameter or invalid filters |
| 500 | Internal Server Error | Solr or Memcached error |
| 503 | Service Unavailable | Solr is down or unreachable |

### Error Format

```json
{
  "error": "Error message description"
}
```

## Environment Variables

```bash
# Server
PORT=8083

# Apache Solr
SOLR_URL=http://solr:8983/solr
SOLR_CORE=activities

# Memcached
MEMCACHED_HOST=memcached:11211

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin123@rabbitmq:5672/
RABBITMQ_EXCHANGE=activities_exchange
RABBITMQ_QUEUE=activities_search_queue

# Cache
CACHE_TTL=300  # 5 minutes in seconds

# CORS
CORS_ALLOWED_ORIGINS=*
```

## Service Dependencies

### Required Services

1. **Apache Solr** - Search engine
   - URL: `http://solr:8983/solr`
   - Core: `activities`

2. **Memcached** - Cache layer
   - Host: `memcached:11211`

3. **RabbitMQ** - Message broker
   - URL: `amqp://rabbitmq:5672/`
   - Exchange: `activities_exchange`

### Service Startup Order

1. Solr must be running with `activities` core created
2. Memcached must be running
3. RabbitMQ must be running with exchange configured
4. Then Search API can start

## Development

### Running Locally

```bash
cd backend/search-api

# Install dependencies
go mod tidy

# Run the service
go run main.go
```

### Testing Search

```bash
# Test basic search
curl "http://localhost:8083/search?q=yoga"

# Test with filters
curl "http://localhost:8083/search?q=fitness&category=gym&difficulty=beginner"

# Test health endpoint
curl "http://localhost:8083/health"
```

### Reindexing All Activities

To rebuild the Solr index from scratch:

```bash
# From project root
./scripts/reindex-activities.sh
```

This script:
1. Fetches all activities from Activities API
2. Deletes existing Solr index
3. Reindexes all activities
4. Clears Memcached

### Monitoring

#### Solr Admin UI

Access Solr's admin interface:
```
http://localhost:8983/solr/#/activities
```

Features:
- Query testing
- Schema browser
- Core statistics
- Cache statistics

#### Cache Statistics

Check Memcached statistics:

```bash
docker-compose exec memcached nc localhost 11211
stats
```

Important metrics:
- `curr_items`: Current items in cache
- `get_hits`: Cache hits
- `get_misses`: Cache misses
- `bytes_used`: Memory usage

## Common Operations

### Clear All Cache

```bash
docker-compose exec memcached nc localhost 11211
flush_all
quit
```

### Delete Solr Index

```bash
curl "http://localhost:8983/solr/activities/update?commit=true" \
  -H "Content-Type: application/json" \
  -d '{"delete":{"query":"*:*"}}'
```

### Check Indexed Document Count

```bash
curl "http://localhost:8983/solr/activities/select?q=*:*&rows=0"
```

Look for `numFound` in the response.

## Example Usage

### Search Yoga Classes

```bash
curl "http://localhost:8083/search?q=yoga&rows=5"
```

### Find Beginner Activities

```bash
curl "http://localhost:8083/search?q=*&difficulty=beginner&rows=10"
```

### Search by Location

```bash
curl "http://localhost:8083/search?q=*&location=Palermo"
```

### Advanced Search with All Filters

```bash
curl "http://localhost:8083/search?q=fitness&category=gym&difficulty=intermediate&minPrice=1000&maxPrice=3000&rows=20&start=0"
```

## Related Documentation

- [Users API](./users-api.md)
- [Activities API](./activities-api.md)
- [Reservations API](./reservations-api.md)
- [Architecture Overview](../architecture.md)

