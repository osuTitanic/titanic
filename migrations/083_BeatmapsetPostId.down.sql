BEGIN;

LOCK TABLE forum_posts, beatmapsets IN SHARE ROW EXCLUSIVE MODE;

DROP TRIGGER trigger_sync_beatmapset_post_id ON forum_posts;
DROP TRIGGER trigger_set_beatmapset_post_id ON beatmapsets;
DROP FUNCTION sync_beatmapset_post_id();
DROP FUNCTION set_beatmapset_post_id();
ALTER TABLE beatmapsets DROP COLUMN post_id;

COMMIT;
