CREATE TABLE releases_files (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    file_version INT NOT NULL,
    file_hash CHARACTER (32) NOT NULL,
    filesize INT NOT NULL,
    patch_id TEXT,
    url_full TEXT NOT NULL,
    url_patch TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE releases_changelog (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    type TEXT NOT NULL,
    branch TEXT NOT NULL,
    author TEXT NOT NULL,
    area TEXT,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE releases_official (
    id SERIAL PRIMARY KEY,
    version INT NOT NULL,
    stream TEXT NOT NULL,
    subversion INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE releases_official_entries (
    release_id INT NOT NULL REFERENCES releases_official (id) ON DELETE CASCADE,
    file_id INT NOT NULL REFERENCES releases_files (id) ON DELETE CASCADE,
    PRIMARY KEY (release_id, file_id)
);

CREATE INDEX idx_releases_releases_files_file_hash ON releases_files (file_hash);
CREATE INDEX idx_releases_changelog_created ON releases_changelog (created_at DESC);
CREATE INDEX idx_releases_official_entries_file_id ON releases_official_entries (file_id);
CREATE INDEX idx_releases_official_created ON releases_official (created_at DESC);
CREATE INDEX idx_releases_official_version ON releases_official (version DESC);
CREATE INDEX idx_releases_official_id ON releases_official (id);