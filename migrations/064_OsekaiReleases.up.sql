BEGIN;

-- Release b370
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2005-10-03 06:48:40'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2007-08-28 09:21:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, '797e24743937d67d69f28f2cf5052ee8', 2414360, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_797e24743937d67d69f28f2cf5052ee8', NULL, '2007-08-28 09:21:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '797e24743937d67d69f28f2cf5052ee8');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2008-01-28 03:44:14'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2008-01-03 07:01:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2008-02-09 17:59:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'a9ae80509759d27d893f3f6f47fd3da6', 72414720, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_a9ae80509759d27d893f3f6f47fd3da6', NULL, '2008-07-11 00:17:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a9ae80509759d27d893f3f6f47fd3da6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'b350511e572b2314db7f47ab2699ecf8', 676888, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_b350511e572b2314db7f47ab2699ecf8', NULL, '2008-07-11 00:40:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b350511e572b2314db7f47ab2699ecf8');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '722f8e33e0fd60509fcbc20f87a54902', 434584, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_722f8e33e0fd60509fcbc20f87a54902', NULL, '2008-06-17 23:09:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '722f8e33e0fd60509fcbc20f87a54902');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (370, 'stable', 0, '2008-07-11 00:40:46');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 370 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', '797e24743937d67d69f28f2cf5052ee8', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'a9ae80509759d27d893f3f6f47fd3da6', 'b350511e572b2314db7f47ab2699ecf8', '722f8e33e0fd60509fcbc20f87a54902', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b394
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2008-02-10 00:33:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2008-02-10 00:33:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, '797e24743937d67d69f28f2cf5052ee8', 2414360, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_797e24743937d67d69f28f2cf5052ee8', NULL, '2008-02-10 00:33:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '797e24743937d67d69f28f2cf5052ee8');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2008-02-10 00:33:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2008-02-10 00:33:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'c83235328479ac95b56cedb3a8d4d428', 73060352, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_c83235328479ac95b56cedb3a8d4d428', NULL, '2008-08-14 02:43:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c83235328479ac95b56cedb3a8d4d428');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '03060f9f96a391c7baeb3dd6a5f45158', 851456, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_03060f9f96a391c7baeb3dd6a5f45158', NULL, '2008-08-04 15:53:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '03060f9f96a391c7baeb3dd6a5f45158');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2008-02-10 00:33:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (394, 'stable', 1, '2008-08-04 15:53:50');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 394 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', '797e24743937d67d69f28f2cf5052ee8', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'c83235328479ac95b56cedb3a8d4d428', '03060f9f96a391c7baeb3dd6a5f45158', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b595
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'c58b94a166106806b9eeb94f05a1d573', 102063104, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_c58b94a166106806b9eeb94f05a1d573', NULL, '2009-01-26 13:13:14'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c58b94a166106806b9eeb94f05a1d573');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '118829b412d8651dfd7b903f5701f715', 1079296, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_118829b412d8651dfd7b903f5701f715', NULL, '2009-01-28 17:15:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '118829b412d8651dfd7b903f5701f715');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '0bff45035e66fc78a0c41aee64d3316a', 175104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_0bff45035e66fc78a0c41aee64d3316a', NULL, '2009-01-26 14:30:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '0bff45035e66fc78a0c41aee64d3316a');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (595, 'stable', 0, '2009-01-28 17:15:44');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 595 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', '797e24743937d67d69f28f2cf5052ee8', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'c58b94a166106806b9eeb94f05a1d573', '118829b412d8651dfd7b903f5701f715', '0bff45035e66fc78a0c41aee64d3316a', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b639
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2009-03-09 17:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2009-03-09 17:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-03-09 17:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-03-09 17:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-03-09 17:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'a8de94b46c613efb15dfb22e275560e0', 105773568, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_a8de94b46c613efb15dfb22e275560e0', NULL, '2009-03-26 19:28:42'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a8de94b46c613efb15dfb22e275560e0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'b23ee57bbeb207077c0180027418e3ed', 1125376, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_b23ee57bbeb207077c0180027418e3ed', NULL, '2009-03-26 19:40:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b23ee57bbeb207077c0180027418e3ed');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '0bff45035e66fc78a0c41aee64d3316a', 175104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_0bff45035e66fc78a0c41aee64d3316a', NULL, '2009-01-26 15:30:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '0bff45035e66fc78a0c41aee64d3316a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-03-09 17:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (639, 'stable', 0, '2009-03-26 19:40:10');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 639 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'a8de94b46c613efb15dfb22e275560e0', 'b23ee57bbeb207077c0180027418e3ed', '0bff45035e66fc78a0c41aee64d3316a', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b699
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 06:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 06:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 06:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2009-03-09 18:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2009-03-09 18:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-03-09 18:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 06:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-03-09 18:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'eaa319a1648fd9dc9b54dba08048e8db', 114238976, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_eaa319a1648fd9dc9b54dba08048e8db', NULL, '2009-07-26 19:45:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eaa319a1648fd9dc9b54dba08048e8db');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '98c2a99fdc1e93f22b1b351f74d6d70d', 1260032, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_98c2a99fdc1e93f22b1b351f74d6d70d', NULL, '2009-07-03 09:26:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '98c2a99fdc1e93f22b1b351f74d6d70d');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '0bff45035e66fc78a0c41aee64d3316a', 175104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_0bff45035e66fc78a0c41aee64d3316a', NULL, '2009-01-26 16:30:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '0bff45035e66fc78a0c41aee64d3316a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 06:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-03-09 18:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (699, 'stable', 0, '2009-07-03 09:26:32');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 699 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'eaa319a1648fd9dc9b54dba08048e8db', '98c2a99fdc1e93f22b1b351f74d6d70d', '0bff45035e66fc78a0c41aee64d3316a', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b753
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-03-09 16:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 05:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'eaa319a1648fd9dc9b54dba08048e8db', 114238976, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_eaa319a1648fd9dc9b54dba08048e8db', NULL, '2009-07-26 18:45:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eaa319a1648fd9dc9b54dba08048e8db');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '622f9cc5ecdb072ad1d2bfad939fd4ee', 1277440, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_622f9cc5ecdb072ad1d2bfad939fd4ee', NULL, '2009-07-26 22:48:42'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '622f9cc5ecdb072ad1d2bfad939fd4ee');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '783de9554ce91d0c939d1f6050693e8b', 183296, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_783de9554ce91d0c939d1f6050693e8b', NULL, '2009-06-28 22:21:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '783de9554ce91d0c939d1f6050693e8b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-03-09 16:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (753, 'stable', 3, '2009-07-26 22:48:42');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 753 AND stream = 'stable' AND subversion = 3 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'eaa319a1648fd9dc9b54dba08048e8db', '622f9cc5ecdb072ad1d2bfad939fd4ee', '783de9554ce91d0c939d1f6050693e8b', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b833
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-03-09 16:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 05:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'bae32b9788267cf1a4790dd3756aa824', 114238976, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_bae32b9788267cf1a4790dd3756aa824', NULL, '2009-08-18 04:28:26'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bae32b9788267cf1a4790dd3756aa824');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'e3c077f35bf88fc4a0025508436964b1', 1298944, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_e3c077f35bf88fc4a0025508436964b1', NULL, '2009-08-18 04:28:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'e3c077f35bf88fc4a0025508436964b1');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '783de9554ce91d0c939d1f6050693e8b', 183296, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_783de9554ce91d0c939d1f6050693e8b', NULL, '2009-06-28 22:21:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '783de9554ce91d0c939d1f6050693e8b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-03-09 16:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (833, 'stable', 0, '2009-08-18 04:28:34');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 833 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'bae32b9788267cf1a4790dd3756aa824', 'e3c077f35bf88fc4a0025508436964b1', '783de9554ce91d0c939d1f6050693e8b', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b904
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'bbfc7d855252b0211875769bbf667bcd', 96320, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_bbfc7d855252b0211875769bbf667bcd', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbfc7d855252b0211875769bbf667bcd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, 'f9ffe0a23a32b79653e31330764c4231', 26712, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_f9ffe0a23a32b79653e31330764c4231', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f9ffe0a23a32b79653e31330764c4231');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-03-09 16:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 05:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-03-09 16:18:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '9aa2d1bd59ceb4a534db6bbc4ec47b3b', 114180608, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_9aa2d1bd59ceb4a534db6bbc4ec47b3b', NULL, '2009-09-07 22:25:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9aa2d1bd59ceb4a534db6bbc4ec47b3b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '2e7ad1a6986cbcdb9b676a9fb5ea090c', 1310208, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_2e7ad1a6986cbcdb9b676a9fb5ea090c', NULL, '2009-09-07 22:25:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7ad1a6986cbcdb9b676a9fb5ea090c');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '571cdd3d64b11b1363d390762c561120', 200192, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_571cdd3d64b11b1363d390762c561120', NULL, '2009-08-23 19:21:40'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '571cdd3d64b11b1363d390762c561120');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 05:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-03-09 16:18:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (904, 'stable', 0, '2009-09-07 22:25:50');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 904 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'bbfc7d855252b0211875769bbf667bcd', 'f9ffe0a23a32b79653e31330764c4231', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '9aa2d1bd59ceb4a534db6bbc4ec47b3b', '2e7ad1a6986cbcdb9b676a9fb5ea090c', '571cdd3d64b11b1363d390762c561120', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1077
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '883d1226d91ccd400fa7f27b7a1fdf25', 115918512, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_883d1226d91ccd400fa7f27b7a1fdf25', NULL, '2009-10-25 06:49:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '883d1226d91ccd400fa7f27b7a1fdf25');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '9ceb72400e52f8c117269cc29ea914cc', 1401008, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_9ceb72400e52f8c117269cc29ea914cc', NULL, '2009-10-25 07:13:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9ceb72400e52f8c117269cc29ea914cc');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '571cdd3d64b11b1363d390762c561120', 200192, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_571cdd3d64b11b1363d390762c561120', NULL, '2009-08-24 01:21:40'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '571cdd3d64b11b1363d390762c561120');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1077, 'stable', 1, '2009-10-25 07:13:56');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1077 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '883d1226d91ccd400fa7f27b7a1fdf25', '9ceb72400e52f8c117269cc29ea914cc', '571cdd3d64b11b1363d390762c561120', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1122
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '32dbea6172d262ece4802f0653135eee', 116111024, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_32dbea6172d262ece4802f0653135eee', NULL, '2009-11-02 05:21:28'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '32dbea6172d262ece4802f0653135eee');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '41a64e750847105e406e0093c7bba5fc', 1419440, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_41a64e750847105e406e0093c7bba5fc', NULL, '2009-11-02 22:31:02'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '41a64e750847105e406e0093c7bba5fc');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'aa8d96171ae6eea1e628cda540901594', 200368, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_aa8d96171ae6eea1e628cda540901594', NULL, '2009-11-02 22:30:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'aa8d96171ae6eea1e628cda540901594');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1122, 'stable', 0, '2009-11-02 22:31:02');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1122 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '32dbea6172d262ece4802f0653135eee', '41a64e750847105e406e0093c7bba5fc', 'aa8d96171ae6eea1e628cda540901594', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1218
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '9940baa85fb5ecead9e23e55c395cc8b', 113964720, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_9940baa85fb5ecead9e23e55c395cc8b', NULL, '2010-01-19 07:59:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9940baa85fb5ecead9e23e55c395cc8b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'be4a5fbfbca3e46ec0c957daf9100116', 1428144, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_be4a5fbfbca3e46ec0c957daf9100116', NULL, '2010-01-19 07:59:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'be4a5fbfbca3e46ec0c957daf9100116');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1218, 'stable', 0, '2010-01-19 07:59:50');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1218 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '9940baa85fb5ecead9e23e55c395cc8b', 'be4a5fbfbca3e46ec0c957daf9100116', '6be82ef9871fe708e54e509b656e46fe', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1553
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'b65231d094c3b6ceb5e043c1cc706783', 115154096, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_b65231d094c3b6ceb5e043c1cc706783', NULL, '2010-04-02 06:04:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b65231d094c3b6ceb5e043c1cc706783');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '431e5ed0a09bab73f431776296c8252f', 1501872, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_431e5ed0a09bab73f431776296c8252f', NULL, '2010-04-02 09:16:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '431e5ed0a09bab73f431776296c8252f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '6be82ef9871fe708e54e509b656e46fe', 254128, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_6be82ef9871fe708e54e509b656e46fe', NULL, '2010-01-20 03:44:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '6be82ef9871fe708e54e509b656e46fe');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1553, 'stable', 0, '2010-04-02 09:16:10');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1553 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'b65231d094c3b6ceb5e043c1cc706783', '431e5ed0a09bab73f431776296c8252f', '6be82ef9871fe708e54e509b656e46fe', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1596
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'b65231d094c3b6ceb5e043c1cc706783', 115154096, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_b65231d094c3b6ceb5e043c1cc706783', NULL, '2010-04-02 06:04:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b65231d094c3b6ceb5e043c1cc706783');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'd4ffc6a66eb7b46687efac49f2464cbb', 1511600, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_d4ffc6a66eb7b46687efac49f2464cbb', NULL, '2010-04-22 09:07:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd4ffc6a66eb7b46687efac49f2464cbb');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1596, 'stable', 0, '2010-04-22 09:07:54');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1596 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'b65231d094c3b6ceb5e043c1cc706783', 'd4ffc6a66eb7b46687efac49f2464cbb', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1672
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '73847d12efe9c051fe1baec5484cb355', 116242152, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_73847d12efe9c051fe1baec5484cb355', NULL, '2010-10-10 17:30:12'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '73847d12efe9c051fe1baec5484cb355');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'fc55d1acd0980b764a3b8741598611e1', 1546984, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_fc55d1acd0980b764a3b8741598611e1', NULL, '2010-10-10 17:30:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'fc55d1acd0980b764a3b8741598611e1');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '6be82ef9871fe708e54e509b656e46fe', 254128, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_6be82ef9871fe708e54e509b656e46fe', NULL, '2010-01-20 03:44:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '6be82ef9871fe708e54e509b656e46fe');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1672, 'stable', 0, '2010-10-10 17:30:18');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1672 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '73847d12efe9c051fe1baec5484cb355', 'fc55d1acd0980b764a3b8741598611e1', '6be82ef9871fe708e54e509b656e46fe', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1704
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '002f5ad5fa3c2a3f218225a4597f9365', 119223016, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_002f5ad5fa3c2a3f218225a4597f9365', NULL, '2010-11-25 14:47:42'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '002f5ad5fa3c2a3f218225a4597f9365');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'eec9475d17300cebfdce046fb21ad884', 1693928, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_eec9475d17300cebfdce046fb21ad884', NULL, '2010-12-24 14:42:26'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eec9475d17300cebfdce046fb21ad884');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb4ad0730b45b0702f3c8ce080730fea', 254184, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb4ad0730b45b0702f3c8ce080730fea', NULL, '2010-10-13 00:30:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb4ad0730b45b0702f3c8ce080730fea');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1704, 'stable', 0, '2010-12-24 14:42:26');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1704 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '002f5ad5fa3c2a3f218225a4597f9365', 'eec9475d17300cebfdce046fb21ad884', 'bb4ad0730b45b0702f3c8ce080730fea', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1807
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3b9a06bbf45a430f587047f26e6a512b', 121463528, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3b9a06bbf45a430f587047f26e6a512b', NULL, '2011-04-01 16:15:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3b9a06bbf45a430f587047f26e6a512b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '3ee965ff3ab6b9593818efcbcb6666d2', 1849576, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_3ee965ff3ab6b9593818efcbcb6666d2', NULL, '2011-04-05 10:18:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3ee965ff3ab6b9593818efcbcb6666d2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2011-08-27 13:31:12'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1807, 'stable', 0, '2011-04-05 10:18:10');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1807 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3b9a06bbf45a430f587047f26e6a512b', '3ee965ff3ab6b9593818efcbcb6666d2', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1814
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 20:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3b9a06bbf45a430f587047f26e6a512b', 121463528, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3b9a06bbf45a430f587047f26e6a512b', NULL, '2011-03-29 03:53:26'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3b9a06bbf45a430f587047f26e6a512b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '565da88e6cefecb2824b755b9e30f742', 1851112, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_565da88e6cefecb2824b755b9e30f742', NULL, '2011-04-17 20:58:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '565da88e6cefecb2824b755b9e30f742');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '2e78eca5e90defdb4b10ce048900eedf', 296680, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_2e78eca5e90defdb4b10ce048900eedf', NULL, '2011-03-31 04:26:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e78eca5e90defdb4b10ce048900eedf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1814, 'stable', 0, '2011-04-17 20:58:10');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1814 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3b9a06bbf45a430f587047f26e6a512b', '565da88e6cefecb2824b755b9e30f742', '2e78eca5e90defdb4b10ce048900eedf', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1815
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3b9a06bbf45a430f587047f26e6a512b', 121463528, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3b9a06bbf45a430f587047f26e6a512b', NULL, '2011-03-28 15:57:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3b9a06bbf45a430f587047f26e6a512b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '44e088f5b8f343d6d9ecfe66c1bf5edf', 1851624, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_44e088f5b8f343d6d9ecfe66c1bf5edf', NULL, '2011-08-09 20:37:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '44e088f5b8f343d6d9ecfe66c1bf5edf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1815, 'stable', 0, '2011-08-09 20:37:32');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1815 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3b9a06bbf45a430f587047f26e6a512b', '44e088f5b8f343d6d9ecfe66c1bf5edf', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1818
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 20:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '9c468e4f3d4ff9e188f07b88e3f74350', 122149608, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_9c468e4f3d4ff9e188f07b88e3f74350', NULL, '2012-02-28 08:54:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9c468e4f3d4ff9e188f07b88e3f74350');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '230682d73d93769494d6b4c42e6e586b', 1854184, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_230682d73d93769494d6b4c42e6e586b', NULL, '2012-03-13 12:00:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '230682d73d93769494d6b4c42e6e586b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2011-08-20 09:30:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1818, 'stable', 0, '2012-03-13 12:00:56');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1818 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '9c468e4f3d4ff9e188f07b88e3f74350', '230682d73d93769494d6b4c42e6e586b', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1821
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '9c468e4f3d4ff9e188f07b88e3f74350', 122149608, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_9c468e4f3d4ff9e188f07b88e3f74350', NULL, '2012-02-27 23:54:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9c468e4f3d4ff9e188f07b88e3f74350');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '0e9558ad9eb4fda87aa147bb920f5fe7', 1865728, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_0e9558ad9eb4fda87aa147bb920f5fe7', NULL, '2012-03-29 10:29:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '0e9558ad9eb4fda87aa147bb920f5fe7');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osuframework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osuframework.dll/2493953a68a470cff9385597df3bb3a2', NULL, '2012-05-04 18:36:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2011-08-20 00:30:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1821, 'stable', 0, '2012-03-29 10:29:06');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1821 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '9c468e4f3d4ff9e188f07b88e3f74350', '0e9558ad9eb4fda87aa147bb920f5fe7', '2493953a68a470cff9385597df3bb3a2', 'b4c3f19676e345fae5d89c5332741725', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b1844
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3b9a06bbf45a430f587047f26e6a512b', 121463528, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3b9a06bbf45a430f587047f26e6a512b', NULL, '2011-04-01 14:15:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3b9a06bbf45a430f587047f26e6a512b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '44e088f5b8f343d6d9ecfe66c1bf5edf', 1851624, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_44e088f5b8f343d6d9ecfe66c1bf5edf', NULL, '2011-08-17 17:53:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '44e088f5b8f343d6d9ecfe66c1bf5edf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '7033eb078c55d22606a4562973d0068c', 1936104, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_7033eb078c55d22606a4562973d0068c', NULL, '2011-11-14 16:16:02'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '7033eb078c55d22606a4562973d0068c');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2011-08-27 11:31:12'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (1844, 'test', 0, '2011-08-17 17:53:00');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 1844 AND stream = 'test' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3b9a06bbf45a430f587047f26e6a512b', '44e088f5b8f343d6d9ecfe66c1bf5edf', '7033eb078c55d22606a4562973d0068c', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20120522
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 10:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '9c468e4f3d4ff9e188f07b88e3f74350', 122149608, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_9c468e4f3d4ff9e188f07b88e3f74350', NULL, '2012-03-13 11:32:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9c468e4f3d4ff9e188f07b88e3f74350');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '170cb44805adc55b65160908e28642fc', 1873640, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_170cb44805adc55b65160908e28642fc', NULL, '2012-05-21 17:04:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '170cb44805adc55b65160908e28642fc');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 17:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2012-02-22 15:16:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20120522, 'stable', 0, '2012-05-21 17:04:18');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20120522 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '9c468e4f3d4ff9e188f07b88e3f74350', '170cb44805adc55b65160908e28642fc', '2493953a68a470cff9385597df3bb3a2', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20120529
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 07:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 07:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 07:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 07:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '9c468e4f3d4ff9e188f07b88e3f74350', 122149608, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_9c468e4f3d4ff9e188f07b88e3f74350', NULL, '2012-03-13 14:32:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9c468e4f3d4ff9e188f07b88e3f74350');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '7110ea84237601dab2ba7c73a6fd89a3', 1874664, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_7110ea84237601dab2ba7c73a6fd89a3', NULL, '2012-06-03 15:12:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '7110ea84237601dab2ba7c73a6fd89a3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 18:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osumeold.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osumeold.exe/b4c3f19676e345fae5d89c5332741725', NULL, '2012-02-22 18:16:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 07:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20120529, 'stable', 0, '2012-06-03 15:12:34');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20120529 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '9c468e4f3d4ff9e188f07b88e3f74350', '7110ea84237601dab2ba7c73a6fd89a3', '2493953a68a470cff9385597df3bb3a2', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20120704
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'f4d5b7cbe979a1f85e8ebcbb0f77c773', 122164968, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_f4d5b7cbe979a1f85e8ebcbb0f77c773', NULL, '2012-07-03 23:45:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f4d5b7cbe979a1f85e8ebcbb0f77c773');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '30e204adc16edb796af2a178b0bfed43', 1902312, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_30e204adc16edb796af2a178b0bfed43', NULL, '2012-07-03 23:45:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '30e204adc16edb796af2a178b0bfed43');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osuframework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osuframework.dll/2493953a68a470cff9385597df3bb3a2', NULL, '2012-05-04 18:36:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2011-08-20 00:30:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20120704, 'stable', 0, '2012-07-03 23:45:56');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20120704 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'f4d5b7cbe979a1f85e8ebcbb0f77c773', '30e204adc16edb796af2a178b0bfed43', '2493953a68a470cff9385597df3bb3a2', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20120725
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 03:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 03:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 03:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-12 22:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-12 22:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-12 22:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 03:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-12 22:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'f4d5b7cbe979a1f85e8ebcbb0f77c773', 122164968, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_f4d5b7cbe979a1f85e8ebcbb0f77c773', NULL, '2012-07-03 04:31:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f4d5b7cbe979a1f85e8ebcbb0f77c773');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '607179cbd461ce0d3e79adcc1263f2b9', 1915112, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_607179cbd461ce0d3e79adcc1263f2b9', NULL, '2012-07-29 12:46:10'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '607179cbd461ce0d3e79adcc1263f2b9');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-07-21 03:46:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2012-02-22 07:16:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 03:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-12 22:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20120725, 'stable', 0, '2012-07-29 12:46:10');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20120725 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'f4d5b7cbe979a1f85e8ebcbb0f77c773', '607179cbd461ce0d3e79adcc1263f2b9', '2493953a68a470cff9385597df3bb3a2', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20120812
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'd783e8ff5654ce363c38cd02bd298562', 129033960, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_d783e8ff5654ce363c38cd02bd298562', NULL, '2012-08-13 14:28:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd783e8ff5654ce363c38cd02bd298562');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'bbc109ccbcf0fbfb916858a1e29c0e02', 1933032, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_bbc109ccbcf0fbfb916858a1e29c0e02', NULL, '2012-08-13 14:28:40'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bbc109ccbcf0fbfb916858a1e29c0e02');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osuframework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osuframework.dll/2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 18:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'b4c3f19676e345fae5d89c5332741725', 294632, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_b4c3f19676e345fae5d89c5332741725', NULL, '2012-02-22 15:16:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b4c3f19676e345fae5d89c5332741725');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20120812, 'stable', 0, '2012-08-13 14:28:40');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20120812 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'd783e8ff5654ce363c38cd02bd298562', 'bbc109ccbcf0fbfb916858a1e29c0e02', '2493953a68a470cff9385597df3bb3a2', 'b4c3f19676e345fae5d89c5332741725', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20120916
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 10:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, 'e278e4b07c7263d76af7149cc91d6808', 130490600, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_e278e4b07c7263d76af7149cc91d6808', NULL, '2012-09-01 14:35:08'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'e278e4b07c7263d76af7149cc91d6808');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '563763d859e0f227affacecd043f77a1', 1942728, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_563763d859e0f227affacecd043f77a1', NULL, '2012-09-17 04:27:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '563763d859e0f227affacecd043f77a1');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 17:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20120916, 'stable', 0, '2012-09-17 04:27:48');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20120916 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', 'e278e4b07c7263d76af7149cc91d6808', '563763d859e0f227affacecd043f77a1', '2493953a68a470cff9385597df3bb3a2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121003
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 20:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '109b680eb982c19cdd36ca4436e8725b', 128621128, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_109b680eb982c19cdd36ca4436e8725b', NULL, '2012-10-28 13:20:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '109b680eb982c19cdd36ca4436e8725b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-10 03:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '8c015e5fed1361b1ac8bd4b06d08f140', 2233928, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_8c015e5fed1361b1ac8bd4b06d08f140', NULL, '2012-10-03 15:06:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '8c015e5fed1361b1ac8bd4b06d08f140');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '52a4149e399cd17190cc0266c5595af2', 297544, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_52a4149e399cd17190cc0266c5595af2', NULL, '2012-10-11 14:13:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '52a4149e399cd17190cc0266c5595af2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121003, 'test', 0, '2012-10-03 15:06:30');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121003 AND stream = 'test' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '109b680eb982c19cdd36ca4436e8725b', '2493953a68a470cff9385597df3bb3a2', '8c015e5fed1361b1ac8bd4b06d08f140', '52a4149e399cd17190cc0266c5595af2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121023
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-09 21:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-09 21:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-09 21:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-12 16:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-12 16:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-12 16:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-29 21:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-12 16:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '805256d34bd0fbd935dd7f696565c1a9', 128571464, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_805256d34bd0fbd935dd7f696565c1a9', NULL, '2012-10-12 12:25:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '805256d34bd0fbd935dd7f696565c1a9');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'c00d001ae6af55d031d6d490cbd750de', 2252872, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_c00d001ae6af55d031d6d490cbd750de', NULL, '2012-10-25 05:18:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c00d001ae6af55d031d6d490cbd750de');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-05-19 01:35:02'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '52a4149e399cd17190cc0266c5595af2', 297544, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_52a4149e399cd17190cc0266c5595af2', NULL, '2012-10-12 12:25:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '52a4149e399cd17190cc0266c5595af2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-09 21:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-12 16:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121023, 'stable', 0, '2012-10-25 05:18:58');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121023 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '805256d34bd0fbd935dd7f696565c1a9', 'c00d001ae6af55d031d6d490cbd750de', '2493953a68a470cff9385597df3bb3a2', '52a4149e399cd17190cc0266c5595af2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121030
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 20:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '109b680eb982c19cdd36ca4436e8725b', 128621128, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_109b680eb982c19cdd36ca4436e8725b', NULL, '2012-10-28 12:20:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '109b680eb982c19cdd36ca4436e8725b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '794195cfc03fb8282f8d5c652ef51d05', 2193480, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_794195cfc03fb8282f8d5c652ef51d05', NULL, '2012-11-03 08:03:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '794195cfc03fb8282f8d5c652ef51d05');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-10 03:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '52a4149e399cd17190cc0266c5595af2', 297544, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_52a4149e399cd17190cc0266c5595af2', NULL, '2012-10-11 14:13:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '52a4149e399cd17190cc0266c5595af2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121030, 'stable', 0, '2012-11-03 08:03:22');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121030 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '109b680eb982c19cdd36ca4436e8725b', '794195cfc03fb8282f8d5c652ef51d05', '2493953a68a470cff9385597df3bb3a2', '52a4149e399cd17190cc0266c5595af2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121115
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 02:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 02:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 02:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-12 21:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-12 21:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-12 21:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 02:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-12 21:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '109b680eb982c19cdd36ca4436e8725b', 128621128, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_109b680eb982c19cdd36ca4436e8725b', NULL, '2012-10-27 19:20:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '109b680eb982c19cdd36ca4436e8725b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'a344a34f30daf1191a8abd46cea9ed92', 2197064, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_a344a34f30daf1191a8abd46cea9ed92', NULL, '2012-11-15 16:07:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a344a34f30daf1191a8abd46cea9ed92');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 09:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '52a4149e399cd17190cc0266c5595af2', 297544, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_52a4149e399cd17190cc0266c5595af2', NULL, '2012-10-10 20:13:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '52a4149e399cd17190cc0266c5595af2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 02:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-12 21:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121115, 'stable', 0, '2003-10-06 16:07:34');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121115 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '109b680eb982c19cdd36ca4436e8725b', 'a344a34f30daf1191a8abd46cea9ed92', '2493953a68a470cff9385597df3bb3a2', '52a4149e399cd17190cc0266c5595af2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121119
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '109b680eb982c19cdd36ca4436e8725b', 128621128, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_109b680eb982c19cdd36ca4436e8725b', NULL, '2012-10-28 03:20:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '109b680eb982c19cdd36ca4436e8725b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'bcae2db2df8be47f94d501d5d254938f', 2207304, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_bcae2db2df8be47f94d501d5d254938f', NULL, '2012-11-19 11:00:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bcae2db2df8be47f94d501d5d254938f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 18:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '52a4149e399cd17190cc0266c5595af2', 297544, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_52a4149e399cd17190cc0266c5595af2', NULL, '2012-10-11 05:13:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '52a4149e399cd17190cc0266c5595af2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121119, 'stable', 0, '2012-11-19 11:00:54');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121119 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '109b680eb982c19cdd36ca4436e8725b', 'bcae2db2df8be47f94d501d5d254938f', '2493953a68a470cff9385597df3bb3a2', '52a4149e399cd17190cc0266c5595af2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121203
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 10:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '109b680eb982c19cdd36ca4436e8725b', 128621128, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_109b680eb982c19cdd36ca4436e8725b', NULL, '2012-11-19 01:39:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '109b680eb982c19cdd36ca4436e8725b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '193dc771794e8099ac86404b43005505', 2099784, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_193dc771794e8099ac86404b43005505', NULL, '2012-12-03 18:03:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '193dc771794e8099ac86404b43005505');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 17:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '52a4149e399cd17190cc0266c5595af2', 297544, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_52a4149e399cd17190cc0266c5595af2', NULL, '2012-10-11 04:13:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '52a4149e399cd17190cc0266c5595af2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121203, 'stable', 0, '2012-12-03 18:03:56');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121203 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '109b680eb982c19cdd36ca4436e8725b', '193dc771794e8099ac86404b43005505', '2493953a68a470cff9385597df3bb3a2', '52a4149e399cd17190cc0266c5595af2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20121223
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-09 19:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-09 19:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-09 19:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-12 14:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-12 14:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-12 14:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-29 19:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-12 14:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '97939c3380580bcecad7000a33f56831', 122284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_97939c3380580bcecad7000a33f56831', NULL, '2012-12-22 05:05:40'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '97939c3380580bcecad7000a33f56831');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'b2696bc31ab065bd6385097c7af94fd5', 2159688, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_b2696bc31ab065bd6385097c7af94fd5', NULL, '2012-12-23 19:39:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b2696bc31ab065bd6385097c7af94fd5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 02:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-09 19:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-12 14:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20121223, 'stable', 0, '2012-12-23 19:39:24');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20121223 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '97939c3380580bcecad7000a33f56831', 'b2696bc31ab065bd6385097c7af94fd5', '2493953a68a470cff9385597df3bb3a2', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130319
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '4ccd18c0170017e99da4582d1a18d237', 2263624, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_4ccd18c0170017e99da4582d1a18d237', NULL, '2013-03-19 15:11:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '4ccd18c0170017e99da4582d1a18d237');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '2493953a68a470cff9385597df3bb3a2', 12520, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_2493953a68a470cff9385597df3bb3a2', NULL, '2012-04-09 18:42:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2493953a68a470cff9385597df3bb3a2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '7c378c01712baeec8855eb4319d539bd', 8110152, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_7c378c01712baeec8855eb4319d539bd', NULL, '2013-03-19 15:10:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '7c378c01712baeec8855eb4319d539bd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'ff57703002fd479609d3e0d3ccd01808', 8212552, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_ff57703002fd479609d3e0d3ccd01808', NULL, '2013-03-19 15:11:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ff57703002fd479609d3e0d3ccd01808');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '560f66a9b28da605f100b41549670b29', 298568, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_560f66a9b28da605f100b41549670b29', NULL, '2013-03-19 15:10:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '560f66a9b28da605f100b41549670b29');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130319, 'stable', 0, '2013-03-19 15:11:00');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130319 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '4ccd18c0170017e99da4582d1a18d237', '2493953a68a470cff9385597df3bb3a2', '7c378c01712baeec8855eb4319d539bd', 'ff57703002fd479609d3e0d3ccd01808', '560f66a9b28da605f100b41549670b29', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130329
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '9be427a7dd68e7357b4b8264b8ea15c1', 2263112, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_9be427a7dd68e7357b4b8264b8ea15c1', NULL, '2013-03-29 05:40:38'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9be427a7dd68e7357b4b8264b8ea15c1');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '82d0d99818b6a1d05eea47fcc1c13acd', 14920, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_82d0d99818b6a1d05eea47fcc1c13acd', NULL, '2013-03-19 15:58:38'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '82d0d99818b6a1d05eea47fcc1c13acd');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '83387a556c4ab1516c5d657aeedcc2ee', 8111688, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_83387a556c4ab1516c5d657aeedcc2ee', NULL, '2013-03-29 05:40:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '83387a556c4ab1516c5d657aeedcc2ee');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'a6e74f0fcd9b2bfb7181c825fcf62f22', 8223304, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_a6e74f0fcd9b2bfb7181c825fcf62f22', NULL, '2013-03-29 05:40:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a6e74f0fcd9b2bfb7181c825fcf62f22');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '3c7451e7b93af27b15e382ce64768ffb', 298568, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_3c7451e7b93af27b15e382ce64768ffb', NULL, '2013-03-29 05:40:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3c7451e7b93af27b15e382ce64768ffb');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130329, 'stable', 0, '2013-03-29 05:40:38');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130329 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '9be427a7dd68e7357b4b8264b8ea15c1', '82d0d99818b6a1d05eea47fcc1c13acd', '83387a556c4ab1516c5d657aeedcc2ee', 'a6e74f0fcd9b2bfb7181c825fcf62f22', '3c7451e7b93af27b15e382ce64768ffb', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130513
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '693f708adcc29a3778917fe496ae5cb2', 2297928, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_693f708adcc29a3778917fe496ae5cb2', NULL, '2013-05-15 18:18:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '693f708adcc29a3778917fe496ae5cb2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '4c6428d34590e1fd05d8ed28e5efa263', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_4c6428d34590e1fd05d8ed28e5efa263', NULL, '2013-05-15 18:18:28'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '4c6428d34590e1fd05d8ed28e5efa263');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '83387a556c4ab1516c5d657aeedcc2ee', 8111688, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_83387a556c4ab1516c5d657aeedcc2ee', NULL, '2013-03-29 05:40:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '83387a556c4ab1516c5d657aeedcc2ee');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'fa58e536b2b0155554d333ca83d5c54f', 8280136, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_fa58e536b2b0155554d333ca83d5c54f', NULL, '2013-05-15 18:18:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'fa58e536b2b0155554d333ca83d5c54f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130513, 'stable', 0, '2013-05-15 18:18:32');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130513 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '693f708adcc29a3778917fe496ae5cb2', '4c6428d34590e1fd05d8ed28e5efa263', '83387a556c4ab1516c5d657aeedcc2ee', 'fa58e536b2b0155554d333ca83d5c54f', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130603
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '5ba40d7cdcfab111d59b9067f96778c8', 2355784, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_5ba40d7cdcfab111d59b9067f96778c8', NULL, '2013-06-03 17:10:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5ba40d7cdcfab111d59b9067f96778c8');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '4d45b2c0ef63cd844bf1c0fde3c6b202', 10228296, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_4d45b2c0ef63cd844bf1c0fde3c6b202', NULL, '2013-06-03 17:10:02'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '4d45b2c0ef63cd844bf1c0fde3c6b202');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130603, 'stable', 0, '2013-06-03 17:10:00');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130603 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '5ba40d7cdcfab111d59b9067f96778c8', '4c6428d34590e1fd05d8ed28e5efa263', '4d45b2c0ef63cd844bf1c0fde3c6b202', 'fa58e536b2b0155554d333ca83d5c54f', '0091e03a4279d8cf5bd3015771a28eb1', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130626
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '857cfa177d75b35ab171e3633cb22e0c', 2363976, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_857cfa177d75b35ab171e3633cb22e0c', NULL, '2013-06-26 10:38:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '857cfa177d75b35ab171e3633cb22e0c');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-06-26 10:38:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130626, 'stable', 0, '2013-06-26 10:38:32');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130626 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '857cfa177d75b35ab171e3633cb22e0c', '63cf5329e812b270071ddd4b1bc99423', '4d45b2c0ef63cd844bf1c0fde3c6b202', 'fa58e536b2b0155554d333ca83d5c54f', '89222748245550e3cb5b80a62c0e3ba1', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130716
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 20:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-20 00:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '8200c36a46359b1abe8db68a0a8d40ea', 2358272, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_8200c36a46359b1abe8db68a0a8d40ea', NULL, '2013-07-17 23:29:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '8200c36a46359b1abe8db68a0a8d40ea');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-07-13 02:17:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '4d45b2c0ef63cd844bf1c0fde3c6b202', 10228296, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_4d45b2c0ef63cd844bf1c0fde3c6b202', NULL, '2013-06-02 07:24:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '4d45b2c0ef63cd844bf1c0fde3c6b202');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'fa58e536b2b0155554d333ca83d5c54f', 8280136, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_fa58e536b2b0155554d333ca83d5c54f', NULL, '2013-05-16 03:18:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'fa58e536b2b0155554d333ca83d5c54f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '89222748245550e3cb5b80a62c0e3ba1', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_89222748245550e3cb5b80a62c0e3ba1', NULL, '2013-07-13 02:17:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '89222748245550e3cb5b80a62c0e3ba1');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130716, 'stable', 0, '2013-07-17 23:29:36');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130716 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '8200c36a46359b1abe8db68a0a8d40ea', '63cf5329e812b270071ddd4b1bc99423', '4d45b2c0ef63cd844bf1c0fde3c6b202', 'fa58e536b2b0155554d333ca83d5c54f', '89222748245550e3cb5b80a62c0e3ba1', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130802
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '031713d5025403c1f7675a2603ce3149', 2387016, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_031713d5025403c1f7675a2603ce3149', NULL, '2013-08-02 02:04:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '031713d5025403c1f7675a2603ce3149');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-08-01 12:18:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '4d45b2c0ef63cd844bf1c0fde3c6b202', 10228296, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_4d45b2c0ef63cd844bf1c0fde3c6b202', NULL, '2013-06-01 20:24:34'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '4d45b2c0ef63cd844bf1c0fde3c6b202');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '264593ce2e93c5db29466510e75334da', 8346184, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_264593ce2e93c5db29466510e75334da', NULL, '2013-08-01 11:28:08'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '264593ce2e93c5db29466510e75334da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb287fbe2c003beaf5bcd7c6359386da', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb287fbe2c003beaf5bcd7c6359386da', NULL, '2013-08-02 00:53:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb287fbe2c003beaf5bcd7c6359386da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130802, 'stable', 1, '2013-08-02 02:04:56');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130802 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '031713d5025403c1f7675a2603ce3149', '63cf5329e812b270071ddd4b1bc99423', '4d45b2c0ef63cd844bf1c0fde3c6b202', '264593ce2e93c5db29466510e75334da', 'bb287fbe2c003beaf5bcd7c6359386da', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130823
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 13:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '8a0d528ecc6c317f890924e413d00957', 2390600, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_8a0d528ecc6c317f890924e413d00957', NULL, '2013-08-23 06:50:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '8a0d528ecc6c317f890924e413d00957');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-08-01 12:18:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'cabc388e515d7fe32b03f8c6a4ccd51f', 10260040, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_cabc388e515d7fe32b03f8c6a4ccd51f', NULL, '2013-08-13 08:40:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'cabc388e515d7fe32b03f8c6a4ccd51f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '1371a4bb7858672edd7cd5a522630898', 8353352, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_1371a4bb7858672edd7cd5a522630898', NULL, '2013-08-07 21:30:20'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '1371a4bb7858672edd7cd5a522630898');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb287fbe2c003beaf5bcd7c6359386da', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb287fbe2c003beaf5bcd7c6359386da', NULL, '2013-08-02 00:53:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb287fbe2c003beaf5bcd7c6359386da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130823, 'stable', 3, '2013-08-23 06:50:06');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130823 AND stream = 'stable' AND subversion = 3 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '8a0d528ecc6c317f890924e413d00957', '63cf5329e812b270071ddd4b1bc99423', 'cabc388e515d7fe32b03f8c6a4ccd51f', '1371a4bb7858672edd7cd5a522630898', 'bb287fbe2c003beaf5bcd7c6359386da', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20130915
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'e7e63c58f15671e4cf947c2f9a3ae0b2', 2412104, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_e7e63c58f15671e4cf947c2f9a3ae0b2', NULL, '2013-09-15 21:09:14'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'e7e63c58f15671e4cf947c2f9a3ae0b2');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-08-01 12:18:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3e88234c19e4907f1d99597e840987a', 10263624, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3e88234c19e4907f1d99597e840987a', NULL, '2013-09-04 00:54:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3e88234c19e4907f1d99597e840987a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '772e43f7e73ac58379dd2276e8c64ced', 8368712, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_772e43f7e73ac58379dd2276e8c64ced', NULL, '2013-08-29 15:54:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '772e43f7e73ac58379dd2276e8c64ced');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb287fbe2c003beaf5bcd7c6359386da', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb287fbe2c003beaf5bcd7c6359386da', NULL, '2013-08-02 00:53:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb287fbe2c003beaf5bcd7c6359386da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20130915, 'stable', 4, '2013-09-15 21:09:14');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20130915 AND stream = 'stable' AND subversion = 4 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', 'e7e63c58f15671e4cf947c2f9a3ae0b2', '63cf5329e812b270071ddd4b1bc99423', 'a3e88234c19e4907f1d99597e840987a', '772e43f7e73ac58379dd2276e8c64ced', 'bb287fbe2c003beaf5bcd7c6359386da', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20131009
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '7cae49efab9618657096ca5ef318fd9c', 2414664, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_7cae49efab9618657096ca5ef318fd9c', NULL, '2013-10-09 09:15:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '7cae49efab9618657096ca5ef318fd9c');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-08-01 14:18:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3e88234c19e4907f1d99597e840987a', 10263624, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3e88234c19e4907f1d99597e840987a', NULL, '2013-09-04 02:54:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3e88234c19e4907f1d99597e840987a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '772e43f7e73ac58379dd2276e8c64ced', 8368712, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_772e43f7e73ac58379dd2276e8c64ced', NULL, '2013-08-29 17:54:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '772e43f7e73ac58379dd2276e8c64ced');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20131009, 'stable', 0, '2013-10-09 09:15:04');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20131009 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '7cae49efab9618657096ca5ef318fd9c', '63cf5329e812b270071ddd4b1bc99423', 'a3e88234c19e4907f1d99597e840987a', '772e43f7e73ac58379dd2276e8c64ced', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20131024
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'c5a0237e0ced21fd7aa7e9f5a02ebba3', 2417736, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_c5a0237e0ced21fd7aa7e9f5a02ebba3', NULL, '2013-10-24 12:36:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c5a0237e0ced21fd7aa7e9f5a02ebba3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, '63cf5329e812b270071ddd4b1bc99423', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_63cf5329e812b270071ddd4b1bc99423', NULL, '2013-08-01 12:18:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '63cf5329e812b270071ddd4b1bc99423');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3e88234c19e4907f1d99597e840987a', 10263624, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3e88234c19e4907f1d99597e840987a', NULL, '2013-09-04 00:54:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3e88234c19e4907f1d99597e840987a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '772e43f7e73ac58379dd2276e8c64ced', 8368712, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_772e43f7e73ac58379dd2276e8c64ced', NULL, '2013-08-29 15:54:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '772e43f7e73ac58379dd2276e8c64ced');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb287fbe2c003beaf5bcd7c6359386da', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb287fbe2c003beaf5bcd7c6359386da', NULL, '2013-08-02 00:53:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb287fbe2c003beaf5bcd7c6359386da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20131024, 'stable', 1, '2013-10-24 12:36:04');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20131024 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', 'c5a0237e0ced21fd7aa7e9f5a02ebba3', '63cf5329e812b270071ddd4b1bc99423', 'a3e88234c19e4907f1d99597e840987a', '772e43f7e73ac58379dd2276e8c64ced', 'bb287fbe2c003beaf5bcd7c6359386da', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20131113
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '1521c455a4e6d5cf9525c84743fbabd0', 2418760, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_1521c455a4e6d5cf9525c84743fbabd0', NULL, '2013-11-13 03:39:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '1521c455a4e6d5cf9525c84743fbabd0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 02:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3e88234c19e4907f1d99597e840987a', 10263624, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3e88234c19e4907f1d99597e840987a', NULL, '2013-09-04 00:54:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3e88234c19e4907f1d99597e840987a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '122eeec4be4268db6a57a2d2fd7fecd5', 8291400, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_122eeec4be4268db6a57a2d2fd7fecd5', NULL, '2013-11-05 23:18:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '122eeec4be4268db6a57a2d2fd7fecd5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb287fbe2c003beaf5bcd7c6359386da', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb287fbe2c003beaf5bcd7c6359386da', NULL, '2013-08-02 00:53:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb287fbe2c003beaf5bcd7c6359386da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20131113, 'stable', 2, '2013-11-13 03:39:50');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20131113 AND stream = 'stable' AND subversion = 2 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '1521c455a4e6d5cf9525c84743fbabd0', 'f757b90ecd5f55f7b1ade82d878e16b6', 'a3e88234c19e4907f1d99597e840987a', '122eeec4be4268db6a57a2d2fd7fecd5', 'bb287fbe2c003beaf5bcd7c6359386da', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20131129
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 20:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-20 00:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'd47ef68157ed98a4973a7b0d53eddc78', 2419784, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_d47ef68157ed98a4973a7b0d53eddc78', NULL, '2013-11-29 16:58:42'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd47ef68157ed98a4973a7b0d53eddc78');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 11:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3e88234c19e4907f1d99597e840987a', 10263624, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3e88234c19e4907f1d99597e840987a', NULL, '2013-09-04 11:54:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3e88234c19e4907f1d99597e840987a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '122eeec4be4268db6a57a2d2fd7fecd5', 8291400, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_122eeec4be4268db6a57a2d2fd7fecd5', NULL, '2013-11-06 08:18:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '122eeec4be4268db6a57a2d2fd7fecd5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'bb287fbe2c003beaf5bcd7c6359386da', 300104, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_bb287fbe2c003beaf5bcd7c6359386da', NULL, '2013-08-02 11:53:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'bb287fbe2c003beaf5bcd7c6359386da');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 20:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 15:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20131129, 'stable', 1, '2013-11-29 16:58:42');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20131129 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', 'd47ef68157ed98a4973a7b0d53eddc78', 'f757b90ecd5f55f7b1ade82d878e16b6', 'a3e88234c19e4907f1d99597e840987a', '122eeec4be4268db6a57a2d2fd7fecd5', 'bb287fbe2c003beaf5bcd7c6359386da', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20131216
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2014-08-18 08:16:59'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2014-08-18 08:16:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2014-08-18 08:16:59'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2014-08-18 08:17:03'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2014-08-18 08:17:13'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2014-08-18 08:17:15'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2014-08-18 08:17:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '51a225508fc0cfcad5aeaf4543b02d48', 2987592, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_51a225508fc0cfcad5aeaf4543b02d48', NULL, '2013-12-16 13:16:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '51a225508fc0cfcad5aeaf4543b02d48');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 02:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3e88234c19e4907f1d99597e840987a', 10263624, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3e88234c19e4907f1d99597e840987a', NULL, '2013-09-04 00:54:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3e88234c19e4907f1d99597e840987a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '6e811b9a410393e8710c7785882413be', 8311368, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_6e811b9a410393e8710c7785882413be', NULL, '2013-12-16 13:16:32'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '6e811b9a410393e8710c7785882413be');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '55f4feae3a9a14c804b70fa266c44dd3', 301128, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_55f4feae3a9a14c804b70fa266c44dd3', NULL, '2013-03-19 15:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '55f4feae3a9a14c804b70fa266c44dd3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2014-08-18 08:18:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2014-08-18 08:18:03'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20131216, 'stable', 26, '2013-12-16 13:16:32');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20131216 AND stream = 'stable' AND subversion = 26 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '51a225508fc0cfcad5aeaf4543b02d48', 'f757b90ecd5f55f7b1ade82d878e16b6', 'a3e88234c19e4907f1d99597e840987a', '6e811b9a410393e8710c7785882413be', '55f4feae3a9a14c804b70fa266c44dd3', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140114
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 11:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 16:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '432a7697725473690496974f637124ad', 2423368, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_432a7697725473690496974f637124ad', NULL, '2014-01-16 18:06:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '432a7697725473690496974f637124ad');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 03:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a03b036f333f3cd2bcd7760c48a45aab', 10263112, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a03b036f333f3cd2bcd7760c48a45aab', NULL, '2014-01-16 18:06:02'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a03b036f333f3cd2bcd7760c48a45aab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'cec7342df43daa13af316f4d933c60b8', 8284744, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_cec7342df43daa13af316f4d933c60b8', NULL, '2014-01-16 18:06:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'cec7342df43daa13af316f4d933c60b8');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '55f4feae3a9a14c804b70fa266c44dd3', 301128, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_55f4feae3a9a14c804b70fa266c44dd3', NULL, '2013-12-31 15:12:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '55f4feae3a9a14c804b70fa266c44dd3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 11:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 06:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140114, 'stable', 2, '2014-01-16 18:06:30');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140114 AND stream = 'stable' AND subversion = 2 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '432a7697725473690496974f637124ad', 'f757b90ecd5f55f7b1ade82d878e16b6', 'a03b036f333f3cd2bcd7760c48a45aab', 'cec7342df43daa13af316f4d933c60b8', '55f4feae3a9a14c804b70fa266c44dd3', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140119
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 09:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 14:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '68c4b3cef5f01df4a236daee959c3966', 2424392, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_68c4b3cef5f01df4a236daee959c3966', NULL, '2014-01-19 20:06:42'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '68c4b3cef5f01df4a236daee959c3966');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 01:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'a3b05a8471168b574b86ddcf403b3f95', 12227656, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_a3b05a8471168b574b86ddcf403b3f95', NULL, '2014-01-19 20:06:38'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a3b05a8471168b574b86ddcf403b3f95');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'cec7342df43daa13af316f4d933c60b8', 8284744, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_cec7342df43daa13af316f4d933c60b8', NULL, '2014-01-14 16:25:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'cec7342df43daa13af316f4d933c60b8');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '55f4feae3a9a14c804b70fa266c44dd3', 301128, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_55f4feae3a9a14c804b70fa266c44dd3', NULL, '2013-12-16 15:30:54'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '55f4feae3a9a14c804b70fa266c44dd3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 09:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 04:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140119, 'stable', 1, '2014-01-19 20:06:42');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140119 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '68c4b3cef5f01df4a236daee959c3966', 'f757b90ecd5f55f7b1ade82d878e16b6', 'a3b05a8471168b574b86ddcf403b3f95', 'cec7342df43daa13af316f4d933c60b8', '55f4feae3a9a14c804b70fa266c44dd3', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140127
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 10:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2014-01-04 18:56:42'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '5d6f3d09470cd54cc47e0bcf5e8e646a', 2425416, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_5d6f3d09470cd54cc47e0bcf5e8e646a', NULL, '2014-01-26 22:29:40'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5d6f3d09470cd54cc47e0bcf5e8e646a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2014-01-04 18:56:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '9c319b2c3c5b6f0038d7f62fa31b04e0', 12446792, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_9c319b2c3c5b6f0038d7f62fa31b04e0', NULL, '2014-01-23 13:57:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9c319b2c3c5b6f0038d7f62fa31b04e0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '5a0d1105d83bbd85232ef5ea3e89d86a', 8291912, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_5a0d1105d83bbd85232ef5ea3e89d86a', NULL, '2014-01-23 13:58:00'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5a0d1105d83bbd85232ef5ea3e89d86a');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, '55f4feae3a9a14c804b70fa266c44dd3', 301128, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_55f4feae3a9a14c804b70fa266c44dd3', NULL, '2014-01-04 18:56:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '55f4feae3a9a14c804b70fa266c44dd3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140127, 'stable', 0, '2014-01-26 22:29:40');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140127 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', '5d6f3d09470cd54cc47e0bcf5e8e646a', 'f757b90ecd5f55f7b1ade82d878e16b6', '9c319b2c3c5b6f0038d7f62fa31b04e0', '5a0d1105d83bbd85232ef5ea3e89d86a', '55f4feae3a9a14c804b70fa266c44dd3', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140323
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 10:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 14:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'd320f8e871ee5046688c9876996aa20d', 2478152, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_d320f8e871ee5046688c9876996aa20d', NULL, '2014-03-23 20:30:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd320f8e871ee5046688c9876996aa20d');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 01:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'dfc0dd0200502eb939ab4c35adde2c41', 12618312, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_dfc0dd0200502eb939ab4c35adde2c41', NULL, '2014-03-09 23:07:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'dfc0dd0200502eb939ab4c35adde2c41');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'ceace8e1666f7c09f4ac7b82106bab84', 8839240, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_ceace8e1666f7c09f4ac7b82106bab84', NULL, '2014-03-15 12:49:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ceace8e1666f7c09f4ac7b82106bab84');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 13:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 10:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 05:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140323, 'stable', 3, '2014-03-23 20:30:48');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140323 AND stream = 'stable' AND subversion = 3 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', 'd320f8e871ee5046688c9876996aa20d', 'f757b90ecd5f55f7b1ade82d878e16b6', 'dfc0dd0200502eb939ab4c35adde2c41', 'ceace8e1666f7c09f4ac7b82106bab84', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140410
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 08:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '3e1455084499f3e86307a6b646ddb576', 122274888, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_3e1455084499f3e86307a6b646ddb576', NULL, '2013-03-19 13:10:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3e1455084499f3e86307a6b646ddb576');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'd0673235f3aa1d8e4240712333d5b995', 2435656, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_d0673235f3aa1d8e4240712333d5b995', NULL, '2014-04-10 01:34:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd0673235f3aa1d8e4240712333d5b995');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 00:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '20ad75d822caf0a4b8b8a872350e0d42', 12623944, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_20ad75d822caf0a4b8b8a872350e0d42', NULL, '2014-04-08 15:26:56'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '20ad75d822caf0a4b8b8a872350e0d42');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'ceace8e1666f7c09f4ac7b82106bab84', 8839240, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_ceace8e1666f7c09f4ac7b82106bab84', NULL, '2014-03-15 11:49:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ceace8e1666f7c09f4ac7b82106bab84');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 13:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140410, 'stable', 1, '2014-04-10 01:34:30');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140410 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '3e1455084499f3e86307a6b646ddb576', 'd0673235f3aa1d8e4240712333d5b995', 'f757b90ecd5f55f7b1ade82d878e16b6', '20ad75d822caf0a4b8b8a872350e0d42', 'ceace8e1666f7c09f4ac7b82106bab84', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140613
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 07:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '5e9af52e463636f704e64c2d28f9cf40', 50284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_5e9af52e463636f704e64c2d28f9cf40', NULL, '2014-05-07 17:47:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5e9af52e463636f704e64c2d28f9cf40');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '8c9f9031fdd764d3670371c932ff5aae', 2536008, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_8c9f9031fdd764d3670371c932ff5aae', NULL, '2014-06-13 18:49:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '8c9f9031fdd764d3670371c932ff5aae');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 00:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'c55e59bcb6990e94427da02ee7724365', 14651464, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_c55e59bcb6990e94427da02ee7724365', NULL, '2014-06-04 14:49:46'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c55e59bcb6990e94427da02ee7724365');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'f31b02f915b2b186426723f6ee4481c4', 8908360, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_f31b02f915b2b186426723f6ee4481c4', NULL, '2014-05-13 14:22:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f31b02f915b2b186426723f6ee4481c4');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 13:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140613, 'stable', 18, '2014-06-13 18:49:30');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140613 AND stream = 'stable' AND subversion = 18 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '5e9af52e463636f704e64c2d28f9cf40', '8c9f9031fdd764d3670371c932ff5aae', 'f757b90ecd5f55f7b1ade82d878e16b6', 'c55e59bcb6990e94427da02ee7724365', 'f31b02f915b2b186426723f6ee4481c4', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140616
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 00:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 00:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 00:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-12 19:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-12 19:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-12 19:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 00:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-12 19:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '5e9af52e463636f704e64c2d28f9cf40', 50284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_5e9af52e463636f704e64c2d28f9cf40', NULL, '2014-05-07 14:22:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5e9af52e463636f704e64c2d28f9cf40');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-12 15:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, 'c55e59bcb6990e94427da02ee7724365', 14651464, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_c55e59bcb6990e94427da02ee7724365', NULL, '2014-06-04 08:47:48'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c55e59bcb6990e94427da02ee7724365');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '40cf5157c50b14630cc935168b853cfc', 2540104, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_40cf5157c50b14630cc935168b853cfc', NULL, '2014-06-16 07:32:18'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '40cf5157c50b14630cc935168b853cfc');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'e926f9c3564aabb1a5189b6fbafcf659', 9110088, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_e926f9c3564aabb1a5189b6fbafcf659', NULL, '2014-06-15 05:55:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'e926f9c3564aabb1a5189b6fbafcf659');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 04:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 00:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-12 19:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140616, 'test', 0, '2014-06-16 07:32:18');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140616 AND stream = 'test' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '5e9af52e463636f704e64c2d28f9cf40', 'f757b90ecd5f55f7b1ade82d878e16b6', 'c55e59bcb6990e94427da02ee7724365', '40cf5157c50b14630cc935168b853cfc', 'e926f9c3564aabb1a5189b6fbafcf659', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140624
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 08:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '5e9af52e463636f704e64c2d28f9cf40', 50284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_5e9af52e463636f704e64c2d28f9cf40', NULL, '2014-05-07 18:47:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5e9af52e463636f704e64c2d28f9cf40');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '9f1e287f663e89785886aa52aa2d1335', 2546760, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_9f1e287f663e89785886aa52aa2d1335', NULL, '2014-06-24 18:20:30'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '9f1e287f663e89785886aa52aa2d1335');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 00:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '881ce42c3313519c31976028a4e03e56', 16392264, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_881ce42c3313519c31976028a4e03e56', NULL, '2014-06-22 14:09:58'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '881ce42c3313519c31976028a4e03e56');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '10b6c7bccaf2675e70874158f3baa7af', 10088008, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_10b6c7bccaf2675e70874158f3baa7af', NULL, '2014-06-24 18:20:28'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '10b6c7bccaf2675e70874158f3baa7af');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 13:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 08:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 03:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140624, 'stable', 0, '2014-06-24 18:20:30');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140624 AND stream = 'stable' AND subversion = 0 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('d41d8cd98f00b204e9800998ecf8427e', 'b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '5e9af52e463636f704e64c2d28f9cf40', '9f1e287f663e89785886aa52aa2d1335', 'f757b90ecd5f55f7b1ade82d878e16b6', '881ce42c3313519c31976028a4e03e56', '10b6c7bccaf2675e70874158f3baa7af', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140628
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 18:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 18:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 18:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 13:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 13:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 13:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 18:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 13:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '5e9af52e463636f704e64c2d28f9cf40', 50284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_5e9af52e463636f704e64c2d28f9cf40', NULL, '2014-05-08 04:47:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5e9af52e463636f704e64c2d28f9cf40');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, '09afdb34e7618fd23ff385a244380da3', 2549832, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_09afdb34e7618fd23ff385a244380da3', NULL, '2014-06-28 21:06:14'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '09afdb34e7618fd23ff385a244380da3');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 09:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '881ce42c3313519c31976028a4e03e56', 16392264, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_881ce42c3313519c31976028a4e03e56', NULL, '2014-06-22 02:29:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '881ce42c3313519c31976028a4e03e56');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, 'd491140ad8d8ce6228f1e4c19a146655', 11573832, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_d491140ad8d8ce6228f1e4c19a146655', NULL, '2014-06-27 21:45:20'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd491140ad8d8ce6228f1e4c19a146655');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 22:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 18:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 13:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140628, 'stable', 6, '2014-06-28 21:06:14');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140628 AND stream = 'stable' AND subversion = 6 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('d41d8cd98f00b204e9800998ecf8427e', 'b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '5e9af52e463636f704e64c2d28f9cf40', '09afdb34e7618fd23ff385a244380da3', 'f757b90ecd5f55f7b1ade82d878e16b6', '881ce42c3313519c31976028a4e03e56', 'd491140ad8d8ce6228f1e4c19a146655', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140725
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 07:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '5e9af52e463636f704e64c2d28f9cf40', 50284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_5e9af52e463636f704e64c2d28f9cf40', NULL, '2014-05-07 17:47:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5e9af52e463636f704e64c2d28f9cf40');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'fbd0e00d0a6c28e847e03dc31e5e3077', 2565704, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_fbd0e00d0a6c28e847e03dc31e5e3077', NULL, '2014-07-26 01:11:16'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'fbd0e00d0a6c28e847e03dc31e5e3077');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-12 23:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '881ce42c3313519c31976028a4e03e56', 16392264, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_881ce42c3313519c31976028a4e03e56', NULL, '2014-06-21 15:29:06'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '881ce42c3313519c31976028a4e03e56');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '307ea9788c984e3f6117799c4c20a70b', 11808840, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_307ea9788c984e3f6117799c4c20a70b', NULL, '2014-07-22 08:09:50'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '307ea9788c984e3f6117799c4c20a70b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 12:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140725, 'stable', 10, '2014-07-26 01:11:16');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140725 AND stream = 'stable' AND subversion = 10 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('d41d8cd98f00b204e9800998ecf8427e', 'b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '5e9af52e463636f704e64c2d28f9cf40', 'fbd0e00d0a6c28e847e03dc31e5e3077', 'f757b90ecd5f55f7b1ade82d878e16b6', '881ce42c3313519c31976028a4e03e56', '307ea9788c984e3f6117799c4c20a70b', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

-- Release b20140814
INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avcodec-51.dll', 0, 'b22bf1e4ecd4be3d909dc68ccab74eec', 4409856, NULL, 'http://cdn.titanic.sh/r/avcodec-51.dll/f_b22bf1e4ecd4be3d909dc68ccab74eec', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'b22bf1e4ecd4be3d909dc68ccab74eec');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avformat-52.dll', 0, '2e7a800133625f827cf46aa0bb1af800', 711680, NULL, 'http://cdn.titanic.sh/r/avformat-52.dll/f_2e7a800133625f827cf46aa0bb1af800', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '2e7a800133625f827cf46aa0bb1af800');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'avutil-49.dll', 0, 'c870147dff89c95c81f8fbdfbc6344ac', 62464, NULL, 'http://cdn.titanic.sh/r/avutil-49.dll/f_c870147dff89c95c81f8fbdfbc6344ac', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c870147dff89c95c81f8fbdfbc6344ac');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass.dll', 0, 'd7f05d3fa5e745e02e1de41821ccccaf', 98872, NULL, 'http://cdn.titanic.sh/r/bass.dll/f_d7f05d3fa5e745e02e1de41821ccccaf', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'd7f05d3fa5e745e02e1de41821ccccaf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'bass_fx.dll', 0, '3f6de1b4476dd77b42adab8cd8c1f7b5', 29784, NULL, 'http://cdn.titanic.sh/r/bass_fx.dll/f_3f6de1b4476dd77b42adab8cd8c1f7b5', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '3f6de1b4476dd77b42adab8cd8c1f7b5');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'd3dx9_31.dll', 0, 'eea5e428ce63804f9b12d21c97b5968f', 4379984, NULL, 'http://cdn.titanic.sh/r/d3dx9_31.dll/f_eea5e428ce63804f9b12d21c97b5968f', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'eea5e428ce63804f9b12d21c97b5968f');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Ink.dll', 0, 'a02ee61542caae25f8a44c9428d30247', 516096, NULL, 'http://cdn.titanic.sh/r/Microsoft.Ink.dll/f_a02ee61542caae25f8a44c9428d30247', NULL, '2009-06-30 07:28:24'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'a02ee61542caae25f8a44c9428d30247');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'Microsoft.Xna.Framework.dll', 0, '45a786658d3f69717652fed471d03ee0', 749568, NULL, 'http://cdn.titanic.sh/r/Microsoft.Xna.Framework.dll/f_45a786658d3f69717652fed471d03ee0', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '45a786658d3f69717652fed471d03ee0');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu.dll', 0, '5e9af52e463636f704e64c2d28f9cf40', 50284616, NULL, 'http://cdn.titanic.sh/r/osu.dll/f_5e9af52e463636f704e64c2d28f9cf40', NULL, '2014-05-07 17:47:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '5e9af52e463636f704e64c2d28f9cf40');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!.exe', 0, 'f4dce9e69b168fd397255264e85c70c1', 2572360, NULL, 'http://cdn.titanic.sh/r/osu!.exe/f_f4dce9e69b168fd397255264e85c70c1', NULL, '2014-08-14 15:02:02'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f4dce9e69b168fd397255264e85c70c1');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!framework.dll', 0, 'f757b90ecd5f55f7b1ade82d878e16b6', 15944, NULL, 'http://cdn.titanic.sh/r/osu!framework.dll/f_f757b90ecd5f55f7b1ade82d878e16b6', NULL, '2013-11-13 00:33:04'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'f757b90ecd5f55f7b1ade82d878e16b6');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!gameplay.dll', 0, '525a20c6e9d2c58c085c42198308a3bf', 16394312, NULL, 'http://cdn.titanic.sh/r/osu!gameplay.dll/f_525a20c6e9d2c58c085c42198308a3bf', NULL, '2014-08-12 17:30:12'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '525a20c6e9d2c58c085c42198308a3bf');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osu!ui.dll', 0, '307ea9788c984e3f6117799c4c20a70b', 11808840, NULL, 'http://cdn.titanic.sh/r/osu!ui.dll/f_307ea9788c984e3f6117799c4c20a70b', NULL, '2014-07-21 20:44:36'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '307ea9788c984e3f6117799c4c20a70b');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'osume.exe', 0, 'c99710c67fcf9434edaee80b6e07bfab', 299080, NULL, 'http://cdn.titanic.sh/r/osume.exe/f_c99710c67fcf9434edaee80b6e07bfab', NULL, '2014-03-18 13:01:52'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'c99710c67fcf9434edaee80b6e07bfab');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'pthreadGC2.dll', 0, 'ce931021e18f385f519e945a8a10548e', 60273, NULL, 'http://cdn.titanic.sh/r/pthreadGC2.dll/f_ce931021e18f385f519e945a8a10548e', NULL, '2009-05-10 07:01:22'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = 'ce931021e18f385f519e945a8a10548e');

INSERT INTO releases_official_files (filename, file_version, file_hash, filesize, patch_id, url_full, url_patch, timestamp)
SELECT 'x3daudio1_1.dll', 0, '121b131eaa369d8f58dacc5c39a77d80', 15128, NULL, 'http://cdn.titanic.sh/r/x3daudio1_1.dll/f_121b131eaa369d8f58dacc5c39a77d80', NULL, '2009-10-13 02:53:44'
WHERE NOT EXISTS (SELECT 1 FROM releases_official_files WHERE file_hash = '121b131eaa369d8f58dacc5c39a77d80');

INSERT INTO releases_official (version, stream, subversion, created_at)
VALUES (20140814, 'stable', 1, '2014-08-14 15:02:02');

INSERT INTO releases_official_entries (release_id, file_id)
SELECT
    (SELECT id FROM releases_official WHERE version = 20140814 AND stream = 'stable' AND subversion = 1 ORDER BY id DESC LIMIT 1),
    id
FROM releases_official_files
WHERE file_hash IN ('d41d8cd98f00b204e9800998ecf8427e', 'b22bf1e4ecd4be3d909dc68ccab74eec', '2e7a800133625f827cf46aa0bb1af800', 'c870147dff89c95c81f8fbdfbc6344ac', 'd7f05d3fa5e745e02e1de41821ccccaf', '3f6de1b4476dd77b42adab8cd8c1f7b5', 'eea5e428ce63804f9b12d21c97b5968f', 'a02ee61542caae25f8a44c9428d30247', '45a786658d3f69717652fed471d03ee0', '5e9af52e463636f704e64c2d28f9cf40', 'f4dce9e69b168fd397255264e85c70c1', 'f757b90ecd5f55f7b1ade82d878e16b6', '525a20c6e9d2c58c085c42198308a3bf', '307ea9788c984e3f6117799c4c20a70b', 'c99710c67fcf9434edaee80b6e07bfab', 'ce931021e18f385f519e945a8a10548e', '121b131eaa369d8f58dacc5c39a77d80');

COMMIT;