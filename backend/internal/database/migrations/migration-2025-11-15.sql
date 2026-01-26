-- Migration: Add network column to games table

ALTER TABLE games ADD COLUMN network VARCHAR(64);