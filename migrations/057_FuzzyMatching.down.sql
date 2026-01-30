ALTER TABLE beatmapsets DROP COLUMN search_text;

DROP INDEX IF EXISTS idx_forum_topics_title_trgm;
DROP INDEX IF EXISTS idx_forum_posts_content_trgm;
DROP INDEX IF EXISTS idx_beatmapsets_search_trgm;

DROP EXTENSION IF EXISTS pg_trim;