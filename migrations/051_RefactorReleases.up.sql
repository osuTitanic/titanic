-- Rename "releases" to "releases_titanic"
ALTER TABLE releases RENAME TO releases_titanic;

-- Create "releases" view for compatibility
CREATE VIEW releases AS
SELECT
    name,
    version,
    description,
    known_bugs,
    supported,
    preview,
    downloads,
    hashes,
    created_at,
    category,
    screenshots
FROM releases_titanic;

-- Rename "releases_files" to "releases_official_files"
ALTER TABLE releases_files RENAME TO releases_official_files;

-- Create "releases_files" view for compatibility
CREATE VIEW releases_files AS
SELECT
    id,
    filename,
    file_version,
    file_hash,
    filesize,
    patch_id,
    url_full,
    url_patch,
    timestamp
FROM releases_official_files;

-- Rename "releases_changelog" to "releases_official_changelog"
ALTER TABLE releases_changelog RENAME TO releases_official_changelog;

-- Create "releases_changelog" view for compatibility
CREATE VIEW releases_changelog AS
SELECT
    id,
    text,
    type,
    branch,
    author,
    area,
    created_at
FROM releases_official_changelog;

-- Add new way of storing modded release entries
-- This will replace the "hashes" field eventually, and
-- should generally be a better way of tracking modded releases
CREATE TABLE releases_modding_entries (
    id serial PRIMARY KEY,
    mod_name varchar NOT NULL REFERENCES releases_modding (name),
    version varchar NOT NULL,
    stream varchar NOT NULL,
    checksum CHARACTER (32) NOT NULL,
    download_url varchar,
    update_url varchar,
    post_id int REFERENCES forum_posts (id),
    created_at timestamp without time zone NOT NULL DEFAULT now()
);