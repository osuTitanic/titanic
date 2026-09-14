BEGIN;

-- Prevent kudosu changes between the backfill & trigger creation
LOCK TABLE beatmap_modding IN SHARE ROW EXCLUSIVE MODE;

-- Bring the existing column up to date
UPDATE users AS target
SET kudosu = COALESCE((
    SELECT SUM(modding.amount)
    FROM beatmap_modding AS modding
    WHERE modding.target_id = target.id
), 0);

CREATE OR REPLACE FUNCTION sync_user_kudosu()
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
