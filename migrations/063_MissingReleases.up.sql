WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250626, 'cuttingedge', 0, '2025-06-26T06:57:41')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4086, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250627, 'cuttingedge', 0, '2025-06-27T07:07:32')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4087, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250628, 'stable40', 0, '2025-06-28T06:32:38')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4088, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250702, 'cuttingedge', 0, '2025-07-02T06:30:44')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4089, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250702, 'stable40', 1, '2025-07-02T06:33:32')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4090, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250720, 'cuttingedge', 0, '2025-07-19T16:56:23')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4091, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250721, 'cuttingedge', 0, '2025-07-20T17:28:45')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4092, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250814, 'cuttingedge', 0, '2025-08-14T05:09:49')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4093, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20250815, 'stable40', 0, '2025-08-15T05:50:38')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4094, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4040, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251022, 'cuttingedge', 0, '2025-10-22T10:16:22')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4096, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4095, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251022, 'stable40', 1, '2025-10-22T10:30:30')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4097, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4095, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251030, 'cuttingedge', 0, '2025-10-30T05:53:18')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4098, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4095, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251031, 'cuttingedge', 0, '2025-10-30T23:10:10')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4099, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4095, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251102, 'cuttingedge', 0, '2025-11-01T22:40:36')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4101, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4100, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251102, 'stable40', 1, '2025-11-01T22:46:27')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4102, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4100, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251127, 'cuttingedge', 0, '2025-11-27T09:17:22')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4104, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4100, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251128, 'cuttingedge', 0, '2025-11-28T12:01:57')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4105, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4100, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251128, 'stable40', 1, '2025-11-28T12:48:38')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4106, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4100, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251221, 'cuttingedge', 0, '2025-12-21T14:32:13')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4108, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4107, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20251222, 'stable40', 0, '2025-12-22T03:04:50')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4109, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4107, 3891);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20260101, 'cuttingedge', 0, '2026-01-01T07:19:02')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4110, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4107, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20260101, 'stable40', 1, '2026-01-01T07:21:39')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4111, 4067, 3470, 3460, 3468, 1753, 3465, 3466, 4107, 4078);

WITH inserted_release AS (
    INSERT INTO releases_official (version, stream, subversion, created_at)
    VALUES (20260116, 'cuttingedge', 0, '2026-01-16T09:47:49')
    RETURNING id
)
INSERT INTO releases_official_entries (release_id, file_id)
SELECT inserted_release.id, releases_official_files.id
FROM inserted_release, releases_official_files
WHERE releases_official_files.file_version IN (3461, 3462, 3511, 3379, 3464, 3591, 4113, 4112, 3470, 3460, 3468, 1753, 3465, 3466, 4107, 4114);