-- Drop triggers
DROP TRIGGER IF EXISTS trigger_forum_topic_count ON forum_topics;
DROP TRIGGER IF EXISTS trigger_post_counts ON forum_posts;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_forum_topic_count();
DROP FUNCTION IF EXISTS update_post_counts();

-- Remove statistic columns
ALTER TABLE forums DROP COLUMN IF EXISTS topic_count;
ALTER TABLE forums DROP COLUMN IF EXISTS post_count;
ALTER TABLE forum_topics DROP COLUMN IF EXISTS post_count;