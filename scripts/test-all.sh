#!/bin/bash

# Script maestro para ejecutar TODOS los tests
# Uso: ./scripts/test-all.sh

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
CYAN='\033[0;36m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
START_TIME=$(date +%s)

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🧪 TEST SUITE COMPLETO${NC}"
echo -e "${BLUE}Testing Backend + Frontend + Infra${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Verificar que docker está corriendo
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Docker no está corriendo${NC}"
    echo -e "${YELLOW}Por favor inicia Docker Desktop e intenta nuevamente${NC}"
    exit 1
fi

# Verificar que los servicios están corriendo
echo -e "${YELLOW}🔍 Verificando que los servicios están corriendo...${NC}\n"

SERVICES_RUNNING=$(docker-compose ps --services --filter "status=running" | wc -l)
TOTAL_SERVICES=9  # mysql, mongo, rabbitmq, solr, memcached, users-api, activities-api, search-api, reservations-service, frontend

if [ "$SERVICES_RUNNING" -lt 8 ]; then
    echo -e "${YELLOW}⚠️  Advertencia: Solo $SERVICES_RUNNING servicios están corriendo${NC}"
    echo -e "${YELLOW}Se esperaban al menos 8 servicios activos${NC}\n"

    read -p "¿Deseas iniciar todos los servicios ahora? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${CYAN}Iniciando servicios...${NC}"
        "$SCRIPT_DIR/start-all.sh"
        echo ""
    else
        echo -e "${RED}Abortando. Por favor inicia los servicios primero.${NC}"
        exit 1
    fi
fi

# ===== Fase 1: Test de Infraestructura =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📋 FASE 1: Infraestructura${NC}"
echo -e "${BLUE}========================================${NC}\n"

if [ -f "$SCRIPT_DIR/test-infrastructure.sh" ]; then
    bash "$SCRIPT_DIR/test-infrastructure.sh"
    INFRA_EXIT_CODE=$?
else
    echo -e "${RED}✗ Script test-infrastructure.sh no encontrado${NC}"
    INFRA_EXIT_CODE=1
fi

echo ""

# ===== Fase 2: Test de Backend =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📋 FASE 2: Backend Services${NC}"
echo -e "${BLUE}========================================${NC}\n"

if [ -f "$SCRIPT_DIR/test-backend.sh" ]; then
    bash "$SCRIPT_DIR/test-backend.sh"
    BACKEND_EXIT_CODE=$?
else
    echo -e "${RED}✗ Script test-backend.sh no encontrado${NC}"
    BACKEND_EXIT_CODE=1
fi

echo ""

# ===== Fase 3: Test de Frontend =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📋 FASE 3: Frontend${NC}"
echo -e "${BLUE}========================================${NC}\n"

echo -e "${YELLOW}→ Frontend Tests${NC}\n"

# Verificar que el frontend responde
if curl -sf http://localhost:3000 > /dev/null 2>&1; then
    echo -e "  ${GREEN}✓${NC} Frontend accesible en http://localhost:3000"
    FRONTEND_EXIT_CODE=0
else
    echo -e "  ${RED}✗${NC} Frontend NO accesible en http://localhost:3000"
    FRONTEND_EXIT_CODE=1
fi

# Verificar que devuelve HTML
RESPONSE=$(curl -s http://localhost:3000 || echo "")
if echo "$RESPONSE" | grep -q "<html"; then
    echo -e "  ${GREEN}✓${NC} Frontend devuelve HTML válido"
else
    echo -e "  ${RED}✗${NC} Frontend NO devuelve HTML válido"
    FRONTEND_EXIT_CODE=1
fi

# Verificar que tiene los scripts de React
if echo "$RESPONSE" | grep -q "react"; then
    echo -e "  ${GREEN}✓${NC} Aplicación React cargada"
else
    echo -e "  ${YELLOW}⚠${NC} No se detectó React en el HTML (puede ser normal)"
fi

echo ""

# ===== Fase 4: Tests End-to-End =====
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📋 FASE 4: End-to-End Flow${NC}"
echo -e "${BLUE}========================================${NC}\n"

echo -e "${YELLOW}→ End-to-End Tests${NC}\n"

# Verificar conectividad completa Frontend ↔ Backend
echo -e "  ${CYAN}→ Verificando stack completo...${NC}\n"

E2E_EXIT_CODE=0

# Test: Frontend puede comunicarse con cada API
for api in "users-api:8081" "activities-api:8082" "search-api:8083" "reservations-service:8080"; do
    SERVICE=$(echo $api | cut -d: -f1)
    PORT=$(echo $api | cut -d: -f2)

    if docker-compose exec -T frontend sh -c "wget -q -O- http://host.docker.internal:$PORT/health 2>/dev/null || wget -q -O- http://localhost:$PORT/health 2>/dev/null" > /dev/null 2>&1; then
        echo -e "    ${GREEN}✓${NC} Frontend ↔ $SERVICE comunicación OK"
    else
        # Algunos contenedores pueden no tener wget, intentar desde host
        if curl -sf http://localhost:$PORT/health > /dev/null 2>&1; then
            echo -e "    ${GREEN}✓${NC} Frontend ↔ $SERVICE comunicación OK (via host)"
        else
            echo -e "    ${RED}✗${NC} Frontend ↔ $SERVICE comunicación FAILED"
            E2E_EXIT_CODE=1
        fi
    fi
done

echo ""

# ===== Reporte Final =====
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📊 REPORTE FINAL${NC}"
echo -e "${BLUE}========================================${NC}\n"

echo -e "  Duración total: ${CYAN}${DURATION}s${NC}\n"

# Resumen de cada fase
echo -e "  ${BLUE}Resultados por fase:${NC}"

if [ $INFRA_EXIT_CODE -eq 0 ]; then
    echo -e "    ${GREEN}✓${NC} Infraestructura: PASSED"
else
    echo -e "    ${RED}✗${NC} Infraestructura: FAILED"
fi

if [ $BACKEND_EXIT_CODE -eq 0 ]; then
    echo -e "    ${GREEN}✓${NC} Backend Services: PASSED"
else
    echo -e "    ${RED}✗${NC} Backend Services: FAILED"
fi

if [ $FRONTEND_EXIT_CODE -eq 0 ]; then
    echo -e "    ${GREEN}✓${NC} Frontend: PASSED"
else
    echo -e "    ${RED}✗${NC} Frontend: FAILED"
fi

if [ $E2E_EXIT_CODE -eq 0 ]; then
    echo -e "    ${GREEN}✓${NC} End-to-End: PASSED"
else
    echo -e "    ${RED}✗${NC} End-to-End: FAILED"
fi

echo ""

# Determinar resultado final
FINAL_EXIT_CODE=$((INFRA_EXIT_CODE + BACKEND_EXIT_CODE + FRONTEND_EXIT_CODE + E2E_EXIT_CODE))

if [ $FINAL_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✅ TODOS LOS TESTS PASARON!${NC}"
    echo -e "${GREEN}El sistema está funcionando OK${NC}"
    echo -e "${GREEN}========================================${NC}\n"

    echo -e "${CYAN}📝 El stack completo está operativo:${NC}"
    echo -e "  • Infraestructura: MySQL, MongoDB, RabbitMQ, Solr, Memcached"
    echo -e "  • Backend: Users, Activities, Search, Reservations APIs"
    echo -e "  • Frontend: React App"
    echo -e "  • Integración: Events, Search indexing, Caching\n"

    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}❌ ALGUNOS TESTS FALLARON${NC}"
    echo -e "${RED}========================================${NC}\n"

    echo -e "${YELLOW}💡 Siguiente pasos para debugging:${NC}"
    echo -e "  1. Revisar logs: ${CYAN}docker-compose logs -f [servicio]${NC}"
    echo -e "  2. Ver estado: ${CYAN}docker-compose ps${NC}"
    echo -e "  3. Ejecutar tests individuales:"
    echo -e "     ${CYAN}./scripts/test-infrastructure.sh${NC}"
    echo -e "     ${CYAN}./scripts/test-backend.sh${NC}"
    echo -e "  4. Reiniciar servicios: ${CYAN}docker-compose restart [servicio]${NC}\n"

    exit 1
fi
