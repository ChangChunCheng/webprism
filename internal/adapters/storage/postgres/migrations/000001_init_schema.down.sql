-- WEBPRISM Schema Rollback
-- This drops all tables created in the up migration

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS health_checks;
DROP TABLE IF EXISTS auth_configs;
DROP TABLE IF EXISTS api_specs;

-- Drop extension if no longer needed
-- DROP EXTENSION IF EXISTS "uuid-ossp";
