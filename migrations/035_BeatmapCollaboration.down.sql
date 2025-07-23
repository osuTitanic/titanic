DROP INDEX IF EXISTS idx_beatmap_collaboration_user_id;
DROP INDEX IF EXISTS idx_beatmap_collaboration_beatmap_id;
DROP INDEX IF EXISTS idx_beatmap_collaboration_requests_user_id;
DROP INDEX IF EXISTS idx_beatmap_collaboration_requests_target_id;
DROP INDEX IF EXISTS idx_beatmap_collaboration_requests_beatmap_id;
DROP INDEX IF EXISTS idx_beatmap_collaboration_blacklist_user_id;

DROP TABLE IF EXISTS beatmap_collaboration_blacklist;
DROP TABLE IF EXISTS beatmap_collaboration_requests;
DROP TABLE IF EXISTS beatmap_collaboration;