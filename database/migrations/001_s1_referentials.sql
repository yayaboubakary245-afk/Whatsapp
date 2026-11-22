-- Sprint 1: Tables de référence académique et RBAC
-- Responsable: Mengue | Projet: WhatsApp Campus (SUP'PTIC - Groupe 12)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE CHECK (name IN ('etudiant', 'enseignant', 'admin')),
    description TEXT NOT NULL,
    permissions JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS programs (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    code VARCHAR(16) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS academic_levels (
    id SERIAL PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,
    label VARCHAR(64) NOT NULL,
    order_index INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_programs_department ON programs(department_id);
