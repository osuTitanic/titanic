-- Add rating columns to beatmapsets
ALTER TABLE beatmapsets
    ADD COLUMN IF NOT EXISTS rating_average real NOT NULL DEFAULT 0.0,
    ADD COLUMN IF NOT EXISTS rating_count int NOT NULL DEFAULT 0;

-- Create function to update beatmapset rating
CREATE OR REPLACE FUNCTION update_beatmapset_rating_stats()
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

    -- Update the beatmapset rating
    UPDATE beatmapsets
    SET rating_average = COALESCE((
            SELECT AVG(rating)::real
            FROM ratings
            WHERE set_id = target_set_id
        ), 0.0),
        rating_count = (
            SELECT COUNT(*)::int
            FROM ratings
            WHERE set_id = target_set_id
        )
    WHERE id = target_set_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for rating changes
DROP TRIGGER IF EXISTS trigger_update_beatmapset_rating ON ratings;
CREATE TRIGGER trigger_update_beatmapset_rating
    AFTER INSERT OR UPDATE OR DELETE ON ratings
    FOR EACH ROW
    EXECUTE FUNCTION update_beatmapset_rating_stats();

-- Initialize existing beatmapset rating statistics
UPDATE beatmapsets b
SET
    rating_average = COALESCE(stats.avg_rating, 0.0),
    rating_count = COALESCE(stats.total_count, 0)
FROM (
    SELECT
        set_id,
        AVG(rating)::real AS avg_rating,
        COUNT(*)::int AS total_count
    FROM ratings
    GROUP BY set_id
) stats
WHERE b.id = stats.set_id;
