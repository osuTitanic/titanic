-- Drop unused columns
ALTER TABLE releases DROP COLUMN actions;
ALTER TABLE releases DROP COLUMN recommended;

-- Modify screenshots to be an array of strings
ALTER TABLE releases ADD COLUMN screenshots_migration varchar[] NOT NULL DEFAULT '{}';
UPDATE releases SET screenshots_migration = ARRAY(SELECT jsonb_array_elements(screenshots)->>'src');
ALTER TABLE releases DROP COLUMN screenshots;
ALTER TABLE releases RENAME screenshots_migration TO screenshots;

-- Releases table specific to modding
CREATE TABLE releases_modding (
    name varchar NOT NULL PRIMARY KEY,
    description text NOT NULL,
    creator_id int REFERENCES users (id),
    topic_id int REFERENCES forum_topics (id),
    client_version int NOT NULL,
    client_extension varchar NOT NULL,
    downloads varchar[] NOT NULL DEFAULT '{}',
    screenshots varchar[] NOT NULL DEFAULT '{}',
    hashes jsonb NOT NULL DEFAULT '[]',
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

-- "Extra" content for osume updater
CREATE TABLE releases_extra (
    name varchar NOT NULL PRIMARY KEY,
    description text NOT NULL,
    download varchar NOT NULL,
    filename varchar NOT NULL
);