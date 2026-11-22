-- Sprint 2: Utilisateurs et Profils
-- Responsable: Mengue (users), Yaya (profiles)

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(160) NOT NULL UNIQUE CHECK (email ~* '^[A-Za-z0-9._%+-]+@ensptic\.cm$'),
    password_hash TEXT NOT NULL,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_code VARCHAR(12),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_suspended BOOLEAN NOT NULL DEFAULT FALSE,
    is_online BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(80) NOT NULL,
    last_name VARCHAR(80) NOT NULL,
    avatar_url TEXT,
    bio TEXT DEFAULT '',
    department_id INT REFERENCES departments(id) ON DELETE SET NULL,
    program_id INT REFERENCES programs(id) ON DELETE SET NULL,
    academic_level_id INT REFERENCES academic_levels(id) ON DELETE SET NULL,
    matricule VARCHAR(32),
    phone VARCHAR(24),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_profiles_user ON profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_profiles_dept_prog ON profiles(department_id, program_id, academic_level_id);
