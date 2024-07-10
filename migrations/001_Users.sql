INSERT INTO users (name, safe_name, email, pw, country, activated, bot)
VALUES ('BanchoBot', 'banchobot', 'bot@example.com', '------------------------------------------------------------', 'OC', true, true),
       ('peppy', 'peppy', 'pe@ppy.sh', '$2b$12$W5ppLwlSEJ3rpJQRq8UcX.QA5cTm7HvsVpn6MXQHE/6OEO.Iv4DGW', 'AU', true, false);

-- Default stats for BanchoBot
INSERT INTO stats (id, mode)
VALUES (1, 0),
       (1, 1),
       (1, 2),
       (1, 3);

INSERT INTO channels (name, topic, read_permissions, write_permissions)
VALUES ('#osu', 'General discussion.', 1, 1),
       ('#announce', 'Public announcements.', 1, 8),
       ('#lobby', 'Multiplayer lobby discussion room.', 1, 1),
       ('#admin', 'General discussion for administrators.', 16, 16);

INSERT INTO groups (id, bancho_permissions, name, short_name, description, color, hidden)
VALUES (1, '16', 'Admins', 'ADMIN', 'Some cool people.', '#9d6b15', false),
       (2, '8', 'Developers', 'DEV', 'These people actively contribute to the titanic codebase.', '#69159d', false),
       (3, '2', 'Beatmap Approval Team', 'BAT', 'These people handle map requests from users and decide which should be approved.', '#1959b8', false),
       (4, '8', 'Global Moderator Team', 'GMT', 'These people focus on player moderation and ensuring the community is a safe place for everyone.', '#59c51b', false),
       (5, '32', 'Tournament Manager Team', 'TMT', 'These people help organize & referee tournament.', '#2117ab', false),
       (6, 4, 'Donators', 'DONOR', 'These people have donated money to this project and are helping to keep it alive. Eternal thanks to them!', '#cf9c02', false),
       (7, 0, 'Alumni', 'ALM', 'These people have voluntary contributed to this project in some way.', '#c51b54', false),
       (8, 0, 'Playtesters', 'TESTER', 'These people were one of the first to register on this server and have helped us finding bugs and giving feedback to the project.', '#31b6e3', false),
       (9, 0, 'Bots', 'BOT', '*beep boop*', '#00b061', false),
       (997, 0, 'Preview', 'PREVIEW', 'These people can access unverified builds of the game.', '#000000', true),
       (998, '1', 'Verified', 'VERIFIED', 'Verified players.', '#000000', true),
       (999, '4', 'Supporter', 'DIRECT', 'People with access to osu! direct.', '#000000', true),
       (1000, '1', 'Players', 'PLAYER', 'People who play the game.', '#000000', true);

INSERT INTO groups_entries (user_id, group_id)
VALUES (1, 1),    -- BanchoBot -> Admins
       (1, 999),  -- BanchoBot -> Supporter
       (1, 1000), -- BanchoBot -> Players
       (2, 1),    -- peppy -> Admins
       (2, 999),  -- peppy -> Supporter
       (2, 1000); -- peppy -> Players

INSERT INTO clients_verified ("type", "hash")
VALUES (0, 'b4ec3c4334a0249dae95c284ec5983df'), -- "runningunderwine"
       (0, '74be16979710d4c4e7c6647856088456'), -- ""
       (0, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (1, 'ad921d60486366258809553a3db49a4a'), -- "unknown"
       (1, '74be16979710d4c4e7c6647856088456'), -- ""
       (1, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (2, 'ad921d60486366258809553a3db49a4a'), -- "unknown"
       (2, 'dcfcd07e645d245babe887e5e2daa016'), -- "0"
       (2, '28c8edde3d61a0411511d3b1866f0636'), -- "1"
       (2, '74be16979710d4c4e7c6647856088456'), -- ""
       (2, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (2, 'd1c651c36f499849f1c9a5843567e686'); -- toshiba hdd

INSERT INTO user_count (count)
VALUES (0);