#!/bin/bash

# Script interactivo para gestionar datos de prueba
# Permite crear y eliminar usuarios, actividades y reservas

set -e

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Variables
USERS_API="http://localhost:8081/api/v1"
ACTIVITIES_API="http://localhost:8082"
RESERVATIONS_API="http://localhost:8080"

# Función para obtener token de admin
get_admin_token() {
    local email=${1:-"admin@example.com"}
    local password=${2:-"password"}
    
    token=$(curl -s -X POST "$USERS_API/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$email\", \"password\": \"$password\"}" \
        | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$token" ] || [ "$token" = "null" ]; then
        echo -e "${RED}❌ Error al obtener token. Verifica las credenciales.${NC}"
        exit 1
    fi
    
    echo "$token"
}

# Función para obtener token de root
get_root_token() {
    get_admin_token "root@example.com" "password"
}

# ===== GESTIÓN DE USUARIOS =====

create_user() {
    echo -e "\n${CYAN}📝 Crear Nuevo Usuario${NC}"
    echo "=============================="
    
    read -p "Email: " email
    read -p "Username: " username
    read -sp "Password: " password
    echo ""
    read -p "Nombre: " first_name
    read -p "Apellido: " last_name
    echo ""
    echo "Selecciona el rol:"
    echo "1) user"
    echo "2) admin"
    echo "3) root"
    read -p "Opción [1-3]: " role_option
    
    case $role_option in
        1) role="user" ;;
        2) role="admin" ;;
        3) role="root" ;;
        *) echo -e "${RED}Opción inválida${NC}"; return 1 ;;
    esac
    
    # Obtener token de admin/root
    if [ "$role" = "root" ]; then
        echo -e "\n${YELLOW}Creando usuario ROOT (requiere token root)...${NC}"
        TOKEN=$(get_root_token)
        
        response=$(curl -s -X POST "$USERS_API/admin/create-root" \
            -H "Content-Type: application/json" \
            -d "{
                \"email\": \"$email\",
                \"username\": \"$username\",
                \"password\": \"$password\",
                \"first_name\": \"$first_name\",
                \"last_name\": \"$last_name\",
                \"secret_key\": \"create-root-secret-2024\"
            }")
    else
        echo -e "\n${YELLOW}Creando usuario $role...${NC}"
        TOKEN=$(get_admin_token)
        
        response=$(curl -s -X POST "$USERS_API/admin/users" \
            -H "Authorization: Bearer $TOKEN" \
            -H "Content-Type: application/json" \
            -d "{
                \"email\": \"$email\",
                \"username\": \"$username\",
                \"password\": \"$password\",
                \"first_name\": \"$first_name\",
                \"last_name\": \"$last_name\",
                \"role\": \"$role\"
            }")
    fi
    
    if echo "$response" | grep -q "error"; then
        echo -e "${RED}❌ Error: $(echo $response | grep -o '"message":"[^"]*"' | cut -d'"' -f4)${NC}"
    else
        echo -e "${GREEN}✅ Usuario creado exitosamente!${NC}"
        echo -e "   Email: $email"
        echo -e "   Username: $username"
        echo -e "   Rol: $role"
    fi
}

delete_user() {
    echo -e "\n${CYAN}🗑️  Eliminar Usuario${NC}"
    echo "=============================="
    
    read -p "Email del usuario a eliminar: " email
    
    # Obtener token de root (solo root puede eliminar)
    TOKEN=$(get_root_token)
    
    # Obtener ID del usuario
    user_id=$(docker exec arq-soft-ii-mysql-1 mysql -u root -prootpassword -N -e \
        "SELECT id FROM users_db.users WHERE email='$email';" 2>/dev/null)
    
    if [ -z "$user_id" ]; then
        echo -e "${RED}❌ Usuario no encontrado${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Usuario encontrado (ID: $user_id)${NC}"
    read -p "¿Confirmar eliminación? (y/n): " confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        response=$(curl -s -X DELETE "$USERS_API/admin/users/$user_id" \
            -H "Authorization: Bearer $TOKEN")
        
        if echo "$response" | grep -q "error"; then
            echo -e "${RED}❌ Error al eliminar usuario${NC}"
        else
            echo -e "${GREEN}✅ Usuario eliminado exitosamente${NC}"
        fi
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

list_users() {
    echo -e "\n${CYAN}👥 Lista de Usuarios${NC}"
    echo "=============================="
    
    docker exec arq-soft-ii-mysql-1 mysql -u root -prootpassword -e \
        "SELECT id, username, email, role, is_active FROM users_db.users ORDER BY role, username;" \
        2>/dev/null | column -t
}

# ===== GESTIÓN DE ACTIVIDADES =====

create_activity() {
    echo -e "\n${CYAN}🏃 Crear Nueva Actividad${NC}"
    echo "=============================="
    
    read -p "Nombre: " name
    read -p "Descripción: " description
    read -p "Categoría (sports/yoga/fitness/etc): " category
    read -p "Dificultad (beginner/intermediate/advanced): " difficulty
    read -p "Ubicación: " location
    read -p "Precio: " price
    read -p "Duración (minutos): " duration
    read -p "Capacidad máxima: " max_capacity
    read -p "Instructor: " instructor
    
    # Obtener token de admin
    TOKEN=$(get_admin_token)
    
    response=$(curl -s -X POST "$ACTIVITIES_API/activities" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"$name\",
            \"description\": \"$description\",
            \"category\": \"$category\",
            \"difficulty\": \"$difficulty\",
            \"location\": \"$location\",
            \"price\": $price,
            \"duration\": $duration,
            \"max_capacity\": $max_capacity,
            \"instructor\": \"$instructor\",
            \"schedule\": [],
            \"equipment\": [],
            \"image_url\": \"\",
            \"is_active\": true
        }")
    
    if echo "$response" | grep -q "error"; then
        echo -e "${RED}❌ Error al crear actividad${NC}"
    else
        echo -e "${GREEN}✅ Actividad creada exitosamente!${NC}"
        activity_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        echo -e "   ID: $activity_id"
        echo -e "   Nombre: $name"
    fi
}

delete_activity() {
    echo -e "\n${CYAN}🗑️  Eliminar Actividad${NC}"
    echo "=============================="
    
    read -p "ID de la actividad a eliminar: " activity_id
    
    # Obtener token de admin
    TOKEN=$(get_admin_token)
    
    read -p "¿Confirmar eliminación? (y/n): " confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        response=$(curl -s -X DELETE "$ACTIVITIES_API/activities/$activity_id" \
            -H "Authorization: Bearer $TOKEN")
        
        if echo "$response" | grep -q "error"; then
            echo -e "${RED}❌ Error al eliminar actividad${NC}"
        else
            echo -e "${GREEN}✅ Actividad eliminada exitosamente${NC}"
        fi
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

list_activities() {
    echo -e "\n${CYAN}🏃 Lista de Actividades${NC}"
    echo "=============================="
    
    docker exec mongo mongosh activitiesdb --quiet --eval \
        "db.activities.find({}, {_id: 1, name: 1, category: 1, price: 1, is_active: 1, created_by: 1}).forEach(a => print(a._id + ' | ' + a.name.padEnd(30) + ' | ' + a.category.padEnd(15) + ' | $' + a.price + ' | Active: ' + a.is_active + ' | Owner: ' + a.created_by))"
}

# ===== GESTIÓN DE RESERVAS =====

delete_reservation() {
    echo -e "\n${CYAN}🗑️  Eliminar Reserva${NC}"
    echo "=============================="
    
    read -p "ID de la reserva a eliminar: " reservation_id
    
    # Obtener token de admin
    TOKEN=$(get_admin_token)
    
    read -p "¿Confirmar eliminación? (y/n): " confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        response=$(curl -s -X DELETE "$RESERVATIONS_API/reservas/$reservation_id" \
            -H "Authorization: Bearer $TOKEN")
        
        if echo "$response" | grep -q "error"; then
            echo -e "${RED}❌ Error al eliminar reserva${NC}"
        else
            echo -e "${GREEN}✅ Reserva eliminada exitosamente${NC}"
        fi
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

list_reservations() {
    echo -e "\n${CYAN}📅 Lista de Reservas${NC}"
    echo "=============================="
    
    docker exec mongo mongosh reservasdb --quiet --eval \
        "db.reservas.find({}, {_id: 1, actividad: 1, schedule: 1, users_id: 1, status: 1}).forEach(r => print(r._id + ' | Activity: ' + r.actividad + ' | Schedule: ' + (r.schedule || 'N/A') + ' | Users: ' + r.users_id.length + ' | Status: ' + r.status))"
}

# ===== OPERACIONES MASIVAS =====

clear_all_users() {
    echo -e "\n${RED}⚠️  ADVERTENCIA: Eliminar TODOS los usuarios de prueba${NC}"
    echo "=============================="
    echo -e "${YELLOW}Esto eliminará todos los usuarios EXCEPTO root y admin${NC}"
    read -p "¿Estás seguro? Escribe 'CONFIRMAR': " confirm
    
    if [ "$confirm" = "CONFIRMAR" ]; then
        docker exec arq-soft-ii-mysql-1 mysql -u root -prootpassword -e \
            "DELETE FROM users_db.users WHERE email NOT IN ('root@example.com', 'admin@example.com');" \
            2>/dev/null
        echo -e "${GREEN}✅ Usuarios de prueba eliminados${NC}"
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

clear_all_activities() {
    echo -e "\n${RED}⚠️  ADVERTENCIA: Eliminar TODAS las actividades${NC}"
    echo "=============================="
    read -p "¿Estás seguro? Escribe 'CONFIRMAR': " confirm
    
    if [ "$confirm" = "CONFIRMAR" ]; then
        docker exec mongo mongosh activitiesdb --quiet --eval \
            "db.activities.deleteMany({})"
        echo -e "${GREEN}✅ Todas las actividades eliminadas${NC}"
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

clear_all_reservations() {
    echo -e "\n${RED}⚠️  ADVERTENCIA: Eliminar TODAS las reservas${NC}"
    echo "=============================="
    read -p "¿Estás seguro? Escribe 'CONFIRMAR': " confirm
    
    if [ "$confirm" = "CONFIRMAR" ]; then
        docker exec mongo mongosh reservasdb --quiet --eval \
            "db.reservas.deleteMany({})"
        echo -e "${GREEN}✅ Todas las reservas eliminadas${NC}"
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

reload_seed_data() {
    echo -e "\n${CYAN}🔄 Recargar Datos de Seed${NC}"
    echo "=============================="
    read -p "¿Recargar todos los datos de prueba? (y/n): " confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        bash "$(dirname "$0")/seed-data.sh"
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
}

# ===== ESTADÍSTICAS =====

show_stats() {
    echo -e "\n${CYAN}📊 Estadísticas del Sistema${NC}"
    echo "=============================="
    
    echo -e "\n${BLUE}👥 Usuarios por Rol:${NC}"
    docker exec arq-soft-ii-mysql-1 mysql -u root -prootpassword -e \
        "SELECT role, COUNT(*) as total FROM users_db.users GROUP BY role ORDER BY FIELD(role, 'root', 'admin', 'user');" \
        2>/dev/null | column -t
    
    echo -e "\n${BLUE}🏃 Actividades:${NC}"
    total_activities=$(docker exec mongo mongosh activitiesdb --quiet --eval "db.activities.countDocuments()")
    active_activities=$(docker exec mongo mongosh activitiesdb --quiet --eval "db.activities.countDocuments({is_active: true})")
    echo "  Total: $total_activities"
    echo "  Activas: $active_activities"
    echo "  Inactivas: $((total_activities - active_activities))"
    
    echo -e "\n${BLUE}📅 Reservas por Estado:${NC}"
    docker exec mongo mongosh reservasdb --quiet --eval \
        "db.reservas.aggregate([{\\$group: {_id: '\\$status', total: {\\$sum: 1}}}]).forEach(r => print(r._id + ': ' + r.total))"
    
    echo ""
}

# ===== MENÚ PRINCIPAL =====

show_menu() {
    clear
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║   🎮 Gestión de Datos de Prueba       ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${GREEN}👥 USUARIOS:${NC}"
    echo "  1) Crear usuario"
    echo "  2) Eliminar usuario"
    echo "  3) Listar usuarios"
    echo ""
    echo -e "${GREEN}🏃 ACTIVIDADES:${NC}"
    echo "  4) Crear actividad"
    echo "  5) Eliminar actividad"
    echo "  6) Listar actividades"
    echo ""
    echo -e "${GREEN}📅 RESERVAS:${NC}"
    echo "  7) Eliminar reserva"
    echo "  8) Listar reservas"
    echo ""
    echo -e "${YELLOW}🔧 OPERACIONES MASIVAS:${NC}"
    echo "  9) Eliminar todos los usuarios de prueba"
    echo " 10) Eliminar todas las actividades"
    echo " 11) Eliminar todas las reservas"
    echo " 12) Recargar seed data completo"
    echo ""
    echo -e "${CYAN}📊 INFORMACIÓN:${NC}"
    echo " 13) Mostrar estadísticas"
    echo ""
    echo -e "${RED} 0) Salir${NC}"
    echo ""
}

# ===== LOOP PRINCIPAL =====

main() {
    # Verificar que Docker esté corriendo
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}❌ Error: Docker no está corriendo${NC}"
        exit 1
    fi
    
    # Verificar que los contenedores estén corriendo
    if ! docker-compose ps | grep -q "Up"; then
        echo -e "${RED}❌ Error: Los contenedores no están corriendo${NC}"
        echo -e "${YELLOW}Ejecuta: docker-compose up -d${NC}"
        exit 1
    fi
    
    while true; do
        show_menu
        read -p "Selecciona una opción: " option
        
        case $option in
            1) create_user ;;
            2) delete_user ;;
            3) list_users ;;
            4) create_activity ;;
            5) delete_activity ;;
            6) list_activities ;;
            7) delete_reservation ;;
            8) list_reservations ;;
            9) clear_all_users ;;
            10) clear_all_activities ;;
            11) clear_all_reservations ;;
            12) reload_seed_data ;;
            13) show_stats ;;
            0) 
                echo -e "\n${GREEN}👋 ¡Hasta luego!${NC}\n"
                exit 0
                ;;
            *) 
                echo -e "${RED}❌ Opción inválida${NC}"
                ;;
        esac
        
        echo ""
        read -p "Presiona Enter para continuar..."
    done
}

# Ejecutar menú principal
main

