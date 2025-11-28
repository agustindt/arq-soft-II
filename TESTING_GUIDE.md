# Guía de Testing - Nuevas Funcionalidades

## Cambios Implementados

1. **Auto-refresh de tokens JWT**
2. **Verificación de expiración de tokens**
3. **Panel de administración de usuarios**
4. **Gestión completa de usuarios (CRUD)**

## Opción 1: Testing Local (Recomendado)

### Paso 1: Levantar solo los servicios backend

Los contenedores de backend ya están corriendo. Verifica que estén activos:

```bash
docker compose ps
```

Deberías ver:
- users-api (puerto 8081) - ✅ Healthy
- activities-api (puerto 8082) - ✅ Healthy
- search-api (puerto 8083) - ✅ Healthy
- reservations-service (puerto 8080) - ✅ Healthy
- MySQL, MongoDB, RabbitMQ, Solr, Memcached - ✅ Running

### Paso 2: Correr el frontend localmente

```bash
cd frontend
npm install  # o pnpm install si usas pnpm
npm start    # Esto levantará el frontend en http://localhost:3000
```

### Paso 3: Crear un usuario administrador

Opción A - Usando la API directamente:

```bash
# 1. Crear usuario root (solo funciona si no existe)
curl -X POST http://localhost:8081/api/v1/admin/create-root \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "username": "admin",
    "password": "Admin123!",
    "first_name": "Admin",
    "last_name": "User",
    "secret_key": "create-root-secret-2024"
  }'
```

Opción B - Registrarse normalmente y actualizar el rol en la base de datos:

```bash
# 1. Regístrate en http://localhost:3000/register
# 2. Conéctate a MySQL y actualiza el rol:

docker exec -it arq-soft-ii-mysql-1 mysql -u root -prootpassword -e \
  "UPDATE users_db.users SET role='admin' WHERE email='tuemail@test.com';"
```

### Paso 4: Tests de Funcionalidades

#### Test 1: Verificación de Expiración de Token

1. Abre las DevTools del navegador (F12) → Console
2. Inicia sesión en http://localhost:3000/login
3. En la consola, ejecuta:

```javascript
// Ver el token decodificado
const token = localStorage.getItem('token');
console.log('Token:', token);

// Importar la función de verificación (solo si estás en el código)
// O verificar manualmente cuándo expira:
const payload = JSON.parse(atob(token.split('.')[1]));
const expiresAt = new Date(payload.exp * 1000);
console.log('Token expira en:', expiresAt);
console.log('Tiempo restante:', Math.floor((payload.exp * 1000 - Date.now()) / 1000 / 60), 'minutos');
```

#### Test 2: Auto-Refresh de Tokens

El token se refresca automáticamente cuando quedan menos de 5 minutos para expirar. Para testearlo:

**Opción A - Modificar temporalmente el código:**

En `frontend/src/utils/jwtUtils.ts`, cambia temporalmente:
```typescript
const FIVE_MINUTES = 5 * 60; // Cambiar a: const FIVE_MINUTES = 23 * 60; // 23 horas
```

Esto hará que el sistema considere que el token expira "pronto" si tiene menos de 23 horas, forzando el refresh.

**Opción B - Monitorear en DevTools:**

1. Abre DevTools → Network
2. Filtra por "refresh"
3. Navega por la aplicación haciendo varias peticiones
4. Deberías ver peticiones automáticas a `/api/v1/auth/refresh` cuando el token esté próximo a expirar

#### Test 3: Panel de Administración

1. Inicia sesión con el usuario admin
2. Ve a http://localhost:3000/admin

**Deberías ver:**
- 📊 Estadísticas de actividades
- 👥 Estadísticas de usuarios
  - Total de usuarios
  - Usuarios activos
  - Nuevos registros (últimos 7 días)
  - Miembros del staff

3. Haz clic en "Manage Users"
4. Ve a http://localhost:3000/admin/users

#### Test 4: Gestión de Usuarios

En http://localhost:3000/admin/users:

**Crear Usuario:**
1. Clic en "Crear Usuario"
2. Completa el formulario:
   - Email: test@example.com
   - Username: testuser
   - Password: Test123!
   - Nombre: Test
   - Apellido: User
   - Rol: user
3. Clic en "Crear"
4. Verifica que aparezca en la lista

**Editar Rol:**
1. Busca el usuario creado
2. Clic en el ícono de editar (lápiz)
3. Cambia el rol (ej: de "user" a "moderator")
4. Clic en "Actualizar"
5. Verifica que el chip de rol se actualice

**Activar/Desactivar:**
1. Clic en el ícono de estado (✓ o ✗)
2. Verifica que el estado cambie

**Eliminar (solo si eres root):**
1. Si tu usuario es "root", verás el ícono de eliminar (🗑️)
2. Clic en eliminar
3. Confirma la eliminación

#### Test 5: Manejo de Errores 401

1. Abre DevTools → Application → Local Storage
2. Modifica el token manualmente (agrega caracteres random)
3. Intenta hacer cualquier acción (ej: ver perfil)
4. El sistema debería:
   - Detectar el token inválido
   - Intentar refrescarlo (fallará porque es inválido)
   - Redirigirte al login automáticamente

#### Test 6: Persistencia de Sesión

1. Inicia sesión
2. Cierra la pestaña del navegador
3. Abre una nueva pestaña en http://localhost:3000
4. Deberías estar automáticamente logueado
5. El token debería refrescarse si está próximo a expirar

## Verificaciones en Backend

### Ver logs del users-api:

```bash
docker logs arq-soft-ii-users-api-1 --tail 50 -f
```

### Probar endpoint de refresh manualmente:

```bash
# 1. Login y obtener token
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"Admin123!"}' \
  | jq -r '.data.token')

echo "Token obtenido: $TOKEN"

# 2. Refrescar token
curl -X POST http://localhost:8081/api/v1/auth/refresh \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq
```

### Verificar estadísticas de admin:

```bash
curl -X GET http://localhost:8081/api/v1/admin/stats \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq
```

### Listar usuarios:

```bash
curl -X GET "http://localhost:8081/api/v1/admin/users?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq
```

## Posibles Problemas y Soluciones

### Frontend no conecta con backend:
- Verifica que los servicios estén corriendo: `docker compose ps`
- Verifica las URLs en `frontend/.env` o las variables de entorno

### Token expira inmediatamente:
- El backend genera tokens con 24 horas de validez
- Verifica que la hora del sistema esté correcta

### No puedes crear usuarios como admin:
- Verifica que tu usuario tenga rol "admin" o "root"
- Revisa los logs del backend para ver errores

### CORS errors:
- El backend ya tiene CORS configurado
- Si persisten, verifica que estés usando http://localhost:3000 (no 127.0.0.1)

## Opción 2: Reconstruir Frontend en Docker (Más complejo)

Actualmente hay un problema con los Dockerfiles que necesita ajustes en la configuración del monorepo.
Si necesitas usar Docker para el frontend, necesitarás:

1. Corregir los Dockerfiles de los backends para usar módulos independientes
2. O, simplificar el go.mod de la raíz para incluir todos los módulos

Por ahora, **recomiendo usar la Opción 1** para testing rápido.
