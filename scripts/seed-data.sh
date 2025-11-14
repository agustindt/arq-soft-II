#!/bin/bash

# Script para cargar datos de prueba (seed data) en las bases de datos
# Uso: ./scripts/seed-data.sh

set -e

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🌱 Seeding Database Data${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Verificar que docker está corriendo
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Docker no está corriendo${NC}"
    echo -e "${YELLOW}Por favor inicia Docker Desktop e intenta nuevamente${NC}"
    exit 1
fi

# Verificar que los contenedores están corriendo
echo -e "${YELLOW}🔍 Verificando que los contenedores están corriendo...${NC}\n"

if ! docker-compose ps mysql | grep -q "Up"; then
    echo -e "${RED}❌ Error: Contenedor MySQL no está corriendo${NC}"
    echo -e "${YELLOW}Por favor inicia los servicios primero: ${CYAN}docker-compose up -d mysql${NC}"
    exit 1
fi

if ! docker-compose ps mongo | grep -q "Up"; then
    echo -e "${RED}❌ Error: Contenedor MongoDB no está corriendo${NC}"
    echo -e "${YELLOW}Por favor inicia los servicios primero: ${CYAN}docker-compose up -d mongo${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Contenedores verificados${NC}\n"

# ===== MySQL Seed Data =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🗄️  MySQL Seed Data${NC}"
echo -e "${BLUE}========================================${NC}\n"

MYSQL_SEED_FILE="$PROJECT_ROOT/database/mysql/seed.sql"

if [ ! -f "$MYSQL_SEED_FILE" ]; then
    echo -e "${RED}❌ Error: Archivo seed.sql no encontrado en $MYSQL_SEED_FILE${NC}"
    exit 1
fi

echo -e "${YELLOW}→ Cargando datos de prueba en MySQL...${NC}"

# Esperar a que MySQL esté listo
echo -n "   Esperando MySQL: "
until docker-compose exec -T mysql mysqladmin ping -h localhost --silent > /dev/null 2>&1; do
    echo -n "."
    sleep 1
done
echo -e " ${GREEN}✓${NC}"

# Cargar seed data
if docker-compose exec -T mysql mysql -uroot -prootpassword users_db < "$MYSQL_SEED_FILE" 2>/dev/null; then
    echo -e "  ${GREEN}✓${NC} Datos de MySQL cargados exitosamente"
    
    # Verificar datos insertados
    USER_COUNT=$(docker-compose exec -T mysql mysql -uroot -prootpassword users_db -se "SELECT COUNT(*) FROM users WHERE email LIKE '%@test.com';" 2>/dev/null || echo "0")
    echo -e "  ${CYAN}→ Usuarios de prueba: ${USER_COUNT}${NC}"
else
    echo -e "  ${YELLOW}⚠${NC} Advertencia: Algunos datos pueden no haberse cargado (puede ser normal si ya existen)"
fi

echo ""

# ===== MongoDB Seed Data =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🍃 MongoDB Seed Data${NC}"
echo -e "${BLUE}========================================${NC}\n"

MONGO_SEED_FILE="$PROJECT_ROOT/database/mongo/seed.js"

if [ ! -f "$MONGO_SEED_FILE" ]; then
    echo -e "${RED}❌ Error: Archivo seed.js no encontrado en $MONGO_SEED_FILE${NC}"
    exit 1
fi

echo -e "${YELLOW}→ Cargando datos de prueba en MongoDB...${NC}"

# Esperar a que MongoDB esté listo
echo -n "   Esperando MongoDB: "
until docker-compose exec -T mongo mongosh --quiet --eval "db.version()" > /dev/null 2>&1; do
    echo -n "."
    sleep 1
done
echo -e " ${GREEN}✓${NC}"

# Cargar seed data
if docker-compose exec -T mongo mongosh < "$MONGO_SEED_FILE" > /dev/null 2>&1; then
    echo -e "  ${GREEN}✓${NC} Datos de MongoDB cargados exitosamente"
    
    # Verificar datos insertados
    ACTIVITIES_COUNT=$(docker-compose exec -T mongo mongosh activitiesdb --quiet --eval "db.activities.countDocuments()" 2>/dev/null | tr -d '\r\n' || echo "0")
    RESERVATIONS_COUNT=$(docker-compose exec -T mongo mongosh reservasdb --quiet --eval "db.reservas.countDocuments()" 2>/dev/null | tr -d '\r\n' || echo "0")
    
    echo -e "  ${CYAN}→ Actividades: ${ACTIVITIES_COUNT}${NC}"
    echo -e "  ${CYAN}→ Reservas: ${RESERVATIONS_COUNT}${NC}"
else
    echo -e "  ${YELLOW}⚠${NC} Advertencia: Algunos datos pueden no haberse cargado (puede ser normal si ya existen)"
fi

echo ""

# ===== Resumen Final =====
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Seeding completado!${NC}"
echo -e "${GREEN}========================================${NC}\n"

echo -e "${CYAN}📊 Resumen:${NC}"
echo -e "  • MySQL: Datos de usuarios cargados"
echo -e "  • MongoDB: Actividades y reservas cargadas\n"

echo -e "${YELLOW}💡 Notas:${NC}"
echo -e "  • Los datos de prueba se cargan en las bases de datos"
echo -e "  • Si los datos ya existen, algunos inserts pueden fallar (normal)"
echo -e "  • Para limpiar y recargar: ${CYAN}docker-compose down -v && docker-compose up -d${NC}\n"

echo -e "${BLUE}📝 Datos de prueba disponibles:${NC}"
echo -e "  • Usuarios de prueba (MySQL): Ver seed.sql para detalles"
echo -e "  • Actividades de prueba (MongoDB): Ver seed.js para detalles"
echo -e "  • Reservas de prueba (MongoDB): Ver seed.js para detalles\n"

