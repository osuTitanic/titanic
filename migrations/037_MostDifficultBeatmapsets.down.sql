DROP INDEX IF EXISTS idx_beatmaps_set_id_diff;

DROP TRIGGER IF EXISTS recalculate_max_diff_on_insert_or_update ON beatmaps;
DROP TRIGGER IF EXISTS recalculate_max_diff_on_delete ON beatmaps;
DROP FUNCTION IF EXISTS beatmapsets_update_max_diff() CASCADE;

DROP INDEX IF EXISTS idx_beatmapsets_max_diff;
DROP INDEX IF EXISTS idx_beatmapsets_max_diff_asc;

ALTER TABLE beatmapsets DROP COLUMN max_diff;