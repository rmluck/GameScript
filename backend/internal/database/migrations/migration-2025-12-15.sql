-- Migration: Update scenarios table's created_at and updated_at columns to use TIMESTAMP WITH TIME ZONE
ALTER TABLE scenarios
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE USING updated_at AT TIME ZONE 'UTC';