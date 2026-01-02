-- Remove trigger & function
DROP TRIGGER IF EXISTS trigger_update_beatmapset_rating ON ratings;
DROP FUNCTION IF EXISTS update_beatmapset_rating_stats();

-- Remove columns from beatmapsets
ALTER TABLE beatmapsets
    DROP COLUMN IF EXISTS rating_average,
    DROP COLUMN IF EXISTS rating_count;
