-- osume.exe is no longer shipped after b2014016.5cuttingedge
DELETE FROM releases_official_entries WHERE file_id = 3348 AND release_id > (
    SELECT id FROM releases_official WHERE version = 20141016 AND subversion = 0
);

-- Microsoft.Xna.Framework.dll was removed in b20150522cuttingedge
DELETE FROM releases_official_entries WHERE file_id = 1892 AND release_id IN (
    SELECT id FROM releases_official WHERE version >= 20150522 AND stream != 'stable'
);

-- d3dx9_31.dll is not present in non-xna releases (unsure about date)
DELETE FROM releases_official_entries WHERE file_id = 18 AND release_id IN (
    SELECT id FROM releases_official WHERE stream IN ('noxna', 'ce45') OR (version > 20151001 AND stream != 'stable')
);

-- x3daudio1_1.dll is not present in non-xna releases (unsure about date)
DELETE FROM releases_official_entries WHERE file_id = 13 AND release_id IN (
    SELECT id FROM releases_official WHERE stream IN ('noxna', 'ce45') OR (version > 20151001 AND stream != 'stable')
);

-- Lazer stuff should only be in lazer stuff
DELETE FROM releases_official_entries WHERE file_id IN (
    3354, 3606, 3355, 3607, 3356, 3608, 3357, 3609, 3358, 3610, 3359, 3611, 3360, 3612, 3361, 3613, 3362, 3614, 3363,
    3615, 3364, 3616, 3365, 3617, 3367, 3618, 3369, 3619, 3370, 3620, 3371, 3621, 3372, 3622, 3373, 3623
) AND release_id IN (
    SELECT id FROM releases_official WHERE stream != 'lazer'
);