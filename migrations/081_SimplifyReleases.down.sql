-- Removed support flags, extra URLs, and non-executable hashes cannot be recovered.
ALTER TABLE releases_titanic
    ADD COLUMN supported boolean NOT NULL DEFAULT true,
    ADD COLUMN downloads varchar[] NOT NULL DEFAULT '{}',
    ADD COLUMN hashes jsonb NOT NULL DEFAULT '[]',
    ADD COLUMN screenshots varchar[] NOT NULL DEFAULT '{}';

UPDATE releases_titanic
SET downloads = CASE WHEN download_url IS NULL THEN '{}'::varchar[] ELSE ARRAY[download_url] END,
    screenshots = CASE WHEN screenshot IS NULL THEN '{}'::varchar[] ELSE ARRAY[screenshot] END,
    hashes = CASE WHEN cardinality(checksums) = 0 THEN '[]'::jsonb
        ELSE jsonb_build_array(jsonb_build_object('file', 'osu!.exe', 'md5', to_jsonb(checksums))) END;

ALTER TABLE releases_titanic
    DROP COLUMN download_url,
    DROP COLUMN checksums,
    DROP COLUMN screenshot;

-- Modded hashes and downloads now live in entries; leave the legacy fields empty.
ALTER TABLE releases_modding
    ADD COLUMN hashes jsonb NOT NULL DEFAULT '[]',
    ADD COLUMN screenshots varchar[] NOT NULL DEFAULT '{}',
    ADD COLUMN downloads varchar[] NOT NULL DEFAULT '{}';
UPDATE releases_modding
SET screenshots = CASE WHEN screenshot IS NULL THEN '{}'::varchar[] ELSE ARRAY[screenshot] END;
ALTER TABLE releases_modding DROP COLUMN screenshot;
