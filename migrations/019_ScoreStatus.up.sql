-- The "score_status" update was made to track the personal best for both pp and score.
-- This means that players are not able to loose score or pp anymore (in most cases).
ALTER TABLE scores ADD COLUMN "status_score" SMALLINT NOT NULL DEFAULT -1;
UPDATE scores SET status_score = status;

CREATE INDEX IF NOT EXISTS idx_score_user_mode_status_score_pp ON scores (user_id, mode, status_score, pp DESC);
CREATE INDEX IF NOT EXISTS idx_beatmap_mode_status_score ON scores (beatmap_id, mode, status_score);
CREATE INDEX IF NOT EXISTS idx_beatmap_status_score ON scores (status_score);