-- osume.exe is no longer shipped after b2014016.5cuttingedge
DELETE FROM releases_official_entries WHERE file_id = 3348 AND release_id > (
    SELECT id FROM releases_official WHERE version = 20141016 AND subversion = 0
);