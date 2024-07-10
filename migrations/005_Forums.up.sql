INSERT INTO forums (id, name)
VALUES (1, 'Titanic!'),
       (2, 'Beatmaps');

INSERT INTO forums (id, name, description, parent_id, allow_icons)
VALUES (3, 'Development', 'Discuss the future of this project.', 1, true),
       (4, 'Gameplay & Rankings', 'Show off your scores to the world and discuss them.', 1, true),
       (5, 'Skinning', 'Discuss & share skins and other customizations.', 1, true),
       (6, 'Feature Requests', 'Suggest what you would like to see in this project.', 1, true),
       (7, 'Support', 'Need help? You will find it here.', 1, true),
       (8, 'Ranked/Approved Beatmaps', 'New approved beatmaps will show up in here.', 2, false),
       (9, 'Pending Beatmaps', 'New pending beatmaps that are waiting for approval.', 2, false),
       (10, 'Work In Progress/Help Wanted', 'Work-in-progress beatmaps that may need support/help.', 2, false),
       (11, 'Map Requests', 'Request beatmaps from the official servers.', 2, false),
       (12, 'Beatmap Graveyard', 'Beatmaps that havent been active for 4 weeks or more will be moved here.', 2, false),
       (13, 'Mapping Discussion', 'Share the art of mapping with others.', 2, false);

INSERT INTO forums (id, name, parent_id, allow_icons)
VALUES (14, 'Taiko', 4, false),
       (15, 'Catch the Beat', 4, false),
       (16, 'osu!mania', 4, false),
       (17, 'Completed Skins', 5, false),
       (18, 'Bug Reports', 3, false);

INSERT INTO forums (id, name)
VALUES (19, 'Other');

INSERT INTO forums (id, name, description, parent_id, allow_icons)
VALUES (20, 'General Discussion', 'The place where you dont post crap.', 19, false),
       (21, 'Off-Topic', 'The perfect place for brainrot.', 19, false),
       (22, 'Introductions', 'Introduce yourself to other passengers.', 19, false),
       (23, 'Client Modding', 'Discover a new world of osu!.', 19, false),
       (24, 'Video Games', 'Discuss any non-osu! games in here.', 19, false),
       (25, 'Art', 'Share your artistic masterpieces, find new avatars and more.', 19, false);

INSERT INTO forum_icons (name, location)
VALUES ('heart', '/images/icons/forum/heart.gif'),
       ('heartpop', '/images/icons/forum/heartpop.gif'),
       ('bubble', '/images/icons/forum/thinking.gif'),
       ('bubblepop', '/images/icons/forum/bubblepop.png'),
       ('fire', '/images/icons/forum/fire.gif'),
       ('star', '/images/icons/forum/star.gif'),
       ('radioactive', '/images/icons/forum/radioactive.gif'),
       ('alert', '/images/icons/forum/alert.gif'),
       ('info', '/images/icons/forum/info.gif'),
       ('question', '/images/icons/forum/question.gif'),
       ('osu', '/images/icons/forum/osu.gif'),
       ('taiko', '/images/icons/forum/taiko.gif'),
       ('ctb', '/images/icons/forum/ctb.gif'),
       ('mania', '/images/icons/forum/mania.gif');