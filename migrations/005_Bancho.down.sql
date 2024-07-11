DELETE FROM channels
WHERE name IN (
    '#osu',
    '#announce',
    '#lobby',
    '#admin'
);

DROP TABLE IF EXISTS mp_events;
DROP TABLE IF EXISTS mp_matches;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS channels;
DROP TABLE IF EXISTS logins;