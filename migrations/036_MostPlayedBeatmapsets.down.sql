DROP INDEX IF EXISTS idx_beatmapsets_total_playcount;
DROP INDEX IF EXISTS idx_beatmapsets_total_playcount_asc;

ALTER TABLE beatmapsets DROP COLUMN total_playcount;

DROP TRIGGER IF EXISTS beatmap_insert_trg ON beatmaps;
DROP TRIGGER IF EXISTS beatmap_update_trg ON beatmaps;
DROP TRIGGER IF EXISTS beatmap_delete_trg ON beatmaps;

DROP FUNCTION IF EXISTS trg_beatmap_insert();
DROP FUNCTION IF EXISTS trg_beatmap_update();
DROP FUNCTION IF EXISTS trg_beatmap_delete();