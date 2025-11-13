-- Seed data para MySQL (Users Database)
-- Este archivo contiene datos de prueba para development/testing
-- Uso: docker-compose exec mysql mysql -uroot -prootpassword users_db < /seed.sql

USE users_db;

-- Limpiar datos de prueba existentes (SOLO para testing!)
-- DELETE FROM users WHERE email LIKE '%@test.com';

-- Usuarios de prueba
-- Nota: Todas las contraseñas son hasheadas con bcrypt
-- Contraseñas en texto plano para referencia:
--   user1@test.com -> User123!
--   user2@test.com -> User456!
--   moderator@test.com -> Mod789!
--   admin@test.com -> Admin000!

INSERT INTO users (username, email, first_name, last_name, password_hash, user_type) VALUES
-- Usuario normal 1
('johndoe', 'john.doe@test.com', 'John', 'Doe',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'normal'),

-- Usuario normal 2
('janesmth', 'jane.smith@test.com', 'Jane', 'Smith',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'normal'),

-- Usuario normal 3 (activo en deportes)
('sportsfan', 'sports.fan@test.com', 'Carlos', 'Rodriguez',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'normal'),

-- Usuario normal 4
('yogalover', 'yoga.lover@test.com', 'Maria', 'Garcia',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'normal'),

-- Usuario normal 5
('runner', 'runner@test.com', 'Pedro', 'Martinez',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'normal'),

-- Moderador
('moderator', 'moderator@test.com', 'Mod', 'Moderator',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'admin'),

-- Admin de testing
('testadmin', 'testadmin@test.com', 'Test', 'Admin',
 '$2a$10$N9qo8uLOickgx2ZMRZoMye', 'admin')

ON DUPLICATE KEY UPDATE username=username;

-- Verificar datos insertados
SELECT
    username,
    email,
    CONCAT(first_name, ' ', last_name) as full_name,
    user_type,
    created_at
FROM users
WHERE email LIKE '%@test.com'
ORDER BY created_at DESC;

-- Estadísticas
SELECT
    user_type,
    COUNT(*) as total
FROM users
GROUP BY user_type;
