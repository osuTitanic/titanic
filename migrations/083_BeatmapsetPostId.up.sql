BEGIN;

-- prevent other services from modifying table while we update it
LOCK TABLE forum_posts, beatmapsets IN SHARE ROW EXCLUSIVE MODE;

ALTER TABLE beatmapsets
    ADD COLUMN post_id bigint REFERENCES forum_posts (id)
        DEFERRABLE INITIALLY DEFERRED; -- check foreign key constraints at COMMIT instead of at the end of each statement

CREATE INDEX beatmapsets_post_id_idx ON beatmapsets (post_id)
    WHERE post_id IS NOT NULL;

UPDATE beatmapsets AS beatmapset
SET post_id = (
    SELECT id FROM forum_posts
    WHERE topic_id = beatmapset.topic_id AND hidden = false
    ORDER BY id LIMIT 1
);

-- this function handles changes to topic_id inside beatmapsets
CREATE FUNCTION set_beatmapset_post_id()
RETURNS trigger AS $$
BEGIN
    PERFORM id FROM forum_topics
    WHERE id = NEW.topic_id
    FOR NO KEY UPDATE; -- lock the topic to prevent concurrent writers from changing the first post

    SELECT id INTO NEW.post_id FROM forum_posts
    WHERE topic_id = NEW.topic_id AND hidden = false
    ORDER BY id LIMIT 1;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_beatmapset_post_id
    BEFORE INSERT OR UPDATE OF topic_id ON beatmapsets
    FOR EACH ROW EXECUTE FUNCTION set_beatmapset_post_id();

-- this function handles changes to forum_posts that may affect the first post of a beatmapset
CREATE FUNCTION sync_beatmapset_post_id()
RETURNS trigger AS $$
DECLARE
    topic_ids int[];
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- post was created
    topic_ids := ARRAY[NEW.topic_id];
    ELSIF TG_OP = 'DELETE' THEN
        -- post was deleted
        topic_ids := ARRAY[OLD.topic_id];
    ELSE
        -- post was updated
        -- check if the topic_id or hidden status changed
        IF OLD.topic_id IS NOT DISTINCT FROM NEW.topic_id
           AND OLD.id IS NOT DISTINCT FROM NEW.id
           AND OLD.hidden IS NOT DISTINCT FROM NEW.hidden THEN
            RETURN NULL;
        END IF;
        topic_ids := ARRAY[OLD.topic_id, NEW.topic_id];
    END IF;

    PERFORM id FROM forum_topics
    WHERE id = ANY(topic_ids)
    ORDER BY id FOR NO KEY UPDATE; -- lock the topics to prevent concurrent writers from changing the first post

    PERFORM id FROM beatmapsets
    WHERE topic_id = ANY(topic_ids)
    ORDER BY id FOR NO KEY UPDATE; -- lock the beatmapsets to prevent concurrent writers from changing the post_id

    UPDATE beatmapsets AS beatmapset
    SET post_id = (
        SELECT id FROM forum_posts
        WHERE topic_id = beatmapset.topic_id AND hidden = false
        ORDER BY id LIMIT 1
    )
    WHERE beatmapset.topic_id = ANY(topic_ids);

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_sync_beatmapset_post_id
    AFTER INSERT OR DELETE OR UPDATE OF id, topic_id, hidden ON forum_posts
    FOR EACH ROW EXECUTE FUNCTION sync_beatmapset_post_id();

COMMIT;
