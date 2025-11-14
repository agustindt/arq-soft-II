#!/bin/bash

# Script para testear toda la infraestructura
# Uso: ./scripts/test-infrastructure.sh

set -e

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

# Función para reportar test exitoso
pass() {
    echo -e "  ${GREEN}✓${NC} $1"
    ((PASSED++))
}

# Función para reportar test fallido
fail() {
    echo -e "  ${RED}✗${NC} $1"
    ((FAILED++))
}

# Función para ejecutar test
test_command() {
    local description="$1"
    local command="$2"

    if eval "$command" > /dev/null 2>&1; then
        pass "$description"
    else
        fail "$description"
    fi
}

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🧪 Testing Infraestructura${NC}"
echo -e "${BLUE}========================================${NC}\n"

# ===== MySQL Tests =====
echo -e "${YELLOW}🗄️  MySQL Tests${NC}"

test_command "MySQL contenedor corriendo" \
    "docker-compose ps mysql | grep -q 'Up'"

test_command "MySQL puerto 3307 accesible" \
    "nc -z localhost 3307"

test_command "MySQL acepta conexiones" \
    "docker-compose exec -T mysql mysqladmin ping -h localhost --silent"

test_command "Base de datos 'users_db' existe" \
    "docker-compose exec -T mysql mysql -uroot -prootpassword -e 'SHOW DATABASES;' | grep -q 'users_db'"

test_command "Tabla 'users' existe en users_db" \
    "docker-compose exec -T mysql mysql -uroot -prootpassword -D users_db -e 'SHOW TABLES;' | grep -q 'users'"

echo ""

# ===== MongoDB Tests =====
echo -e "${YELLOW}🍃 MongoDB Tests${NC}"

test_command "MongoDB contenedor corriendo" \
    "docker-compose ps mongo | grep -q 'Up'"

test_command "MongoDB puerto 27017 accesible" \
    "nc -z localhost 27017"

test_command "MongoDB acepta conexiones" \
    "docker-compose exec -T mongo mongosh --quiet --eval 'db.version()'"

test_command "Database 'activitiesdb' accesible" \
    "docker-compose exec -T mongo mongosh activitiesdb --quiet --eval 'db.getName()' | grep -q 'activitiesdb'"

test_command "Database 'reservasdb' accesible" \
    "docker-compose exec -T mongo mongosh reservasdb --quiet --eval 'db.getName()' | grep -q 'reservasdb'"

echo ""

# ===== RabbitMQ Tests =====
echo -e "${YELLOW}🐇 RabbitMQ Tests${NC}"

test_command "RabbitMQ contenedor corriendo" \
    "docker-compose ps rabbitmq | grep -q 'Up'"

test_command "RabbitMQ puerto 5672 accesible" \
    "nc -z localhost 5672"

test_command "RabbitMQ puerto 15672 (Management UI) accesible" \
    "nc -z localhost 15672"

test_command "RabbitMQ health check OK" \
    "docker-compose exec -T rabbitmq rabbitmq-diagnostics -q ping"

test_command "RabbitMQ Management UI responde" \
    "curl -sf -u admin:admin123 http://localhost:15672/api/overview"

test_command "RabbitMQ exchange 'entity.events' existe" \
    "curl -sf -u admin:admin123 http://localhost:15672/api/exchanges/%2F/entity.events"

test_command "RabbitMQ queue 'search-sync' existe" \
    "curl -sf -u admin:admin123 http://localhost:15672/api/queues/%2F/search-sync"

echo ""

# ===== Solr Tests =====
echo -e "${YELLOW}🔆 Solr Tests${NC}"

test_command "Solr contenedor corriendo" \
    "docker-compose ps solr | grep -q 'Up'"

test_command "Solr puerto 8983 accesible" \
    "nc -z localhost 8983"

test_command "Solr Admin UI responde" \
    "curl -sf http://localhost:8983/solr/"

test_command "Solr core 'activities' existe" \
    "curl -sf http://localhost:8983/solr/admin/cores?action=STATUS | grep -q 'activities'"

test_command "Solr core 'activities' ping OK" \
    "curl -sf http://localhost:8983/solr/activities/admin/ping | grep -q 'OK'"

test_command "Solr core 'activities' schema cargado" \
    "curl -sf http://localhost:8983/solr/activities/schema | grep -q 'name'"

echo ""

# ===== Memcached Tests =====
echo -e "${YELLOW}💾 Memcached Tests${NC}"

test_command "Memcached contenedor corriendo" \
    "docker-compose ps memcached | grep -q 'Up'"

test_command "Memcached puerto 11211 accesible" \
    "nc -z localhost 11211"

# Test de escritura/lectura en Memcached
test_command "Memcached acepta comandos (set/get)" \
    "(echo -e 'set test_key 0 60 5\r\nhello\r\nget test_key\r\nquit\r' | nc localhost 11211 | grep -q 'hello')"

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
    echo -e "${GREEN}✅ Todos los tests de infraestructura pasaron!${NC}\n"
    exit 0
else
    echo -e "${RED}❌ Algunos tests fallaron. Revisa la configuración.${NC}\n"
    exit 1
fi
