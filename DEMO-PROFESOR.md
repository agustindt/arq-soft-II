# 🎓 GUÍA DE DEMOSTRACIÓN - Users API

## 📋 Información del Proyecto

- **Proyecto**: Sports Activities Platform - Users API
- **Branch**: `feature/users-api-jwt-auth`
- **Versión**: 2.1.0
- **Tecnologías**: Go, Gin Framework, GORM, MySQL, JWT, Docker

---

## 🚀 Cómo Ejecutar el Proyecto

### Prerrequisitos
- Docker y Docker Compose instalados
- Puertos disponibles: 8081 (Users API), 3306 (MySQL), 3000 (Frontend)

### Pasos para Levantar la Aplicación

```bash
# 1. Clonar el repositorio (si es necesario)
git clone <repo-url>
cd arq-soft-II

# 2. Cambiar a la branch correcta
git checkout feature/users-api-jwt-auth

# 3. Levantar todos los servicios
docker-compose up --build

# 4. Verificar que users-api esté corriendo
curl http://localhost:8081/api/v1/health
```

### Detener la Aplicación

```bash
docker-compose down
```

---

## 🎯 Endpoints Implementados

### 1. **Health Check** ✅
```bash
curl http://localhost:8081/api/v1/health
```

**Respuesta Esperada:**
```json
{
  "status": "ok",
  "message": "Users API is running",
  "service": "users-api"
}
```

---

### 2. **Documentación de la API** 📚
```bash
curl http://localhost:8081/
```

Muestra todos los endpoints disponibles, características y versión.

---

### 3. **Registro de Usuario** 🆕
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "estudiante",
    "email": "estudiante@test.com",
    "password": "password123",
    "first_name": "Juan",
    "last_name": "Pérez"
  }'
```

**Características:**
- ✅ Validación de datos de entrada
- ✅ Hash seguro de contraseña con bcrypt
- ✅ Retorna JWT token automáticamente
- ✅ Verifica que email y username no estén duplicados

---

### 4. **Login** 🔐
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "estudiante@test.com",
    "password": "password123"
  }'
```

**Características:**
- ✅ Validación de credenciales
- ✅ Retorna JWT token con claims (user_id, email, username)
- ✅ Token válido por 24 horas

---

### 5. **Ver Perfil** (Protegido con JWT) 👤
```bash
# Primero hacer login para obtener el token
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"estudiante@test.com","password":"password123"}' \
  | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['token'])")

# Usar el token para ver el perfil
curl -X GET http://localhost:8081/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Características:**
- ✅ Requiere JWT token válido
- ✅ Retorna información completa del perfil
- ✅ Incluye campos extendidos (bio, avatar, ubicación, etc.)

---

### 6. **Actualizar Perfil** ✏️
```bash
curl -X PUT http://localhost:8081/api/v1/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Estudiante de Arquitectura de Software",
    "location": "Buenos Aires, Argentina",
    "phone": "+54 11 1234-5678",
    "fitness_level": "intermediate",
    "sports_interests": "[\"fútbol\", \"running\", \"natación\"]"
  }'
```

---

### 7. **Listar Usuarios Públicos** 👥
```bash
curl http://localhost:8081/api/v1/users
```

Retorna lista de usuarios con información pública (sin datos sensibles).

---

### 8. **Obtener Usuario por ID** 🔍
```bash
curl http://localhost:8081/api/v1/users/1
```

---

### 9. **Cambiar Contraseña** 🔑
```bash
curl -X PUT http://localhost:8081/api/v1/profile/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "password123",
    "new_password": "newpassword456"
  }'
```

---

### 10. **Crear Usuario Root** (Admin) 👑
```bash
curl -X POST http://localhost:8081/api/v1/admin/create-root \
  -H "Content-Type: application/json" \
  -d '{
    "username": "root",
    "email": "root@admin.com",
    "password": "rootpassword123",
    "first_name": "Root",
    "last_name": "Admin",
    "secret_key": "your-super-secret-jwt-key-here"
  }'
```

---

### 11. **Endpoints de Administración** (Requieren rol admin/root)

#### Ver Estadísticas del Sistema
```bash
curl -X GET http://localhost:8081/api/v1/admin/stats \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Listar Todos los Usuarios (con datos completos)
```bash
curl -X GET http://localhost:8081/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Actualizar Rol de Usuario
```bash
curl -X PUT http://localhost:8081/api/v1/admin/users/2/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'
```

---

## 🏗️ Arquitectura de la Aplicación

```
users-api/
├── config/               # Configuración de base de datos
│   └── database.go
├── handlers/             # Controladores HTTP
│   ├── auth.go          # Register, Login, Refresh
│   ├── users.go         # Perfil, Lista de usuarios
│   └── admin.go         # Panel de administración
├── middleware/           # Middleware de autenticación
│   ├── auth.go          # Validación JWT
│   └── role.go          # Control de acceso por rol
├── models/              # Modelos de datos
│   └── user.go          # Modelo User con GORM
├── utils/               # Utilidades
│   ├── jwt.go           # Generación/validación JWT
│   └── validator.go     # Validaciones customizadas
├── go.mod               # Dependencias Go
├── main.go              # Punto de entrada
└── Dockerfile           # Configuración Docker
```

---

## 🔐 Características de Seguridad Implementadas

- ✅ **Contraseñas hasheadas** con bcrypt (cost 10)
- ✅ **JWT Tokens** firmados con HMAC-SHA256
- ✅ **Validación de entrada** en todos los endpoints
- ✅ **Control de acceso basado en roles** (user, moderator, admin, root)
- ✅ **CORS configurado** para permitir frontend
- ✅ **Soft delete** para usuarios (no se eliminan realmente de la DB)
- ✅ **Timestamps automáticos** (created_at, updated_at)
- ✅ **Índices en base de datos** para optimización

---

## 🗄️ Base de Datos

### Tabla: `users`

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    
    -- Perfil extendido
    avatar_url VARCHAR(500),
    bio TEXT,
    phone VARCHAR(20),
    birth_date DATE,
    location VARCHAR(100),
    gender ENUM('male', 'female', 'other'),
    height DECIMAL(5,2),
    weight DECIMAL(5,2),
    sports_interests JSON,
    fitness_level ENUM('beginner', 'intermediate', 'advanced'),
    social_links JSON,
    
    -- Control de acceso
    role ENUM('user', 'moderator', 'admin', 'root') DEFAULT 'user',
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Índices
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_role (role)
);
```

---

## 📊 Logs de la Aplicación

Para ver los logs en tiempo real:

```bash
# Ver logs de users-api
docker logs -f arq-soft-ii-users-api-1

# Ver logs de todos los servicios
docker-compose logs -f
```

---

## ✅ Lista de Verificación para el Profesor

### Funcionalidades Básicas
- [ ] La API levanta correctamente con `docker-compose up`
- [ ] El health check responde en `/api/v1/health`
- [ ] Se puede registrar un nuevo usuario
- [ ] Se puede hacer login y obtener JWT token
- [ ] El token JWT permite acceder a endpoints protegidos

### Autenticación y Seguridad
- [ ] Las contraseñas se guardan hasheadas (bcrypt)
- [ ] Los tokens JWT contienen claims correctos
- [ ] Los endpoints protegidos rechazan requests sin token
- [ ] Los endpoints de admin rechazan usuarios sin rol admin

### Base de Datos
- [ ] La conexión a MySQL funciona correctamente
- [ ] Las migraciones se aplican automáticamente
- [ ] Los usuarios se guardan con todos los campos
- [ ] Los índices están creados

### Arquitectura
- [ ] El código sigue arquitectura en capas (handlers, models, middleware)
- [ ] Separación de responsabilidades clara
- [ ] Uso correcto de GORM para ORM
- [ ] Middleware de autenticación bien implementado

### Documentación
- [ ] El endpoint raíz (`/`) muestra documentación completa
- [ ] Logs informativos y claros
- [ ] Código comentado donde es necesario

---

## 🐛 Troubleshooting

### La API no levanta
```bash
# Verificar que los puertos no estén en uso
lsof -i :8081

# Limpiar contenedores anteriores
docker-compose down -v
docker-compose up --build
```

### Error de conexión a MySQL
```bash
# Verificar que MySQL esté healthy
docker-compose ps

# Ver logs de MySQL
docker logs arq-soft-ii-mysql-1
```

### Token JWT inválido
- Verificar que JWT_SECRET esté configurado en docker-compose.yml
- El token expira después de 24 horas

---

## 📞 Contacto

Para cualquier consulta sobre la implementación, referirse al código fuente en la branch `feature/users-api-jwt-auth`.

---

## 🎯 Resumen de Endpoints

| Método | Endpoint | Autenticación | Descripción |
|--------|----------|---------------|-------------|
| GET | `/api/v1/health` | No | Health check |
| GET | `/` | No | Documentación API |
| POST | `/api/v1/auth/register` | No | Registrar usuario |
| POST | `/api/v1/auth/login` | No | Iniciar sesión |
| POST | `/api/v1/auth/refresh` | No | Renovar token |
| GET | `/api/v1/users` | No | Listar usuarios (público) |
| GET | `/api/v1/users/:id` | No | Ver usuario por ID |
| GET | `/api/v1/profile` | JWT | Ver mi perfil |
| PUT | `/api/v1/profile` | JWT | Actualizar perfil |
| PUT | `/api/v1/profile/password` | JWT | Cambiar contraseña |
| POST | `/api/v1/profile/avatar` | JWT | Subir avatar |
| DELETE | `/api/v1/profile/avatar` | JWT | Eliminar avatar |
| POST | `/api/v1/admin/create-root` | Secret Key | Crear usuario root |
| GET | `/api/v1/admin/users` | Admin | Ver todos los usuarios |
| POST | `/api/v1/admin/users` | Admin | Crear usuario |
| PUT | `/api/v1/admin/users/:id/role` | Admin | Cambiar rol |
| PUT | `/api/v1/admin/users/:id/status` | Admin | Cambiar estado |
| GET | `/api/v1/admin/stats` | Admin | Estadísticas |
| DELETE | `/api/v1/admin/users/:id` | Root | Eliminar usuario |

---

**Total de Endpoints Implementados: 18+**

✅ **La API está 100% funcional y lista para evaluación**

