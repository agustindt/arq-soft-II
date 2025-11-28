#!/bin/bash

# Script para testear todos los servicios backend
# Uso: ./scripts/test-backend.sh

set -e

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0
TOKEN=""
ADMIN_TOKEN=""
ACTIVITY_ID=""

# Función para reportar test exitoso
pass() {
    echo -e "  ${GREEN}✓${NC} $1"
    ((PASSED++))
}

# Función para reportar test fallido
fail() {
    echo -e "  ${RED}✗${NC} $1"
    echo -e "    ${CYAN}Response: $2${NC}" | head -n 3
    ((FAILED++))
}

# Función para ejecutar test HTTP
test_http() {
    local description="$1"
    local method="$2"
    local url="$3"
    local expected_code="$4"
    local data="${5:-}"
    local headers="${6:-}"

    local response
    local http_code

    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            ${headers:+-H "$headers"} \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            ${headers:+-H "$headers"})
    fi

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "$expected_code" ]; then
        pass "$description (HTTP $http_code)"
        echo "$body"
    else
        fail "$description (Expected $expected_code, got $http_code)" "$body"
        echo ""
    fi
}

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🧪 Testing Backend Services${NC}"
echo -e "${BLUE}========================================${NC}\n"

# ===== Health Checks =====
echo -e "${YELLOW}🏥 Health Checks${NC}"

test_http "Users API health check" "GET" "http://localhost:8081/api/v1/health" "200"
test_http "Activities API health check" "GET" "http://localhost:8082/healthz" "200"
test_http "Search API health check" "GET" "http://localhost:8083/health" "200"
test_http "Reservations API health check" "GET" "http://localhost:8080/healthz" "200"

echo ""

# ===== Users API Tests =====
echo -e "${YELLOW}👤 Users API Tests${NC}"

# Registrar usuario
echo -e "\n  ${CYAN}→ Registrando usuario de prueba...${NC}"
REGISTER_RESPONSE=$(test_http "Registrar usuario" "POST" "http://localhost:8081/users" "201" \
    '{
        "username": "testuser",
        "email": "test@example.com",
        "password": "Test123!",
        "role": "user"
    }')

# Login
echo -e "\n  ${CYAN}→ Iniciando sesión...${NC}"
LOGIN_RESPONSE=$(test_http "Login usuario" "POST" "http://localhost:8081/auth/login" "200" \
    '{
        "email": "test@example.com",
        "password": "Test123!"
    }')

# Extraer token
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || echo "")

if [ -n "$TOKEN" ]; then
    pass "Token JWT obtenido"
    echo -e "    ${CYAN}Token: ${TOKEN:0:50}...${NC}"
else
    fail "No se pudo obtener token JWT" "$LOGIN_RESPONSE"
fi

# Obtener perfil con token
echo -e "\n  ${CYAN}→ Obteniendo perfil del usuario...${NC}"
if [ -n "$TOKEN" ]; then
    test_http "Obtener perfil autenticado" "GET" "http://localhost:8081/users/profile" "200" "" "Authorization: Bearer $TOKEN"
else
    fail "No se puede obtener perfil sin token" ""
fi

# Registrar admin
echo -e "\n  ${CYAN}→ Registrando usuario admin...${NC}"
ADMIN_REGISTER=$(test_http "Registrar admin" "POST" "http://localhost:8081/users" "201" \
    '{
        "username": "adminuser",
        "email": "admin@example.com",
        "password": "Admin123!",
        "role": "admin"
    }')

# Login admin
echo -e "\n  ${CYAN}→ Login como admin...${NC}"
ADMIN_LOGIN=$(test_http "Login admin" "POST" "http://localhost:8081/auth/login" "200" \
    '{
        "email": "admin@example.com",
        "password": "Admin123!"
    }')

ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || echo "")

if [ -n "$ADMIN_TOKEN" ]; then
    pass "Admin token obtenido"
else
    fail "No se pudo obtener admin token" "$ADMIN_LOGIN"
fi

echo ""

# ===== Activities API Tests =====
echo -e "${YELLOW}🏃 Activities API Tests${NC}"

# Crear actividad (como admin)
echo -e "\n  ${CYAN}→ Creando actividad deportiva...${NC}"
if [ -n "$ADMIN_TOKEN" ]; then
    CREATE_ACTIVITY=$(test_http "Crear actividad" "POST" "http://localhost:8082/activities" "201" \
        '{
            "name": "Fútbol 5 en Parque",
            "description": "Partido amistoso de fútbol 5 en el parque central",
            "category": "football",
            "difficulty": "intermediate",
            "duration": 90,
            "maxParticipants": 10,
            "location": "Parque Central",
            "price": 500.00
        }' "Authorization: Bearer $ADMIN_TOKEN")

    ACTIVITY_ID=$(echo "$CREATE_ACTIVITY" | grep -o '"id":"[^"]*' | cut -d'"' -f4 || echo "")

    if [ -n "$ACTIVITY_ID" ]; then
        pass "Actividad creada con ID: $ACTIVITY_ID"
    else
        fail "No se pudo obtener ID de actividad" "$CREATE_ACTIVITY"
    fi
else
    fail "No se puede crear actividad sin admin token" ""
fi

# Esperar a que RabbitMQ propague el evento
echo -e "\n  ${CYAN}→ Esperando propagación de evento (3s)...${NC}"
sleep 3

# Obtener todas las actividades
echo -e "\n  ${CYAN}→ Listando todas las actividades...${NC}"
test_http "Listar actividades" "GET" "http://localhost:8082/activities" "200"

# Obtener actividad por ID
if [ -n "$ACTIVITY_ID" ]; then
    echo -e "\n  ${CYAN}→ Obteniendo actividad por ID...${NC}"
    test_http "Obtener actividad por ID" "GET" "http://localhost:8082/activities/$ACTIVITY_ID" "200"
fi

# Filtrar por categoría
echo -e "\n  ${CYAN}→ Filtrando por categoría...${NC}"
test_http "Filtrar por categoría football" "GET" "http://localhost:8082/activities?category=football" "200"

echo ""

# ===== Search API Tests =====
echo -e "${YELLOW}🔎 Search API Tests${NC}"

# Esperar indexación en Solr
echo -e "\n  ${CYAN}→ Esperando indexación en Solr (5s)...${NC}"
sleep 5

# Buscar actividades
echo -e "\n  ${CYAN}→ Buscando actividades...${NC}"
test_http "Buscar 'Fútbol'" "GET" "http://localhost:8083/search?q=Fútbol" "200"

test_http "Buscar 'parque'" "GET" "http://localhost:8083/search?q=parque" "200"

# Buscar con filtro de categoría
echo -e "\n  ${CYAN}→ Buscando con filtros...${NC}"
test_http "Buscar con categoría football" "GET" "http://localhost:8083/search?q=*&category=football" "200"

test_http "Buscar con dificultad intermediate" "GET" "http://localhost:8083/search?q=*&difficulty=intermediate" "200"

# Test de cache (segunda búsqueda debe venir del cache)
echo -e "\n  ${CYAN}→ Testing cache (segunda búsqueda)...${NC}"
test_http "Buscar 'Fútbol' (desde cache)" "GET" "http://localhost:8083/search?q=Fútbol" "200"

# Verificar que está en Solr
echo -e "\n  ${CYAN}→ Verificando indexación en Solr...${NC}"
SOLR_RESPONSE=$(curl -s "http://localhost:8983/solr/activities/select?q=*:*&rows=10")
SOLR_COUNT=$(echo "$SOLR_RESPONSE" | grep -o '"numFound":[0-9]*' | cut -d':' -f2 || echo "0")

if [ "$SOLR_COUNT" -gt 0 ]; then
    pass "Documentos indexados en Solr: $SOLR_COUNT"
else
    fail "No hay documentos en Solr" "$SOLR_RESPONSE"
fi

echo ""

# ===== Reservations API Tests =====
echo -e "${YELLOW}📅 Reservations API Tests${NC}"

# Crear reserva (como usuario autenticado)
if [ -n "$TOKEN" ] && [ -n "$ACTIVITY_ID" ]; then
    echo -e "\n  ${CYAN}→ Creando reserva...${NC}"
    CREATE_RESERVATION=$(test_http "Crear reserva" "POST" "http://localhost:8080/reservations" "201" \
        '{
            "activityId": "'"$ACTIVITY_ID"'",
            "userId": "testuser",
            "date": "2025-12-01T10:00:00Z",
            "participants": 5
        }' "Authorization: Bearer $TOKEN")

    RESERVATION_ID=$(echo "$CREATE_RESERVATION" | grep -o '"id":"[^"]*' | cut -d'"' -f4 || echo "")

    if [ -n "$RESERVATION_ID" ]; then
        pass "Reserva creada con ID: $RESERVATION_ID"
    fi

    # Listar reservas
    echo -e "\n  ${CYAN}→ Listando reservas...${NC}"
    test_http "Listar reservas" "GET" "http://localhost:8080/reservations" "200" "" "Authorization: Bearer $TOKEN"

    # Obtener reserva por ID
    if [ -n "$RESERVATION_ID" ]; then
        echo -e "\n  ${CYAN}→ Obteniendo reserva por ID...${NC}"
        test_http "Obtener reserva por ID" "GET" "http://localhost:8080/reservations/$RESERVATION_ID" "200" "" "Authorization: Bearer $TOKEN"
    fi
else
    fail "No se puede crear reserva sin token o activity ID" ""
fi

echo ""

# ===== Integration Tests =====
echo -e "${YELLOW}🔗 Integration Tests${NC}"

# Verificar que RabbitMQ tiene mensajes procesados
echo -e "\n  ${CYAN}→ Verificando mensajes en RabbitMQ...${NC}"
QUEUE_INFO=$(curl -s -u admin:admin123 http://localhost:15672/api/queues/%2F/search-sync)
MESSAGES_READY=$(echo "$QUEUE_INFO" | grep -o '"messages_ready":[0-9]*' | cut -d':' -f2 || echo "0")

if [ "$MESSAGES_READY" -eq 0 ]; then
    pass "Todos los mensajes de RabbitMQ fueron procesados"
else
    echo -e "  ${YELLOW}⚠${NC} Hay $MESSAGES_READY mensajes pendientes en RabbitMQ"
fi

# Verificar Memcached tiene entradas
echo -e "\n  ${CYAN}→ Verificando cache en Memcached...${NC}"
CACHE_STATS=$(echo -e "stats\r\nquit\r" | nc localhost 11211 | grep curr_items)
if [ -n "$CACHE_STATS" ]; then
    pass "Memcached tiene estadísticas disponibles"
    echo -e "    ${CYAN}$CACHE_STATS${NC}"
else
    fail "No se pudo obtener estadísticas de Memcached" ""
fi

# Test de flujo completo: Create Activity → Event → Index → Search
echo -e "\n  ${CYAN}→ Test de flujo completo (Create → Event → Index → Search)...${NC}"
if [ -n "$ADMIN_TOKEN" ]; then
    # Crear actividad única
    UNIQUE_NAME="Yoga Matutino $(date +%s)"
    CREATE_FLOW=$(curl -s -X POST "http://localhost:8082/activities" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -d '{
            "name": "'"$UNIQUE_NAME"'",
            "description": "Sesión de yoga al amanecer",
            "category": "yoga",
            "difficulty": "beginner",
            "duration": 60,
            "maxParticipants": 20,
            "location": "Playa",
            "price": 300.00
        }')

    FLOW_ACTIVITY_ID=$(echo "$CREATE_FLOW" | grep -o '"id":"[^"]*' | cut -d'"' -f4 || echo "")

    if [ -n "$FLOW_ACTIVITY_ID" ]; then
        pass "Actividad creada para test de flujo: $FLOW_ACTIVITY_ID"

        # Esperar indexación
        echo -e "    ${CYAN}Esperando indexación (8s)...${NC}"
        sleep 8

        # Buscar en Search API
        SEARCH_FLOW=$(curl -s "http://localhost:8083/search?q=Yoga%20Matutino")

        if echo "$SEARCH_FLOW" | grep -q "$FLOW_ACTIVITY_ID"; then
            pass "Actividad encontrada en Search API (flujo completo exitoso!)"
        else
            fail "Actividad NO encontrada en Search API después de 8s" "$SEARCH_FLOW"
        fi
    else
        fail "No se pudo crear actividad para test de flujo" "$CREATE_FLOW"
    fi
fi

echo ""

# ===== Resumen =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📊 Resumen de Tests${NC}"
echo -e "${BLUE}========================================${NC}\n"

TOTAL=$((PASSED + FAILED))

echo -e "  Total de tests: ${BLUE}$TOTAL${NC}"
echo -e "  Tests pasados:  ${GREEN}$PASSED${NC}"
echo -e "  Tests fallidos: ${RED}$FAILED${NC}\n"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ Todos los tests de backend pasaron!${NC}\n"
    echo -e "${CYAN}📝 Datos de prueba creados:${NC}"
    echo -e "  - Usuario: test@example.com / Test123!"
    echo -e "  - Admin: admin@example.com / Admin123!"
    [ -n "$ACTIVITY_ID" ] && echo -e "  - Actividad ID: $ACTIVITY_ID"
    [ -n "$RESERVATION_ID" ] && echo -e "  - Reserva ID: $RESERVATION_ID"
    echo ""
    exit 0
else
    echo -e "${RED}❌ Algunos tests fallaron. Revisa los logs de los servicios.${NC}\n"
    echo -e "${YELLOW}Comandos útiles para debugging:${NC}"
    echo -e "  docker-compose logs users-api"
    echo -e "  docker-compose logs activities-api"
    echo -e "  docker-compose logs search-api"
    echo -e "  docker-compose logs reservations-service"
    echo ""
    exit 1
fi
