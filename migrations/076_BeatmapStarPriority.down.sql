BEGIN;

LOCK TABLE beatmap_modding, beatmapset_stars, users, beatmapsets, forum_topics
    IN SHARE ROW EXCLUSIVE MODE;

DELETE FROM beatmapset_stars;

DROP TRIGGER IF EXISTS trigger_sync_beatmapset_star ON beatmapset_stars;
DROP FUNCTION IF EXISTS sync_beatmapset_star();
DROP TABLE beatmapset_stars;

DROP TRIGGER IF EXISTS trigger_mirror_beatmapset_star_priority_insert ON beatmapsets;
DROP TRIGGER IF EXISTS trigger_mirror_beatmapset_star_priority_update ON beatmapsets;
DROP FUNCTION IF EXISTS mirror_beatmapset_star_priority();

DROP TRIGGER IF EXISTS trigger_sync_user_kudosu ON beatmap_modding;
DROP FUNCTION IF EXISTS sync_user_kudosu();
DROP FUNCTION IF EXISTS recalculate_user_kudosu(int);

UPDATE users AS target
SET kudosu = COALESCE((
    SELECT SUM(modding.amount)
    FROM beatmap_modding AS modding
    WHERE modding.target_id = target.id
), 0);

ALTER TABLE users
    DROP COLUMN kudosu_spent,
    DROP COLUMN kudosu_earned;

ALTER TABLE forum_topics
    DROP COLUMN star_priority;

CREATE FUNCTION sync_user_kudosu()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users
        SET kudosu = kudosu + NEW.amount
        WHERE id = NEW.target_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users
        SET kudosu = kudosu - OLD.amount
        WHERE id = OLD.target_id;
    ELSIF OLD.target_id = NEW.target_id THEN
        UPDATE users
        SET kudosu = kudosu + NEW.amount - OLD.amount
        WHERE id = NEW.target_id;
    ELSE
        UPDATE users
        SET kudosu = kudosu + CASE
            WHEN id = OLD.target_id THEN -OLD.amount
            ELSE NEW.amount
        END
        WHERE id IN (OLD.target_id, NEW.target_id);
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_sync_user_kudosu
    AFTER INSERT OR UPDATE OR DELETE ON beatmap_modding
    FOR EACH ROW
    EXECUTE FUNCTION sync_user_kudosu();

COMMIT;
