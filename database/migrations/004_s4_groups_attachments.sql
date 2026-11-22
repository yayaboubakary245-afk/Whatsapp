-- Sprint 4: Groupes universitaires (7 types) et Pièces jointes pédagogiques
-- Responsable: Mengue

CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL UNIQUE REFERENCES conversations(id) ON DELETE CASCADE,
    name VARCHAR(120) NOT NULL,
    description TEXT DEFAULT '',
    group_type VARCHAR(24) NOT NULL CHECK (
        group_type IN ('classe', 'filiere', 'niveau', 'departement', 'td', 'projet', 'pedagogique')
    ),
    department_id INT REFERENCES departments(id) ON DELETE SET NULL,
    program_id INT REFERENCES programs(id) ON DELETE SET NULL,
    academic_level_id INT REFERENCES academic_levels(id) ON DELETE SET NULL,
    avatar_url TEXT,
    only_admins_can_post BOOLEAN NOT NULL DEFAULT FALSE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    member_role VARCHAR(16) NOT NULL DEFAULT 'member' CHECK (member_role IN ('member', 'moderator', 'admin')),
    can_post BOOLEAN NOT NULL DEFAULT TRUE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_group_member UNIQUE (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    uploader_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(120) NOT NULL,
    file_category VARCHAR(24) NOT NULL CHECK (file_category IN ('pdf', 'doc', 'presentation', 'image', 'video', 'other')),
    file_size BIGINT NOT NULL CHECK (file_size >= 0),
    storage_path TEXT NOT NULL,
    checksum_sha256 VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_groups_type ON groups(group_type);
CREATE INDEX IF NOT EXISTS idx_group_members_user ON group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_attachments_message ON attachments(message_id);
CREATE INDEX IF NOT EXISTS idx_attachments_uploader ON attachments(uploader_id);
