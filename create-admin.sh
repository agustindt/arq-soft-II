#!/bin/bash

# Script para crear un usuario administrador

echo "🔧 Script de creación de usuario administrador"
echo "=============================================="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Variables
API_URL="http://localhost:8081/api/v1"

echo "Selecciona una opción:"
echo "1) Crear usuario ROOT (requiere secret key)"
echo "2) Crear usuario ADMIN normal"
echo "3) Promover usuario existente a ADMIN"
echo ""
read -p "Opción [1-3]: " option

case $option in
  1)
    echo ""
    echo "${YELLOW}Creando usuario ROOT...${NC}"
    read -p "Email: " email
    read -p "Username: " username
    read -sp "Password: " password
    echo ""
    read -p "Nombre: " first_name
    read -p "Apellido: " last_name

    response=$(curl -s -X POST "$API_URL/admin/create-root" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$email\",
        \"username\": \"$username\",
        \"password\": \"$password\",
        \"first_name\": \"$first_name\",
        \"last_name\": \"$last_name\",
        \"secret_key\": \"create-root-secret-2024\"
      }")

    if echo "$response" | grep -q "error"; then
      echo "${RED}❌ Error: $(echo $response | jq -r '.message')${NC}"
    else
      echo "${GREEN}✅ Usuario ROOT creado exitosamente!${NC}"
      echo ""
      echo "Credenciales:"
      echo "  Email: $email"
      echo "  Username: $username"
      echo ""
      echo "Puedes iniciar sesión en: http://localhost:3000/login"
    fi
    ;;

  2)
    echo ""
    echo "${YELLOW}Nota: Primero necesitas iniciar sesión como admin existente${NC}"
    read -p "Email del admin actual: " admin_email
    read -sp "Password del admin actual: " admin_password
    echo ""

    # Login
    token=$(curl -s -X POST "$API_URL/auth/login" \
      -H "Content-Type: application/json" \
      -d "{\"email\": \"$admin_email\", \"password\": \"$admin_password\"}" \
      | jq -r '.data.token')

    if [ "$token" = "null" ] || [ -z "$token" ]; then
      echo "${RED}❌ Error al iniciar sesión. Verifica las credenciales.${NC}"
      exit 1
    fi

    echo "${GREEN}✅ Sesión iniciada${NC}"
    echo ""

    # Crear nuevo usuario
    read -p "Email del nuevo usuario: " email
    read -p "Username: " username
    read -sp "Password: " password
    echo ""
    read -p "Nombre: " first_name
    read -p "Apellido: " last_name

    response=$(curl -s -X POST "$API_URL/admin/users" \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$email\",
        \"username\": \"$username\",
        \"password\": \"$password\",
        \"first_name\": \"$first_name\",
        \"last_name\": \"$last_name\",
        \"role\": \"admin\"
      }")

    if echo "$response" | grep -q "error"; then
      echo "${RED}❌ Error: $(echo $response | jq -r '.message')${NC}"
    else
      echo "${GREEN}✅ Usuario ADMIN creado exitosamente!${NC}"
      echo ""
      echo "Credenciales:"
      echo "  Email: $email"
      echo "  Username: $username"
      echo "  Rol: admin"
    fi
    ;;

  3)
    echo ""
    read -p "Email del usuario a promover: " user_email

    # Obtener ID del usuario desde MySQL
    user_id=$(docker exec arq-soft-ii-mysql-1 mysql -u root -prootpassword -N -e \
      "SELECT id FROM users_db.users WHERE email='$user_email';")

    if [ -z "$user_id" ]; then
      echo "${RED}❌ Usuario no encontrado${NC}"
      exit 1
    fi

    echo "Usuario encontrado con ID: $user_id"
    echo ""
    echo "Selecciona el nuevo rol:"
    echo "1) user"
    echo "2) moderator"
    echo "3) admin"
    echo "4) super_admin"
    echo "5) root"
    read -p "Opción [1-5]: " role_option

    case $role_option in
      1) new_role="user" ;;
      2) new_role="moderator" ;;
      3) new_role="admin" ;;
      4) new_role="super_admin" ;;
      5) new_role="root" ;;
      *) echo "${RED}Opción inválida${NC}"; exit 1 ;;
    esac

    # Actualizar directamente en la base de datos
    docker exec arq-soft-ii-mysql-1 mysql -u root -prootpassword -e \
      "UPDATE users_db.users SET role='$new_role' WHERE email='$user_email';"

    echo "${GREEN}✅ Usuario promovido a $new_role exitosamente!${NC}"
    echo ""
    echo "El usuario debe cerrar sesión y volver a iniciar para que los cambios tomen efecto."
    ;;

  *)
    echo "${RED}Opción inválida${NC}"
    exit 1
    ;;
esac

echo ""
echo "=============================================="
echo "Para acceder al panel de administración:"
echo "  http://localhost:3000/admin"
echo ""
