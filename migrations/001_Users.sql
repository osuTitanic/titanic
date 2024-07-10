INSERT INTO users (name, safe_name, email, pw, country, activated, bot)
VALUES ('BanchoBot', 'banchobot', 'bot@example.com', '------------------------------------------------------------', 'OC', true, true),
       ('peppy', 'peppy', 'pe@ppy.sh', '$2b$12$W5ppLwlSEJ3rpJQRq8UcX.QA5cTm7HvsVpn6MXQHE/6OEO.Iv4DGW', 'AU', true, false);

-- Default stats for BanchoBot
INSERT INTO stats (id, mode)
VALUES (1, 0),
       (1, 1),
       (1, 2),
       (1, 3);

INSERT INTO user_count (count)
VALUES (0);