-- Sprint 3: Contacts et Messagerie privée temps réel
-- Responsable: Mengue

CREATE TABLE IF NOT EXISTS contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nickname VARCHAR(80),
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_owner_contact UNIQUE (owner_id, contact_user_id),
    CONSTRAINT chk_no_self_contact CHECK (owner_id <> contact_user_id)
);

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(16) NOT NULL CHECK (type IN ('private', 'group')),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    last_message_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS conversation_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_read_at TIMESTAMPTZ,
    CONSTRAINT uq_conversation_member UNIQUE (conversation_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reply_to_message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    message_type VARCHAR(16) NOT NULL DEFAULT 'text' CHECK (message_type IN ('text', 'file', 'image', 'video', 'system')),
    content TEXT NOT NULL DEFAULT '',
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS message_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(16) NOT NULL DEFAULT 'sent' CHECK (status IN ('sent', 'delivered', 'read')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_message_user_status UNIQUE (message_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_contacts_owner ON contacts(owner_id);
CREATE INDEX IF NOT EXISTS idx_conv_members_user ON conversation_members(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_conv_created ON messages(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_message_status_msg ON message_status(message_id);
