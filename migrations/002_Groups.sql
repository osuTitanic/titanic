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