DELETE FROM releases_extra;
DELETE FROM releases_modding;

DROP TABLE releases_extra;
DROP TABLE releases_modding;

ALTER TABLE releases ADD COLUMN actions boolean NOT NULL DEFAULT false;
ALTER TABLE releases ADD COLUMN recommended boolean NOT NULL DEFAULT false;

ALTER TABLE releases ADD COLUMN screenshots_migration jsonb NOT NULL DEFAULT '[]';
UPDATE releases SET screenshots_migration = jsonb_build_array(jsonb_build_object('src', screenshots[1], 'width', 0, 'height', 0));
ALTER TABLE releases DROP COLUMN screenshots;
ALTER TABLE releases RENAME screenshots_migration TO screenshots;