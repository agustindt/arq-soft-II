# 🧪 Testing Documentation

This document provides comprehensive information about testing the Sports Activities Platform.

## Table of Contents

- [Overview](#overview)
- [Test Scripts](#test-scripts)
- [Running Tests](#running-tests)
- [Test Details](#test-details)
- [Troubleshooting](#troubleshooting)
- [Expected Outputs](#expected-outputs)

## Overview

The platform includes three main test scripts that verify different aspects of the system:

1. **test-infrastructure.sh** - Tests infrastructure components (databases, message queues, search engine)
2. **test-backend.sh** - Tests all backend API services and their integration
3. **test-all.sh** - Master script that runs all tests in sequence

## Test Scripts

### test-infrastructure.sh

Tests all infrastructure components to ensure they are properly configured and accessible.

**What it tests:**
- MySQL database connectivity and schema
- MongoDB connectivity and databases
- RabbitMQ message broker and queues
- Apache Solr search engine and core
- Memcached cache service

**Usage:**
```bash
./scripts/test-infrastructure.sh
# or
make test-infra
```

**Expected Duration:** 10-30 seconds

### test-backend.sh

Comprehensive backend testing including API endpoints, authentication, and integration flows.

**What it tests:**
- Health checks for all APIs
- User registration and authentication
- JWT token generation and validation
- Activity CRUD operations
- Search API functionality
- Reservation creation and management
- End-to-end flow: Create Activity → Event → Index → Search
- RabbitMQ message processing
- Memcached caching

**Usage:**
```bash
./scripts/test-backend.sh
# or
make test-backend
```

**Expected Duration:** 30-60 seconds

**Note:** This script creates test data (users, activities, reservations) that will persist in your databases.

### test-all.sh

Master script that orchestrates all test suites and provides a comprehensive report.

**What it does:**
1. Verifies Docker is running
2. Checks that all services are up
3. Runs infrastructure tests
4. Runs backend tests
5. Tests frontend accessibility
6. Performs end-to-end connectivity tests
7. Generates a final report

**Usage:**
```bash
./scripts/test-all.sh
# or
make test
```

**Expected Duration:** 1-2 minutes

## Running Tests

### Prerequisites

Before running tests, ensure:

1. **Docker is running**
   ```bash
   docker info
   ```

2. **All services are started**
   ```bash
   make start
   # or
   ./scripts/start-all.sh
   ```

3. **Services are healthy**
   ```bash
   docker-compose ps
   ```
   All services should show "Up" status.

### Quick Start

```bash
# Start the entire stack
make start

# Wait for services to be ready (about 30-60 seconds)
# Then run all tests
make test
```

### Running Individual Test Suites

```bash
# Test only infrastructure
make test-infra

# Test only backend APIs
make test-backend

# Test everything
make test
```

## Test Details

### Infrastructure Tests

#### MySQL Tests
- ✅ Container is running
- ✅ Port 3307 is accessible
- ✅ Database accepts connections
- ✅ Database `users_db` exists
- ✅ Table `users` exists

#### MongoDB Tests
- ✅ Container is running
- ✅ Port 27017 is accessible
- ✅ Database accepts connections
- ✅ Database `activitiesdb` is accessible
- ✅ Database `reservasdb` is accessible

#### RabbitMQ Tests
- ✅ Container is running
- ✅ Port 5672 (AMQP) is accessible
- ✅ Port 15672 (Management UI) is accessible
- ✅ Health check passes
- ✅ Management UI responds
- ✅ Exchange `entity.events` exists
- ✅ Queue `search-sync` exists

#### Solr Tests
- ✅ Container is running
- ✅ Port 8983 is accessible
- ✅ Admin UI responds
- ✅ Core `activities` exists
- ✅ Core ping is OK
- ✅ Schema is loaded

#### Memcached Tests
- ✅ Container is running
- ✅ Port 11211 is accessible
- ✅ Accepts set/get commands

### Backend Tests

#### Health Checks
- ✅ Users API health endpoint
- ✅ Activities API health endpoint
- ✅ Search API health endpoint
- ✅ Reservations API health endpoint

#### Users API Tests
- ✅ User registration
- ✅ User login and JWT token generation
- ✅ Get user profile (authenticated)
- ✅ Admin user registration
- ✅ Admin login

#### Activities API Tests
- ✅ Create activity (admin only)
- ✅ List all activities
- ✅ Get activity by ID
- ✅ Filter activities by category

#### Search API Tests
- ✅ Search activities by text
- ✅ Search with category filter
- ✅ Search with difficulty filter
- ✅ Cache functionality
- ✅ Solr indexing verification

#### Reservations API Tests
- ✅ Create reservation (authenticated)
- ✅ List reservations
- ✅ Get reservation by ID

#### Integration Tests
- ✅ RabbitMQ message processing
- ✅ Memcached statistics
- ✅ Complete flow: Create → Event → Index → Search

### Frontend Tests

#### Basic Checks
- ✅ Frontend is accessible on port 3000
- ✅ Returns valid HTML
- ✅ React application is loaded

#### End-to-End Tests
- ✅ Frontend can communicate with Users API
- ✅ Frontend can communicate with Activities API
- ✅ Frontend can communicate with Search API
- ✅ Frontend can communicate with Reservations API

## Troubleshooting

### Common Issues

#### 1. Tests Fail: "Docker no está corriendo"

**Problem:** Docker Desktop is not running.

**Solution:**
```bash
# Start Docker Desktop, then verify
docker info
```

#### 2. Tests Fail: "Servicios no están corriendo"

**Problem:** Services haven't been started yet.

**Solution:**
```bash
# Start all services
make start

# Wait for services to be ready
# Check status
docker-compose ps

# Then run tests again
make test
```

#### 3. Infrastructure Tests Fail: "MySQL no acepta conexiones"

**Problem:** MySQL container is still initializing.

**Solution:**
```bash
# Wait a bit longer, MySQL can take 30-60 seconds to be ready
# Check MySQL logs
docker-compose logs mysql

# Restart MySQL if needed
docker-compose restart mysql
```

#### 4. Backend Tests Fail: "Token JWT no obtenido"

**Problem:** Users API might not be ready or there's an authentication issue.

**Solution:**
```bash
# Check Users API logs
docker-compose logs users-api

# Verify Users API is responding
curl http://localhost:8081/health

# Restart Users API if needed
docker-compose restart users-api
```

#### 5. Search Tests Fail: "No hay documentos en Solr"

**Problem:** Solr indexing might be delayed or the consumer isn't processing messages.

**Solution:**
```bash
# Check Search API logs
docker-compose logs search-api

# Check RabbitMQ queues
# Visit http://localhost:15672 (admin/admin123)
# Check if messages are in the queue

# Restart Search API to trigger re-indexing
docker-compose restart search-api

# Wait a few seconds for indexing
sleep 10

# Run tests again
make test-backend
```

#### 6. Tests Timeout

**Problem:** Services are taking too long to respond.

**Solution:**
```bash
# Check resource usage
docker stats

# Restart services
make restart

# Wait longer before testing
# Some services need more time on slower machines
```

### Debugging Commands

```bash
# View logs from all services
make logs

# View logs from specific service
docker-compose logs -f users-api
docker-compose logs -f activities-api
docker-compose logs -f search-api

# Check service status
docker-compose ps

# Check service health
curl http://localhost:8081/health
curl http://localhost:8082/healthz
curl http://localhost:8083/health
curl http://localhost:8080/health

# Check RabbitMQ management UI
# Open http://localhost:15672 (admin/admin123)

# Check Solr admin UI
# Open http://localhost:8983/solr

# Connect to databases
make db-mysql
make db-mongo
```

### Resetting Test Environment

If tests are failing due to corrupted state:

```bash
# Stop everything and remove volumes
make clean

# Start fresh
make start

# Wait for services to be ready
sleep 60

# Run tests
make test
```

## Expected Outputs

### test-infrastructure.sh Output

```
========================================
🧪 Testing Infraestructura
========================================

🗄️  MySQL Tests
  ✓ MySQL contenedor corriendo
  ✓ MySQL puerto 3307 accesible
  ✓ MySQL acepta conexiones
  ✓ Base de datos 'users_db' existe
  ✓ Tabla 'users' existe en users_db

🍃 MongoDB Tests
  ✓ MongoDB contenedor corriendo
  ✓ MongoDB puerto 27017 accesible
  ✓ MongoDB acepta conexiones
  ✓ Database 'activitiesdb' accesible
  ✓ Database 'reservasdb' accesible

[... more tests ...]

========================================
📊 Resumen de Tests
========================================

  Total de tests: 25
  Tests pasados:  25
  Tests fallidos: 0

✅ Todos los tests de infraestructura pasaron!
```

### test-backend.sh Output

```
========================================
🧪 Testing Backend Services
========================================

🏥 Health Checks
  ✓ Users API health check (HTTP 200)
  ✓ Activities API health check (HTTP 200)
  ✓ Search API health check (HTTP 200)
  ✓ Reservations API health check (HTTP 200)

👤 Users API Tests
  → Registrando usuario de prueba...
  ✓ Registrar usuario (HTTP 201)
  → Iniciando sesión...
  ✓ Login usuario (HTTP 200)
  ✓ Token JWT obtenido
  → Obteniendo perfil del usuario...
  ✓ Obtener perfil autenticado (HTTP 200)

[... more tests ...]

========================================
📊 Resumen de Tests
========================================

  Total de tests: 45
  Tests pasados:  45
  Tests fallidos: 0

✅ Todos los tests de backend pasaron!
```

### test-all.sh Output

```
========================================
🧪 TEST SUITE COMPLETO
Testing Backend + Frontend + Infra
========================================

🔍 Verificando que los servicios están corriendo...

========================================
📋 FASE 1: Infraestructura
========================================
[... infrastructure test output ...]

========================================
📋 FASE 2: Backend Services
========================================
[... backend test output ...]

========================================
📋 FASE 3: Frontend
========================================
  ✓ Frontend accesible en http://localhost:3000
  ✓ Frontend devuelve HTML válido
  ✓ Aplicación React cargada

========================================
📋 FASE 4: End-to-End Flow
========================================
  ✓ Frontend ↔ users-api comunicación OK
  ✓ Frontend ↔ activities-api comunicación OK
  ✓ Frontend ↔ search-api comunicación OK
  ✓ Frontend ↔ reservations-service comunicación OK

========================================
📊 REPORTE FINAL
========================================

  Duración total: 95s

  Resultados por fase:
    ✓ Infraestructura: PASSED
    ✓ Backend Services: PASSED
    ✓ Frontend: PASSED
    ✓ End-to-End: PASSED

========================================
✅ TODOS LOS TESTS PASARON!
El sistema está funcionando OK
========================================
```

## Best Practices

1. **Run tests after starting services:** Always wait for services to be fully ready before running tests.

2. **Clean state for consistent results:** If you're experiencing flaky tests, try cleaning and restarting:
   ```bash
   make clean
   make start
   make test
   ```

3. **Check logs when tests fail:** Always check service logs when tests fail to understand the root cause.

4. **Run individual test suites during development:** Use `make test-infra` or `make test-backend` for faster feedback during development.

5. **Verify seed data:** After running `make seed`, verify the data was loaded correctly:
   ```bash
   make db-mysql
   # Then: SELECT * FROM users WHERE email LIKE '%@test.com';
   
   make db-mongo
   # Then: db.activities.countDocuments()
   ```

## Continuous Integration

These test scripts are designed to be run in CI/CD pipelines. Example GitHub Actions workflow:

```yaml
name: Test Suite

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Start services
        run: make start
      - name: Wait for services
        run: sleep 60
      - name: Run tests
        run: make test
```

---

For more information, see the main [README.md](./README.md) file.

