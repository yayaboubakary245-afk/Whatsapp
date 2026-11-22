-- Sprint 6: Administration, Audit, Signalements et Appels WebRTC
-- Responsable: Mengue

CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reported_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reported_message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    reason VARCHAR(120) NOT NULL,
    details TEXT DEFAULT '',
    status VARCHAR(24) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'resolved', 'dismissed')),
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolution_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(80) NOT NULL,
    target_type VARCHAR(48) NOT NULL,
    target_id VARCHAR(64),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS calls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id UUID REFERENCES conversations(id) ON DELETE SET NULL,
    call_type VARCHAR(16) NOT NULL CHECK (call_type IN ('audio', 'video')),
    status VARCHAR(16) NOT NULL DEFAULT 'ringing' CHECK (status IN ('ringing', 'active', 'ended', 'missed', 'rejected')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    answered_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    duration_seconds INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS call_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    call_id UUID NOT NULL REFERENCES calls(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    direction VARCHAR(16) NOT NULL CHECK (direction IN ('outgoing', 'incoming', 'missed')),
    status VARCHAR(16) NOT NULL,
    duration_seconds INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS call_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    call_id UUID NOT NULL REFERENCES calls(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ,
    is_muted BOOLEAN NOT NULL DEFAULT FALSE,
    is_video_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT uq_call_participant UNIQUE (call_id, user_id)
);

CREATE TABLE IF NOT EXISTS ice_servers (
    id SERIAL PRIMARY KEY,
    urls TEXT NOT NULL,
    username VARCHAR(120),
    credential VARCHAR(120),
    server_type VARCHAR(16) NOT NULL DEFAULT 'stun' CHECK (server_type IN ('stun', 'turn')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    priority INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_calls_caller ON calls(caller_id, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_call_history_user ON call_history(user_id, created_at DESC);
