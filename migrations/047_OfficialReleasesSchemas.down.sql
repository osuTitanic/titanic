DROP INDEX IF EXISTS idx_releases_releases_files_file_hash;
DROP INDEX IF EXISTS idx_releases_changelog_created;
DROP INDEX IF EXISTS idx_releases_official_entries_file_id;
DROP INDEX IF EXISTS idx_releases_official_created;
DROP INDEX IF EXISTS idx_releases_official_version;
DROP INDEX IF EXISTS idx_releases_official_id;

DROP TABLE IF EXISTS releases_changelog;
DROP TABLE IF EXISTS releases_files;
DROP TABLE IF EXISTS releases_official;
DROP TABLE IF EXISTS releases_official_entries;