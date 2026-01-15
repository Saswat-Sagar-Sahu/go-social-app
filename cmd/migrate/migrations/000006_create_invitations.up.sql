-- create invitations table
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inviter_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    invitee_email TEXT NOT NULL,
    message TEXT,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    accepted_at TIMESTAMP WITH TIME ZONE,
    redeemed_by BIGINT REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_invitations_token ON invitations(token);
CREATE INDEX idx_invitations_invitee_email ON invitations(invitee_email);
