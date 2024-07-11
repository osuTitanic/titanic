ALTER TABLE beatmapsets DROP COLUMN IF EXISTS search;
ALTER TABLE beatmaps DROP COLUMN IF EXISTS search;

ALTER TABLE beatmapsets ADD search tsvector
GENERATED ALWAYS AS (
    setweight(to_tsvector('simple', coalesce(title, '')), 'B') ||
    setweight(to_tsvector('simple', coalesce(title_unicode, '')), 'A') ||
    setweight(to_tsvector('simple', coalesce(artist, '')), 'B') ||
    setweight(to_tsvector('simple', coalesce(artist_unicode, '')), 'A') ||
    setweight(to_tsvector('simple', coalesce(creator, '')), 'B') ||
    setweight(to_tsvector('simple', coalesce(source, '')), 'B') ||
    setweight(to_tsvector('simple', coalesce(tags, '')), 'B') :: tsvector
) STORED;

ALTER TABLE beatmaps ADD search tsvector
GENERATED ALWAYS AS (
    to_tsvector('simple', coalesce(version, '')) :: tsvector
) STORED;

CREATE INDEX beatmapsets_search_idx ON beatmapsets USING gin (search);
CREATE INDEX beatmaps_search_idx ON beatmaps USING gin (search);