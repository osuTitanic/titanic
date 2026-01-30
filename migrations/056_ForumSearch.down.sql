DROP INDEX IF EXISTS idx_forum_topics_search;
DROP INDEX IF EXISTS idx_forum_posts_search;

ALTER TABLE forum_topics DROP COLUMN IF EXISTS search_vector;
ALTER TABLE forum_posts DROP COLUMN IF EXISTS search_vector;