#!/bin/bash

# Script para levantar todo el stack completo
# Uso: ./scripts/start-all.sh

# Ensure script is run with bash
if [ -z "$BASH_VERSION" ]; then
    exec bash "$0" "$@"
fi

set -e

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Helper function for colored output (using printf for better compatibility)
print_colored() {
    local color=$1
    shift
    printf "${color}%s${NC}\n" "$@"
}

print_colored_no_newline() {
    local color=$1
    shift
    printf "${color}%s${NC}" "$@"
}

# Helper function to wait for service with timeout
wait_for_service() {
    local service_name=$1
    local check_command=$2
    local max_attempts=${3:-60}  # Default 60 attempts = 120 seconds
    local attempt=0
    
    print_colored_no_newline "$YELLOW" "   $service_name: "
    
    while [ $attempt -lt $max_attempts ]; do
        if eval "$check_command" > /dev/null 2>&1; then
            printf " ${GREEN}✓${NC}\n"
            return 0
        fi
        printf "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    printf " ${RED}✗${NC}\n"
    print_colored "$YELLOW" "   ⚠️  Timeout esperando $service_name (${max_attempts}s). Continuando..."
    return 1
}

print_colored "$BLUE" "========================================"
print_colored "$BLUE" "🚀 Iniciando Stack Completo"
print_colored "$BLUE" "========================================"
printf "\n"

# Verificar que docker está corriendo
if ! docker info > /dev/null 2>&1; then
    print_colored "$RED" "❌ Error: Docker no está corriendo"
    print_colored "$YELLOW" "Por favor inicia Docker Desktop e intenta nuevamente"
    exit 1
fi

# Detener contenedores anteriores si existen
print_colored "$YELLOW" "🧹 Limpiando contenedores anteriores..."
docker-compose down 2>/dev/null || true

# Limpiar volúmenes órfanos
print_colored "$YELLOW" "🧹 Limpiando volúmenes no utilizados..."
docker volume prune -f > /dev/null 2>&1 || true

printf "\n"
print_colored "$BLUE" "========================================"
print_colored "$BLUE" "📦 Fase 1: Infraestructura"
print_colored "$BLUE" "========================================"
printf "\n"

# Levantar infraestructura primero
print_colored "$YELLOW" "🗄️  Levantando bases de datos..."
docker-compose up -d mysql mongo

print_colored "$YELLOW" "🐇 Levantando RabbitMQ..."
docker-compose up -d rabbitmq

print_colored "$YELLOW" "🔆 Levantando Solr..."
docker-compose up -d solr

print_colored "$YELLOW" "💾 Levantando Memcached..."
docker-compose up -d memcached

# Esperar a que la infraestructura esté lista
printf "\n"
print_colored "$YELLOW" "⏳ Esperando a que la infraestructura esté lista..."

# Esperar MySQL
wait_for_service "MySQL" "docker-compose exec -T mysql mysqladmin ping -h localhost --silent"

# Esperar RabbitMQ
wait_for_service "RabbitMQ" "docker-compose exec -T rabbitmq rabbitmq-diagnostics -q ping"

# Esperar Solr con lógica mejorada
print_colored_no_newline "$YELLOW" "   Solr: "

# First check if Solr container is running
SOLR_ATTEMPT=0
MAX_SOLR_ATTEMPTS=90  # 180 seconds for Solr (it can take longer)
SOLR_READY=0

while [ $SOLR_ATTEMPT -lt $MAX_SOLR_ATTEMPTS ]; do
    # Step 1: Check if container is running
    if ! docker-compose ps solr | grep -q "Up"; then
        printf "."
        sleep 2
        SOLR_ATTEMPT=$((SOLR_ATTEMPT + 1))
        continue
    fi
    
    # Step 2: Check if Solr admin UI responds
    if ! curl -sf http://localhost:8983/solr/ > /dev/null 2>&1; then
        printf "."
        sleep 2
        SOLR_ATTEMPT=$((SOLR_ATTEMPT + 1))
        continue
    fi
    
    # Step 3: Check if core exists and ping works
    if curl -sf http://localhost:8983/solr/activities/admin/ping > /dev/null 2>&1; then
        printf " ${GREEN}✓${NC}\n"
        SOLR_READY=1
        break
    fi
    
    printf "."
    sleep 2
    SOLR_ATTEMPT=$((SOLR_ATTEMPT + 1))
done

if [ $SOLR_READY -eq 0 ]; then
    printf " ${YELLOW}⚠${NC}\n"
    print_colored "$YELLOW" "   ⚠️  Solr tardó más de lo esperado. Puede que el core aún se esté inicializando."
    print_colored "$YELLOW" "   Continuando... (puedes verificar manualmente: http://localhost:8983)"
fi

printf "\n"
print_colored "$BLUE" "========================================"
print_colored "$BLUE" "🔧 Fase 2: Backend Services"
print_colored "$BLUE" "========================================"
printf "\n"

# Levantar backend services
print_colored "$YELLOW" "👤 Levantando Users API..."
docker-compose up -d users-api

# Esperar users-api
wait_for_service "Users API" "curl -sf http://localhost:8081/api/v1/health"

print_colored "$YELLOW" "🏃 Levantando Activities API..."
docker-compose up -d activities-api

# Esperar activities-api
wait_for_service "Activities API" "curl -sf http://localhost:8082/healthz"

print_colored "$YELLOW" "🔎 Levantando Search API..."
docker-compose up -d search-api

# Esperar search-api
wait_for_service "Search API" "curl -sf http://localhost:8083/health"

print_colored "$YELLOW" "📅 Levantando Reservations API..."
docker-compose up -d reservations-service

# Esperar reservations-service
wait_for_service "Reservations API" "curl -sf http://localhost:8080/healthz"

printf "\n"
print_colored "$BLUE" "========================================"
print_colored "$BLUE" "⚛️  Fase 3: Frontend"
print_colored "$BLUE" "========================================"
printf "\n"

print_colored "$YELLOW" "🌐 Levantando Frontend React..."
docker-compose up -d frontend

# Esperar frontend
wait_for_service "Frontend" "curl -sf http://localhost:3000"

printf "\n"
print_colored "$GREEN" "========================================"
print_colored "$GREEN" "✅ Stack completo levantado exitosamente!"
print_colored "$GREEN" "========================================"
printf "\n"

print_colored "$BLUE" "📍 URLs de acceso:"
printf "\n"
print_colored "$YELLOW" "  Frontend:"
printf "    🌐 React App:           ${GREEN}http://localhost:3000${NC}\n"
printf "\n"

print_colored "$YELLOW" "  Backend APIs:"
printf "    👤 Users API:           ${GREEN}http://localhost:8081${NC}\n"
printf "    🏃 Activities API:      ${GREEN}http://localhost:8082${NC}\n"
printf "    🔎 Search API:          ${GREEN}http://localhost:8083${NC}\n"
printf "    📅 Reservations API:    ${GREEN}http://localhost:8080${NC}\n"
printf "\n"

print_colored "$YELLOW" "  Infraestructura:"
printf "    🗄️  MySQL:              ${GREEN}localhost:3307${NC}\n"
printf "    🍃 MongoDB:             ${GREEN}localhost:27017${NC}\n"
printf "    🐇 RabbitMQ Management: ${GREEN}http://localhost:15672${NC} (admin/admin123)\n"
printf "    🔆 Solr Admin:          ${GREEN}http://localhost:8983${NC}\n"
printf "    💾 Memcached:           ${GREEN}localhost:11211${NC}\n"
printf "\n"

print_colored "$BLUE" "📋 Comandos útiles:"
printf "  ${YELLOW}Ver logs:${NC}              docker-compose logs -f [servicio]\n"
printf "  ${YELLOW}Ver todos los logs:${NC}    docker-compose logs -f\n"
printf "  ${YELLOW}Detener todo:${NC}          docker-compose down\n"
printf "  ${YELLOW}Reiniciar servicio:${NC}    docker-compose restart [servicio]\n"
printf "  ${YELLOW}Ver estado:${NC}            docker-compose ps\n"
printf "\n"

printf "${YELLOW}💡 Tip:${NC} Ejecuta ${GREEN}./scripts/test-all.sh${NC} para testear todo el stack\n"
printf "\n"
