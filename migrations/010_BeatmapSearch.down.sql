DROP INDEX IF EXISTS beatmapsets_search_idx;
DROP INDEX IF EXISTS beatmaps_search_idx;

ALTER TABLE beatmapsets DROP COLUMN IF EXISTS search;
ALTER TABLE beatmaps DROP COLUMN IF EXISTS search;