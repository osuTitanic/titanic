-- Add favourite_count column to beatmapsets
ALTER TABLE beatmapsets
    ADD COLUMN IF NOT EXISTS favourite_count int NOT NULL DEFAULT 0;

-- Create function to update beatmapset favourites count
CREATE OR REPLACE FUNCTION update_beatmapset_favourite_count()
RETURNS TRIGGER AS $$
DECLARE
    target_set_id int;
BEGIN
    -- Determine which set_id to update based on operation
    IF TG_OP = 'DELETE' THEN
        target_set_id := OLD.set_id;
    ELSE
        target_set_id := NEW.set_id;
    END IF;

    -- Update the beatmapset favourites count
    UPDATE beatmapsets
    SET favourite_count = (
            SELECT COUNT(*)::int
            FROM favourites
            WHERE set_id = target_set_id
        )
    WHERE id = target_set_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for favourites changes
DROP TRIGGER IF EXISTS trigger_update_beatmapset_favourites ON favourites;
CREATE TRIGGER trigger_update_beatmapset_favourites
    AFTER INSERT OR UPDATE OR DELETE ON favourites
    FOR EACH ROW
    EXECUTE FUNCTION update_beatmapset_favourite_count();

-- Initialize existing beatmapset favourites count
UPDATE beatmapsets b
SET favourite_count = COALESCE(stats.total_count, 0)
FROM (
    SELECT
        set_id,
        COUNT(*)::int AS total_count
    FROM favourites
    GROUP BY set_id
) stats
WHERE b.id = stats.set_id;
