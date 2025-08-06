-- Add a seperate read-only column to beatmapsets for total playcount
-- This will be used to optimize queries for most played beatmapsets
ALTER TABLE beatmapsets ADD COLUMN total_playcount bigint NOT NULL DEFAULT 0;

CREATE FUNCTION trg_beatmap_insert() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  UPDATE beatmapsets
     SET total_playcount = total_playcount + NEW.playcount
   WHERE id = NEW.set_id;
  RETURN NEW;
END;
$$;

CREATE FUNCTION trg_beatmap_update() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
  delta bigint := NEW.playcount - OLD.playcount;
BEGIN
  UPDATE beatmapsets
     SET total_playcount = total_playcount + delta
   WHERE id = NEW.set_id;
  RETURN NEW;
END;
$$;

CREATE FUNCTION trg_beatmap_delete() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  UPDATE beatmapsets
     SET total_playcount = total_playcount - OLD.playcount
   WHERE id = OLD.set_id;
  RETURN OLD;
END;
$$;

-- Create triggers for beatmap insert, update, and delete
CREATE TRIGGER beatmap_insert_trg
  AFTER INSERT ON beatmaps
  FOR EACH ROW EXECUTE FUNCTION trg_beatmap_insert();

CREATE TRIGGER beatmap_update_trg
  AFTER UPDATE OF playcount ON beatmaps
  FOR EACH ROW EXECUTE FUNCTION trg_beatmap_update();

CREATE TRIGGER beatmap_delete_trg
  AFTER DELETE ON beatmaps
  FOR EACH ROW EXECUTE FUNCTION trg_beatmap_delete();

-- Backfill existing data
UPDATE beatmapsets AS bs
SET total_playcount = sub.sum_playcount
FROM (
  SELECT set_id, SUM(playcount) AS sum_playcount
  FROM beatmaps
  GROUP BY set_id
) AS sub
WHERE bs.id = sub.set_id;

-- Create index for faster queries on total_playcount
CREATE INDEX idx_beatmapsets_total_playcount ON beatmapsets (total_playcount DESC);
CREATE INDEX idx_beatmapsets_total_playcount_asc ON beatmapsets (total_playcount ASC);