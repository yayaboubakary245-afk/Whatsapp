-- Sprint 5: Annonces universitaires et Notifications
-- Responsable: Mengue

CREATE TABLE IF NOT EXISTS announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(180) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(32) NOT NULL CHECK (
        category IN ('administrative', 'examen', 'inscription', 'emploi_du_temps', 'evenement', 'pedagogique')
    ),
    target_department_id INT REFERENCES departments(id) ON DELETE SET NULL,
    target_program_id INT REFERENCES programs(id) ON DELETE SET NULL,
    target_level_id INT REFERENCES academic_levels(id) ON DELETE SET NULL,
    attachment_url TEXT,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(32) NOT NULL CHECK (
        type IN ('message', 'group_message', 'mention', 'announcement', 'group_add', 'call_missed', 'system')
    ),
    title VARCHAR(160) NOT NULL,
    body TEXT NOT NULL,
    reference_id UUID,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_announcements_published ON announcements(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_announcements_category ON announcements(category);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id, is_read, created_at DESC);
