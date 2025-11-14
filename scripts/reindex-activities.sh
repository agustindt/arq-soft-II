#!/bin/bash

# Script para reindexar todas las actividades de MongoDB a Solr
# Uso: ./scripts/reindex-activities.sh

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🔄 Reindexando Actividades en Solr${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Verificar que los servicios estén corriendo
if ! curl -sf http://localhost:8082/healthz > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Activities API no está corriendo${NC}"
    exit 1
fi

if ! curl -sf http://localhost:8983/solr/activities/admin/ping > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Solr no está corriendo${NC}"
    exit 1
fi

echo -e "${YELLOW}→ Obteniendo actividades de MongoDB...${NC}"

# Obtener todas las actividades y guardarlas en un archivo temporal
TEMP_FILE=$(mktemp)
curl -s http://localhost:8082/activities > "$TEMP_FILE"

# Contar actividades usando jq si está disponible, o python
if command -v jq &> /dev/null; then
    COUNT=$(jq '.activities | length' "$TEMP_FILE")
else
    COUNT=$(python3 -c "import json, sys; data=json.load(open('$TEMP_FILE')); print(len(data.get('activities', [])))" 2>/dev/null || echo "0")
fi

if [ "$COUNT" -eq "0" ]; then
    echo -e "${YELLOW}⚠ No hay actividades para indexar${NC}"
    rm "$TEMP_FILE"
    exit 0
fi

echo -e "${GREEN}✓ Encontradas ${COUNT} actividades${NC}\n"
echo -e "${YELLOW}→ Indexando actividades en Solr...${NC}"

# Usar Python para procesar y indexar
python3 << PYTHON_SCRIPT
import json
import subprocess
import sys

# Leer actividades
with open('$TEMP_FILE', 'r') as f:
    data = json.load(f)
    activities = data.get('activities', [])

indexed = 0
errors = 0

for activity in activities:
    # Crear documento para Solr
    solr_doc = {
        "id": activity["id"],
        "name": activity["name"],
        "description": activity["description"],
        "category": activity["category"],
        "difficulty": activity["difficulty"],
        "location": activity["location"],
        "price": float(activity["price"]),
        "date_created": activity["created_at"]
    }
    
    # Indexar usando curl
    json_data = json.dumps([solr_doc])
    result = subprocess.run(
        ['curl', '-s', '-X', 'POST',
         'http://localhost:8983/solr/activities/update?commit=true',
         '-H', 'Content-Type: application/json',
         '-d', json_data],
        capture_output=True,
        text=True
    )
    
    if result.returncode == 0:
        try:
            response = json.loads(result.stdout)
            if response.get('responseHeader', {}).get('status') == 0:
                indexed += 1
                print(f"  ✓ {activity['name']}")
            else:
                errors += 1
                print(f"  ✗ {activity['name']}: Error en Solr")
        except:
            errors += 1
            print(f"  ✗ {activity['name']}: Error parseando respuesta")
    else:
        errors += 1
        print(f"  ✗ {activity['name']}: Error en curl")

print(f"\n${GREEN}✅ Indexadas {indexed}/{len(activities)} actividades${NC}")
if errors > 0:
    print(f"${YELLOW}⚠ {errors} errores${NC}")
PYTHON_SCRIPT

rm "$TEMP_FILE"

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Reindexación completada!${NC}"
echo -e "${GREEN}========================================${NC}\n"
