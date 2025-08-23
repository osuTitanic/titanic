-- Add a seperate read-only column to beatmapsets for max difficulty
-- This will be used to optimize queries for most difficult beatmapsets
ALTER TABLE beatmapsets ADD COLUMN max_diff real NOT NULL DEFAULT 0;

-- Create index to make following query faster
CREATE INDEX idx_beatmaps_set_id_diff ON beatmaps (set_id, diff DESC);

-- Backfill existing beatmapsets with max difficulty
UPDATE beatmapsets b
SET max_diff = m.max_diff
FROM (SELECT set_id, MAX(diff) AS max_diff FROM beatmaps GROUP BY set_id) m
WHERE b.id = m.set_id;

CREATE OR REPLACE FUNCTION beatmapsets_update_max_diff()
RETURNS TRIGGER AS $$
BEGIN
  -- Recompute max(diff) for the affected set
  UPDATE beatmapsets
    SET max_diff = COALESCE((SELECT MAX(diff) FROM beatmaps WHERE set_id = NEW.set_id), 0)
    WHERE id = NEW.set_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for beatmap insert, update, and delete
CREATE TRIGGER recalculate_max_diff_on_insert_or_update
  AFTER INSERT OR UPDATE ON beatmaps
  FOR EACH ROW
  WHEN (NEW.set_id IS NOT NULL)
  EXECUTE FUNCTION beatmapsets_update_max_diff();

CREATE TRIGGER recalculate_max_diff_on_delete
  AFTER DELETE ON beatmaps
  FOR EACH ROW
  WHEN (OLD.set_id IS NOT NULL)
  EXECUTE FUNCTION beatmapsets_update_max_diff();

-- Create index for faster queries on max_diff
CREATE INDEX idx_beatmapsets_max_diff ON beatmapsets (max_diff DESC);
CREATE INDEX idx_beatmapsets_max_diff_asc ON beatmapsets (max_diff ASC);