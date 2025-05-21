CREATE TABLE user_permissions (
    id serial PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission varchar(255) NOT NULL,
    rejected boolean DEFAULT false,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_permissions (
    id serial PRIMARY KEY,
    group_id int NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission varchar(255) NOT NULL,
    rejected boolean DEFAULT false,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_permissions_user_id ON user_permissions (user_id);
CREATE INDEX idx_group_permissions_group_id ON group_permissions (group_id);
CREATE INDEX idx_user_permissions_user_id_rejected ON user_permissions (user_id, rejected);
CREATE INDEX idx_group_permissions_group_id_rejected ON group_permissions (group_id, rejected);

CREATE INDEX idx_groups_entries_user_id ON groups_entries (user_id, group_id);
CREATE INDEX idx_user_perms_u_r_p ON user_permissions(user_id, rejected) INCLUDE (permission);
CREATE INDEX idx_group_perms_g_r_p ON group_permissions (group_id, rejected) INCLUDE (permission);

INSERT INTO group_permissions (group_id, permission, rejected) VALUES
    -- Admins
    (1, '*', false),
    -- Developers
    (2, 'commands.moderation.*', false),
    (2, 'users.moderation.*', false),
    (2, 'forum.moderation.*', false),
    (2, 'chat.moderation.*', false),
    (2, 'beatmaps.*', false),
    -- BAT
    (3, 'beatmaps.*', false),
    (3, 'forum.topics.lock', false),
    (3, 'forum.posts.lock', false),
    (3, 'forum.posts.edit_locked', false),
    (3, 'forum.beatmaps.edit_icon', false),
    (3, 'forum.kudosu.*', false),
    -- Moderators
    (4, 'commands.moderation.*', false),
    (4, 'users.moderation.*', false),
    (4, 'forum.moderation.*', false),
    (4, 'chat.moderation.*', false),
    (4, 'beatmaps.nuke', false),
    -- Tournament Manager Team
    (5, 'commands.tournaments.create_match', false),
    (5, 'commands.tournaments.force_invite', false),
    (5, 'commands.tournaments.addref', false),
    (5, 'commands.tournaments.remref', false),
    -- Donator
    (6, 'clients.preview', false),
    -- Bots
    (9, 'chat.bypass_filter', false),
    (9, 'chat.bypass_spam', false),
    -- Preview
    (997, 'clients.validation.bypass', false),
    -- Bypass
    (998, 'scores.validation.bypass', false),
    -- "Supporter"
    (999, 'beatmaps.direct.search', false),
    (999, 'beatmaps.direct.download', false),
    (999, 'beatmaps.leaderboards.country', false),
    (999, 'beatmaps.leaderboards.friends', false),
    (999, 'beatmaps.leaderboards.mods', false),
    (999, 'beatmaps.favourites.extended_limit', false),
    (999, 'beatmaps.upload.extended_limit', false),
    (999, 'users.friends.extended_limit', false),
    -- Regular
    (1000, 'users.notifications.delete', false),
    (1000, 'users.profile.update', false),
    (1000, 'users.friends.delete', false),
    (1000, 'users.friends.add', false),
    (1000, 'beatmaps.search', false),
    (1000, 'beatmaps.upload', false),
    (1000, 'beatmaps.revive', false),
    (1000, 'beatmaps.delete', false),
    (1000, 'beatmaps.download', false),
    (1000, 'beatmaps.rating.submit', false),
    (1000, 'beatmaps.rating.view', false),
    (1000, 'beatmaps.favourites.create', false),
    (1000, 'beatmaps.favourites.delete', false),
    (1000, 'beatmaps.comments.view', false),
    (1000, 'beatmaps.comments.submit', false),
    (1000, 'beatmaps.leaderboards.global', false),
    (1000, 'commands.standard.*', false),
    (1000, 'commands.multiplayer.standard.*', false),
    (1000, 'benchmarks.view', false),
    (1000, 'benchmarks.submit', false),
    (1000, 'chat.search', false),
    (1000, 'chat.messages.view', false),
    (1000, 'chat.messages.create', false),
    (1000, 'chat.messages.private.view', false),
    (1000, 'chat.messages.private.create', false),
    (1000, 'forum.search', false),
    (1000, 'forum.topics.create', false),
    (1000, 'forum.posts.create', false),
    (1000, 'forum.posts.edit', false),
    (1000, 'forum.posts.delete', false),
    (1000, 'forum.bookmarks.create', false),
    (1000, 'forum.bookmarks.delete', false),
    (1000, 'forum.subscriptions.create', false),
    (1000, 'forum.subscriptions.delete', false),
    (1000, 'forum.kudosu.reward', false),
    (1000, 'scores.submit', false),
    (1000, 'scores.download', false),
    (1000, 'scores.pins.create', false),
    (1000, 'scores.pins.delete', false),
    (1000, 'bancho.login', false),
    (1000, 'bancho.coins.earn', false),
    (1000, 'bancho.coins.use', false),
    (1000, 'bancho.coins.recharge', false),
    (1000, 'bancho.errors.submit', false),
    (1000, 'bancho.screenshots.upload', false),
    (1000, 'bancho.matches.create', false),
    (1000, 'bancho.matches.join', false),
    (1000, 'bancho.spectating.start', false),
    (1000, 'bancho.spectating.send_frames', false);