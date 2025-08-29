ALTER TABLE profile_play_history
ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();

UPDATE profile_play_history
SET created_at = make_date(year, month, 1);

ALTER TABLE profile_play_history
ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE profile_replay_history
ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();

UPDATE profile_replay_history
SET created_at = make_date(year, month, 1);

ALTER TABLE profile_replay_history
ALTER COLUMN created_at SET NOT NULL;

CREATE INDEX idx_replay_history_user_mode_created_at
ON profile_replay_history (user_id, mode, created_at DESC);

CREATE INDEX idx_play_history_user_mode_created_at
ON profile_play_history (user_id, mode, created_at DESC);