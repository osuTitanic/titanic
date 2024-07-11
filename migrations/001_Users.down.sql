DELETE FROM groups_entries
WHERE user_id = 1 AND group_id IN (1, 999, 1000) OR
      user_id = 2 AND group_id IN (1, 999, 1000);

DELETE FROM groups
WHERE id IN (
    1, 2, 3, 4, 5, 6, 7, 8,
    9, 997, 998, 999, 1000
);

DELETE FROM stats
WHERE id in (1, 2);

DELETE FROM users
WHERE name = 'BanchoBot' OR name = 'peppy';

DELETE FROM user_count;

DROP INDEX IF EXISTS users_name_idx;
DROP INDEX IF EXISTS users_id_idx;
DROP INDEX IF EXISTS stats_id_idx;

DROP TABLE IF EXISTS screenshots;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS verifications;
DROP TABLE IF EXISTS infringements;
DROP TABLE IF EXISTS relationships;
DROP TABLE IF EXISTS user_count;
DROP TABLE IF EXISTS groups_entries;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS stats;
DROP TABLE IF EXISTS users;