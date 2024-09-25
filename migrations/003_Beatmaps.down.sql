DELETE FROM resource_mirrors;

ALTER SEQUENCE beatmapsets_id_seq RESTART WITH 1;
ALTER SEQUENCE beatmaps_id_seq RESTART WITH 1;

DROP INDEX IF EXISTS beatmapsets_id_idx;
DROP INDEX IF EXISTS beatmaps_filename_idx;
DROP INDEX IF EXISTS beatmaps_md5_idx;
DROP INDEX IF EXISTS beatmaps_id_idx;

DROP TABLE IF EXISTS resource_mirrors;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS favourites;
DROP TABLE IF EXISTS plays;
DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS beatmap_modding;
DROP TABLE IF EXISTS beatmap_nominations;
DROP TABLE IF EXISTS beatmaps;
DROP TABLE IF EXISTS beatmapsets;