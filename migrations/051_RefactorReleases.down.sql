-- Drop the modding entries table
DROP TABLE IF EXISTS releases_modding_entries;

-- Drop compatibility views
DROP VIEW IF EXISTS releases_changelog;
DROP VIEW IF EXISTS releases_files;
DROP VIEW IF EXISTS releases;

-- Restore original table names
ALTER TABLE releases_official_changelog RENAME TO releases_changelog;
ALTER TABLE releases_official_files RENAME TO releases_files;
ALTER TABLE releases_titanic RENAME TO releases;
