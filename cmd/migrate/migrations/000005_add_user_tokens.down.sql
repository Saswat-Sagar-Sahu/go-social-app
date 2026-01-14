-- Rollback for 000005_add_user_tokens.up.sql

DROP INDEX IF EXISTS idx_user_tokens_token;
DROP TABLE IF EXISTS user_tokens;

ALTER TABLE users
DROP COLUMN IF EXISTS activated;
