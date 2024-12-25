DROP INDEX IF EXISTS idx_beatmap_status_score;
DROP INDEX IF EXISTS idx_beatmap_mode_status_score;
DROP INDEX IF EXISTS idx_score_user_mode_status_score_pp;

ALTER TABLE scores DROP COLUMN "status_score";