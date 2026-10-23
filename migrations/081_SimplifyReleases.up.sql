ALTER TABLE releases_titanic
    ADD COLUMN download_url varchar,
    ADD COLUMN checksums varchar[] NOT NULL DEFAULT '{}',
    ADD COLUMN screenshot varchar;

-- downloads -> download
-- screenshots -> screenshot
-- checksums (json) -> checksums (list)
UPDATE releases_titanic
SET download_url = downloads[array_lower(downloads, 1)],
    screenshot = screenshots[array_lower(screenshots, 1)],
    checksums = ARRAY(
        SELECT hash
        FROM jsonb_array_elements(hashes) WITH ORDINALITY AS files(file, file_index)
        CROSS JOIN LATERAL jsonb_array_elements_text(COALESCE(file->'md5', '[]'::jsonb))
            WITH ORDINALITY AS hashes(hash, hash_index)
        WHERE file->>'file' = 'osu!.exe'
        ORDER BY file_index, hash_index
    );

ALTER TABLE releases_titanic
    DROP COLUMN supported,
    DROP COLUMN downloads,
    DROP COLUMN hashes,
    DROP COLUMN screenshots;

ALTER TABLE releases_modding ADD COLUMN screenshot varchar;
UPDATE releases_modding SET screenshot = screenshots[array_lower(screenshots, 1)];
ALTER TABLE releases_modding
    DROP COLUMN hashes,
    DROP COLUMN screenshots,
    DROP COLUMN downloads;
