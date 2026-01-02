-- Remove trigger & function
DROP TRIGGER IF EXISTS trigger_update_beatmapset_favourites ON favourites;
DROP FUNCTION IF EXISTS update_beatmapset_favourites_count();

-- Remove column from beatmapsets
ALTER TABLE beatmapsets
    DROP COLUMN IF EXISTS favourites_count;
