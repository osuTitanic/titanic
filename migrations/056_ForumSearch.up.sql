-- Add generated tsvector columns
ALTER TABLE forum_topics ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (to_tsvector('english', coalesce(title, ''))) STORED;

ALTER TABLE forum_posts ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (to_tsvector('english', coalesce(content, ''))) STORED;

-- Create gin indexes
CREATE INDEX idx_forum_topics_search ON forum_topics USING GIN (search_vector);
CREATE INDEX idx_forum_posts_search ON forum_posts USING GIN (search_vector);