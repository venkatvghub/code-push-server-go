-- Initialize database schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE os_type AS ENUM ('ios', 'android');
CREATE TYPE platform_type AS ENUM ('react_native', 'cordova');

-- Check if the database exists, and create it if it doesn't
DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_database
      WHERE datname = 'codepush'
   ) THEN
      PERFORM dblink_exec('dbname=' || current_database(), 'CREATE DATABASE codepush');
   END IF;
END
$$;

