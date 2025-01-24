-- Column will determine the server to download from in `resource_mirrors` table.
-- This is useful when we want to have custom osz downloads for beatmaps from osu.ppy.sh.
ALTER TABLE beatmapsets ADD COLUMN download_server smallint DEFAULT 0;

-- Set all existing beatmapsets to download from the same server as the main server.
UPDATE beatmapsets SET download_server = server;

-- Create a trigger function to set the download_server value based on the server column
CREATE OR REPLACE FUNCTION set_default_download_server()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.download_server = 0 THEN
    NEW.download_server := NEW.server;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function
CREATE TRIGGER set_download_server
BEFORE INSERT ON beatmapsets
FOR EACH ROW
EXECUTE FUNCTION set_default_download_server();