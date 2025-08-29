DROP INDEX IF EXISTS idx_replay_history_user_mode_created_at;
DROP INDEX IF EXISTS idx_play_history_user_mode_created_at;

ALTER TABLE profile_play_history
DROP COLUMN created_at;

ALTER TABLE profile_replay_history
DROP COLUMN created_at;