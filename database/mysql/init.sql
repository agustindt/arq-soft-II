-- Inicialización de base de datos MySQL para Users API
CREATE DATABASE IF NOT EXISTS users_db;
USE users_db;

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin', 'root') NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Índices para optimización
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);

-- Usuarios por defecto para desarrollo/testing
-- Contraseña para todos: password

-- Usuario root (máxima jerarquía)
INSERT INTO users (username, email, first_name, last_name, password_hash, role)
VALUES ('root', 'root@example.com', 'Root', 'Administrator', '$2a$10$dyX0fZvCTxYIuntXCbAtO.PMpEkc94lTAF30H7r/Y1H9MTos5wZP2', 'root')
ON DUPLICATE KEY UPDATE username=username;

-- Usuario admin
INSERT INTO users (username, email, first_name, last_name, password_hash, role)
VALUES ('admin', 'admin@example.com', 'Admin', 'User', '$2a$10$dyX0fZvCTxYIuntXCbAtO.PMpEkc94lTAF30H7r/Y1H9MTos5wZP2', 'admin')
ON DUPLICATE KEY UPDATE username=username;

-- Usuario regular para testing
INSERT INTO users (username, email, first_name, last_name, password_hash, role)
VALUES ('user', 'user@example.com', 'Regular', 'User', '$2a$10$dyX0fZvCTxYIuntXCbAtO.PMpEkc94lTAF30H7r/Y1H9MTos5wZP2', 'user')
ON DUPLICATE KEY UPDATE username=username;