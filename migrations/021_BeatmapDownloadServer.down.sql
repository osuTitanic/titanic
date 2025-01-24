ALTER TABLE beatmapsets DROP COLUMN IF EXISTS download_server;

DROP TRIGGER IF EXISTS set_download_server ON beatmapsets;
DROP FUNCTION IF EXISTS set_default_download_server();