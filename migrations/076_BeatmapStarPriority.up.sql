BEGIN;

LOCK TABLE beatmap_modding, users, beatmapsets, forum_topics
    IN SHARE ROW EXCLUSIVE MODE;

ALTER TABLE users
    ADD COLUMN kudosu_earned int NOT NULL DEFAULT 0,
    ADD COLUMN kudosu_spent int NOT NULL DEFAULT 0;

ALTER TABLE forum_topics
    ADD COLUMN star_priority int NOT NULL DEFAULT 0;

CREATE TABLE beatmapset_stars
(
    id bigserial NOT NULL PRIMARY KEY,
    set_id int NOT NULL REFERENCES beatmapsets (id) ON DELETE CASCADE,
    user_id int NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kudosu_cost int NOT NULL DEFAULT 1 CHECK (kudosu_cost > 0),
    star_priority int NOT NULL DEFAULT 1 CHECK (star_priority > 0),
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX beatmapset_stars_set_id_created_at_idx
    ON beatmapset_stars (set_id, created_at, id);
CREATE INDEX beatmapset_stars_user_id_idx
    ON beatmapset_stars (user_id);

WITH kudosu_balances AS (
    SELECT
        users.id,
        COALESCE(modding.earned, 0)::int AS earned,
        COALESCE(stars.spent, 0)::int AS spent
    FROM users
    LEFT JOIN (
        SELECT target_id AS user_id, SUM(amount) AS earned
        FROM beatmap_modding
        GROUP BY target_id
    ) AS modding ON modding.user_id = users.id
    LEFT JOIN (
        SELECT user_id, SUM(kudosu_cost) AS spent
        FROM beatmapset_stars
        GROUP BY user_id
    ) AS stars ON stars.user_id = users.id
)
UPDATE users
SET kudosu_earned = kudosu_balances.earned,
    kudosu_spent = kudosu_balances.spent,
    kudosu = kudosu_balances.earned - kudosu_balances.spent
FROM kudosu_balances
WHERE users.id = kudosu_balances.id;

UPDATE forum_topics AS topic
SET star_priority = beatmapset.star_priority
FROM beatmapsets AS beatmapset
WHERE beatmapset.topic_id = topic.id;

CREATE FUNCTION recalculate_user_kudosu(user_id int)
RETURNS void AS $$
    UPDATE users AS target
    SET kudosu_earned = COALESCE((
            SELECT SUM(modding.amount)
            FROM beatmap_modding AS modding
            WHERE modding.target_id = target.id
        ), 0)::int,
        kudosu = COALESCE((
            SELECT SUM(modding.amount)
            FROM beatmap_modding AS modding
            WHERE modding.target_id = target.id
        ), 0)::int - target.kudosu_spent
    WHERE target.id = user_id;
$$ LANGUAGE sql;

CREATE OR REPLACE FUNCTION sync_user_kudosu()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM recalculate_user_kudosu(NEW.target_id);
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM recalculate_user_kudosu(OLD.target_id);
    ELSIF OLD.target_id = NEW.target_id THEN
        PERFORM recalculate_user_kudosu(NEW.target_id);
    ELSE
        PERFORM recalculate_user_kudosu(OLD.target_id);
        PERFORM recalculate_user_kudosu(NEW.target_id);
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION sync_beatmapset_star()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        RAISE EXCEPTION USING
            ERRCODE = 'P0001',
            MESSAGE = 'beatmapset_star_immutable';
    END IF;

    IF TG_OP = 'INSERT' THEN
        UPDATE users
        SET kudosu_spent = kudosu_spent + NEW.kudosu_cost,
            kudosu = kudosu - NEW.kudosu_cost
        WHERE id = NEW.user_id
          AND kudosu >= NEW.kudosu_cost;

        IF NOT FOUND THEN
            RAISE EXCEPTION USING
                ERRCODE = 'P0001',
                MESSAGE = 'insufficient_kudosu';
        END IF;

        UPDATE beatmapsets
        SET star_priority = star_priority + NEW.star_priority
        WHERE id = NEW.set_id;

        RETURN NEW;
    END IF;

    UPDATE users
    SET kudosu_spent = kudosu_spent - OLD.kudosu_cost,
        kudosu = kudosu + OLD.kudosu_cost
    WHERE id = OLD.user_id;

    UPDATE beatmapsets
    SET star_priority = GREATEST(star_priority - OLD.star_priority, 0)
    WHERE id = OLD.set_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_sync_beatmapset_star
    BEFORE INSERT OR UPDATE OR DELETE ON beatmapset_stars
    FOR EACH ROW
    EXECUTE FUNCTION sync_beatmapset_star();

CREATE FUNCTION mirror_beatmapset_star_priority()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE'
       AND OLD.topic_id IS DISTINCT FROM NEW.topic_id
       AND OLD.topic_id IS NOT NULL THEN
        UPDATE forum_topics
        SET star_priority = 0
        WHERE id = OLD.topic_id;
    END IF;

    IF NEW.topic_id IS NOT NULL THEN
        UPDATE forum_topics
        SET star_priority = NEW.star_priority
        WHERE id = NEW.topic_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_mirror_beatmapset_star_priority_insert
    AFTER INSERT ON beatmapsets
    FOR EACH ROW
    EXECUTE FUNCTION mirror_beatmapset_star_priority();

CREATE TRIGGER trigger_mirror_beatmapset_star_priority_update
    AFTER UPDATE OF star_priority, topic_id ON beatmapsets
    FOR EACH ROW
    EXECUTE FUNCTION mirror_beatmapset_star_priority();

COMMIT;
