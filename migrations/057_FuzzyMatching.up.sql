CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add a combined text column for trgm matching
ALTER TABLE beatmapsets ADD COLUMN IF NOT EXISTS search_text text
GENERATED ALWAYS AS (
    coalesce(title, '') || ' ' ||
    coalesce(title_unicode, '') || ' ' ||
    coalesce(artist, '') || ' ' ||
    coalesce(artist_unicode, '') || ' ' ||
    coalesce(creator, '') || ' ' ||
    coalesce(source, '') || ' ' ||
    coalesce(tags, '')
) STORED;

-- Create trgm index for text column
CREATE INDEX idx_beatmapsets_search_trgm ON beatmapsets USING GIN (search_text gin_trgm_ops);

-- Add trgm index on existing post/topic content
CREATE INDEX idx_forum_topics_title_trgm ON forum_topics USING GIN (title gin_trgm_ops);
CREATE INDEX idx_forum_posts_content_trgm ON forum_posts USING GIN (content gin_trgm_ops);
