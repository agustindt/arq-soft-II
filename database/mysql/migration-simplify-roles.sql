-- ================================================
-- Migration: Simplify Role System from 5 to 3 Roles
-- ================================================
-- Author: Claude Code
-- Date: 2025-11-28
-- Description: Migrates user roles from 5 levels (user, moderator, admin, super_admin, root)
--              to 3 levels (user, admin, root)
--
-- Migration Strategy:
--   - moderator → user
--   - super_admin → admin
--   - user, admin, root remain unchanged
-- ================================================

USE users_db;

-- Show current role distribution BEFORE migration
SELECT '=== BEFORE MIGRATION ===' AS status;
SELECT role, COUNT(*) as count FROM users GROUP BY role ORDER BY role;

-- Backup current state (optional but recommended)
-- CREATE TABLE users_backup_20251128 AS SELECT * FROM users;

-- Migrate moderator to user
UPDATE users
SET role = 'user'
WHERE role = 'moderator';

-- Migrate super_admin to admin
UPDATE users
SET role = 'admin'
WHERE role = 'super_admin';

-- Show role distribution AFTER migration
SELECT '=== AFTER MIGRATION ===' AS status;
SELECT role, COUNT(*) as count FROM users GROUP BY role ORDER BY role;

-- Verify no orphaned roles exist
SELECT '=== VERIFICATION: Check for invalid roles ===' AS status;
SELECT * FROM users WHERE role NOT IN ('user', 'admin', 'root');

-- Update role column to use ENUM (more strict validation)
-- Note: This will prevent any future inserts with invalid roles
DROP INDEX idx_users_role ON users;

ALTER TABLE users
MODIFY COLUMN role ENUM('user', 'admin', 'root') NOT NULL DEFAULT 'user';

-- Recreate index on role column
CREATE INDEX idx_users_role ON users(role);

-- Show final statistics
SELECT '=== FINAL STATISTICS ===' AS status;
SELECT
    SUM(CASE WHEN role = 'user' THEN 1 ELSE 0 END) as users,
    SUM(CASE WHEN role = 'admin' THEN 1 ELSE 0 END) as admins,
    SUM(CASE WHEN role = 'root' THEN 1 ELSE 0 END) as roots,
    COUNT(*) as total
FROM users;

SELECT '=== MIGRATION COMPLETED SUCCESSFULLY ===' AS status;
