-- Recreate compatibility views
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

CREATE VIEW user_count AS
SELECT
    time,
    (osu_count + irc_count) AS count
FROM user_activity;