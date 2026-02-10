-- Add statistic columns
ALTER TABLE forums ADD COLUMN topic_count int NOT NULL DEFAULT 0;
ALTER TABLE forums ADD COLUMN post_count int NOT NULL DEFAULT 0;
ALTER TABLE forum_topics ADD COLUMN post_count int NOT NULL DEFAULT 0;

-- Initialize existing statistics
UPDATE forums f SET topic_count = (
    SELECT COUNT(*) FROM forum_topics t WHERE t.forum_id = f.id
);
UPDATE forums f SET post_count = (
    SELECT COUNT(*) FROM forum_posts p WHERE p.forum_id = f.id
);
UPDATE forum_topics t SET post_count = (
    SELECT COUNT(*) FROM forum_posts p WHERE p.topic_id = t.id
);

-- For topic changes
CREATE OR REPLACE FUNCTION update_forum_topic_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE forums SET topic_count = topic_count + 1 WHERE id = NEW.forum_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE forums SET topic_count = topic_count - 1 WHERE id = OLD.forum_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.forum_id != NEW.forum_id THEN
        UPDATE forums SET topic_count = topic_count - 1 WHERE id = OLD.forum_id;
        UPDATE forums SET topic_count = topic_count + 1 WHERE id = NEW.forum_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- For post changes
CREATE OR REPLACE FUNCTION update_post_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE forum_topics SET post_count = post_count + 1 WHERE id = NEW.topic_id;
        UPDATE forums SET post_count = post_count + 1 WHERE id = NEW.forum_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE forum_topics SET post_count = post_count - 1 WHERE id = OLD.topic_id;
        UPDATE forums SET post_count = post_count - 1 WHERE id = OLD.forum_id;
    ELSIF TG_OP = 'UPDATE' THEN
        IF OLD.topic_id != NEW.topic_id THEN
            UPDATE forum_topics SET post_count = post_count - 1 WHERE id = OLD.topic_id;
            UPDATE forum_topics SET post_count = post_count + 1 WHERE id = NEW.topic_id;
        END IF;
        IF OLD.forum_id != NEW.forum_id THEN
            UPDATE forums SET post_count = post_count - 1 WHERE id = OLD.forum_id;
            UPDATE forums SET post_count = post_count + 1 WHERE id = NEW.forum_id;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER trigger_forum_topic_count
AFTER INSERT OR UPDATE OR DELETE ON forum_topics
FOR EACH ROW EXECUTE FUNCTION update_forum_topic_count();

CREATE TRIGGER trigger_post_counts
AFTER INSERT OR UPDATE OR DELETE ON forum_posts
FOR EACH ROW EXECUTE FUNCTION update_post_counts();