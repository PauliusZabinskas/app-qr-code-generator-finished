-- WiFi QR Code Generator - Database Schema
-- PostgreSQL 16+
-- ============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- TABLES
-- ============================================================================

-- Table: users
-- Stores user account information for authentication and authorization
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Table: wifi_qr_codes
-- Stores WiFi credentials and generated QR code data
CREATE TABLE wifi_qr_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    ssid VARCHAR(32) NOT NULL,
    encrypted_password BYTEA NULL, -- NULL for open networks (nopass)
    security_type VARCHAR(10) NOT NULL CHECK (security_type IN ('WPA', 'WPA2', 'WEP', 'nopass')),
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    qr_code_data TEXT NOT NULL, -- WiFi QR code string format
    qr_code_image_url VARCHAR(500) NULL, -- Optional: URL to stored image
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,

    -- Foreign key constraint
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Users table indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- WiFi QR codes table indexes
CREATE INDEX idx_wifi_qr_codes_user_id ON wifi_qr_codes(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_wifi_qr_codes_created_at ON wifi_qr_codes(created_at DESC);
CREATE INDEX idx_wifi_qr_codes_user_created ON wifi_qr_codes(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_wifi_qr_codes_deleted_at ON wifi_qr_codes(deleted_at);

-- Full-text search index on SSID (for admin search functionality)
CREATE INDEX idx_wifi_qr_codes_ssid_trgm ON wifi_qr_codes USING gin(ssid gin_trgm_ops);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for wifi_qr_codes table
CREATE TRIGGER update_wifi_qr_codes_updated_at
    BEFORE UPDATE ON wifi_qr_codes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA (Development/Testing)
-- ============================================================================

-- Admin user
-- Email: admin@wifiqr.com
-- Password: Admin@123
-- bcrypt hash with cost 10
INSERT INTO users (id, email, password_hash, role) VALUES
(
    '550e8400-e29b-41d4-a716-446655440000',
    'admin@wifiqr.com',
    '$2a$10$8K1p/a0Dq7HU0Fv0YMcXOumQg5sDNHo9gvj5yxKp7xKm.GZ9DqG3m',
    'admin'
);

-- Regular user
-- Email: user@example.com
-- Password: User@123
INSERT INTO users (id, email, password_hash, role) VALUES
(
    '550e8400-e29b-41d4-a716-446655440001',
    'user@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'user'
);

-- Sample WiFi QR code (Note: encrypted_password is placeholder, actual encryption done by application)
-- SSID: MyHomeWiFi
-- Password: SecurePass123
-- Security: WPA2
INSERT INTO wifi_qr_codes (id, user_id, ssid, encrypted_password, security_type, is_hidden, qr_code_data) VALUES
(
    '650e8400-e29b-41d4-a716-446655440010',
    '550e8400-e29b-41d4-a716-446655440001',
    'MyHomeWiFi',
    NULL, -- Encrypted by application
    'WPA2',
    false,
    'WIFI:T:WPA2;S:MyHomeWiFi;P:SecurePass123;H:false;;'
);

-- ============================================================================
-- VIEWS
-- ============================================================================

-- View: Active users (excluding soft-deleted)
CREATE OR REPLACE VIEW active_users AS
SELECT
    id,
    email,
    role,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL;

-- View: Active QR codes with user information
CREATE OR REPLACE VIEW active_qr_codes_with_users AS
SELECT
    qr.id,
    qr.user_id,
    u.email as user_email,
    u.role as user_role,
    qr.ssid,
    qr.security_type,
    qr.is_hidden,
    qr.qr_code_data,
    qr.created_at,
    qr.updated_at
FROM wifi_qr_codes qr
INNER JOIN users u ON qr.user_id = u.id
WHERE qr.deleted_at IS NULL AND u.deleted_at IS NULL;

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Function to get user QR code count
CREATE OR REPLACE FUNCTION get_user_qr_count(p_user_id UUID)
RETURNS INTEGER AS $$
BEGIN
    RETURN (
        SELECT COUNT(*)
        FROM wifi_qr_codes
        WHERE user_id = p_user_id AND deleted_at IS NULL
    );
END;
$$ LANGUAGE plpgsql;

-- Function to get total system statistics
CREATE OR REPLACE FUNCTION get_system_stats()
RETURNS TABLE(
    total_users BIGINT,
    total_admins BIGINT,
    total_qr_codes BIGINT,
    active_users BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        (SELECT COUNT(*) FROM users WHERE deleted_at IS NULL) as total_users,
        (SELECT COUNT(*) FROM users WHERE role = 'admin' AND deleted_at IS NULL) as total_admins,
        (SELECT COUNT(*) FROM wifi_qr_codes WHERE deleted_at IS NULL) as total_qr_codes,
        (SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at > NOW() - INTERVAL '30 days') as active_users;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- SECURITY
-- ============================================================================

-- Revoke public access
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA public FROM PUBLIC;

-- Create application role (used by backend)
CREATE ROLE wifiqr_app WITH LOGIN PASSWORD 'change_me_in_production';

-- Grant necessary permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON users TO wifiqr_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON wifi_qr_codes TO wifiqr_app;
GRANT SELECT ON active_users TO wifiqr_app;
GRANT SELECT ON active_qr_codes_with_users TO wifiqr_app;
GRANT EXECUTE ON FUNCTION get_user_qr_count(UUID) TO wifiqr_app;
GRANT EXECUTE ON FUNCTION get_system_stats() TO wifiqr_app;

-- Grant sequence usage for UUID generation
GRANT USAGE ON SCHEMA public TO wifiqr_app;

-- ============================================================================
-- MAINTENANCE
-- ============================================================================

-- Query to permanently delete soft-deleted records older than 90 days
-- Run this periodically as a scheduled job
CREATE OR REPLACE FUNCTION cleanup_old_deleted_records()
RETURNS void AS $$
BEGIN
    DELETE FROM wifi_qr_codes
    WHERE deleted_at IS NOT NULL
    AND deleted_at < NOW() - INTERVAL '90 days';

    DELETE FROM users
    WHERE deleted_at IS NOT NULL
    AND deleted_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- ANALYTICS QUERIES (Examples for admin dashboard)
-- ============================================================================

-- Most popular security types
CREATE OR REPLACE VIEW security_type_stats AS
SELECT
    security_type,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
FROM wifi_qr_codes
WHERE deleted_at IS NULL
GROUP BY security_type
ORDER BY count DESC;

-- QR codes created per day (last 30 days)
CREATE OR REPLACE VIEW daily_qr_creation_stats AS
SELECT
    DATE(created_at) as date,
    COUNT(*) as qr_codes_created
FROM wifi_qr_codes
WHERE created_at > NOW() - INTERVAL '30 days'
AND deleted_at IS NULL
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- Top users by QR code count
CREATE OR REPLACE VIEW top_users_by_qr_count AS
SELECT
    u.id,
    u.email,
    COUNT(qr.id) as qr_code_count
FROM users u
LEFT JOIN wifi_qr_codes qr ON u.id = qr.user_id AND qr.deleted_at IS NULL
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.email
ORDER BY qr_code_count DESC
LIMIT 10;

-- ============================================================================
-- BACKUP RECOMMENDATIONS
-- ============================================================================

-- Recommended backup strategy:
-- 1. Full backup daily: pg_dump -Fc wifiqr_db > backup_$(date +%Y%m%d).dump
-- 2. Point-in-time recovery: Enable WAL archiving
-- 3. Retention: Keep daily backups for 30 days, weekly for 3 months
-- 4. Test restores monthly

-- Example backup command:
-- pg_dump -h localhost -U wifiqr -Fc -f wifiqr_backup_$(date +%Y%m%d_%H%M%S).dump wifiqr_db

-- Example restore command:
-- pg_restore -h localhost -U wifiqr -d wifiqr_db -c wifiqr_backup_YYYYMMDD_HHMMSS.dump
