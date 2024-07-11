DELETE FROM resource_mirrors
WHERE url IN (
    'https://api.osu.direct/osu/{}',
    'https://old.ppy.sh/osu/{}',
    'https://b.ppy.sh/thumb/{}l.jpg',
    'https://b.ppy.sh/thumb/{}.jpg',
    'https://b.ppy.sh/preview/{}.mp3',
    'https://api.nerinyan.moe/d/{}?noVideo=true',
    'https://api.nerinyan.moe/d/{}',
    'https://api.osu.direct/d/{}?noVideo=',
    'https://api.osu.direct/d/{}',
    '/api/beatmaps/osz/{}',
    '/api/beatmaps/osz/{}?noVideo=true',
    '/api/beatmaps/osu/{}',
    '/api/beatmaps/mp3/{}',
    '/api/beatmaps/mt/{}',
    '/api/beatmaps/mt/{}?large=true'
);

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