-- These are a bunch of indexes I have applied over the years to the prod db
-- Never actually added them to migrations, so this migration is here to fix that

-- beatmap_modding
CREATE INDEX IF NOT EXISTS beatmap_modding_sender_id_id_idx ON beatmap_modding USING btree (sender_id, id DESC);
CREATE INDEX IF NOT EXISTS beatmap_modding_target_id_id_idx ON beatmap_modding USING btree (target_id, id DESC);
CREATE INDEX IF NOT EXISTS idx_beatmap_modding_sender_target ON beatmap_modding USING btree (sender_id, target_id);
CREATE INDEX IF NOT EXISTS idx_modding_post_id ON beatmap_modding USING btree (post_id);
CREATE INDEX IF NOT EXISTS idx_modding_set_id ON beatmap_modding USING btree (set_id);

-- beatmap_nominations
CREATE INDEX IF NOT EXISTS idx_beatmap_nominations_user_time ON beatmap_nominations USING btree (user_id, "time" DESC);
CREATE INDEX IF NOT EXISTS idx_nominations_set_id ON beatmap_nominations USING btree (set_id);

-- beatmaps
CREATE INDEX IF NOT EXISTS idx_beatmaps_difficulty ON beatmaps USING btree (diff DESC);
CREATE INDEX IF NOT EXISTS idx_beatmaps_id_desc ON beatmaps USING btree (id DESC);
CREATE INDEX IF NOT EXISTS idx_beatmaps_mode ON beatmaps USING btree (mode);
CREATE INDEX IF NOT EXISTS idx_beatmaps_playcount_desc ON beatmaps USING btree (playcount DESC);
CREATE INDEX IF NOT EXISTS idx_beatmaps_set_id_playcount_include ON beatmaps USING btree (set_id) INCLUDE (playcount);
CREATE INDEX IF NOT EXISTS idx_beatmaps_set_id_playcount_with_lb ON beatmaps USING btree (set_id, playcount DESC) WHERE (status > 0);
CREATE INDEX IF NOT EXISTS idx_beatmaps_set_id_playcount_without_lb ON beatmaps USING btree (set_id, playcount DESC) WHERE (status <= 0);
CREATE INDEX IF NOT EXISTS idx_beatmaps_set_mode ON beatmaps USING btree (set_id, mode);
CREATE INDEX IF NOT EXISTS idx_beatmaps_status_id_alt ON beatmaps USING btree (status, id);

-- beatmapsets
CREATE INDEX IF NOT EXISTS beatmapsets_approved_date_idx ON beatmapsets USING btree (approved_date DESC) WHERE (submission_status <> '-3'::integer);
CREATE INDEX IF NOT EXISTS idx_beatmapset_text_search ON beatmapsets USING gin (to_tsvector('simple'::regconfig, (((((((((title)::text || ' '::text) || (artist)::text) || ' '::text) || (creator)::text) || ' '::text) || (source)::text) || ' '::text) || (tags)::text)));
CREATE INDEX IF NOT EXISTS idx_beatmapsets_approved_date_partial ON beatmapsets USING btree (approved_date DESC) WHERE (submission_status > '-3'::integer);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_approved_date_with_lb ON beatmapsets USING btree (approved_date DESC) WHERE (submission_status > 0);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_artist_desc ON beatmapsets USING btree (artist DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_artist_with_lb ON beatmapsets USING btree (artist DESC) WHERE (submission_status > 0);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_artist_without_lb ON beatmapsets USING btree (artist DESC) WHERE (submission_status <= 0);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_created_with_lb ON beatmapsets USING btree (submission_date DESC) WHERE (submission_status > 0);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_creator_id ON beatmapsets USING btree (creator_id);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_download_server_fastpath ON beatmapsets USING btree (id) INCLUDE (download_server);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_enhanced ON beatmapsets USING btree (enhanced);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_favourite_count_desc ON beatmapsets USING btree (favourite_count DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_genre ON beatmapsets USING btree (genre_id);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_id_desc ON beatmapsets USING btree (id DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_language ON beatmapsets USING btree (language_id);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_last_update ON beatmapsets USING btree (last_updated DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_rating_avg_desc ON beatmapsets USING btree (rating_average DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_rating_count_desc ON beatmapsets USING btree (rating_count DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_search_text_gist ON beatmapsets USING gist (search_text gist_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_server_private ON beatmapsets USING btree (server) WHERE (server = 1);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_status_approved_id ON beatmapsets USING btree (submission_status, approved_date DESC, id);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_status_genre_lang_server ON beatmapsets USING btree (submission_status, genre_id, language_id, server);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_storyboard ON beatmapsets USING btree (has_storyboard) WHERE (has_storyboard = true);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_submission_date ON beatmapsets USING btree (submission_date DESC);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_submission_status_id ON beatmapsets USING btree (submission_status, id);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_title_graveyard ON beatmapsets USING btree (title DESC) WHERE (submission_status = '-2'::integer);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_title_with_lb ON beatmapsets USING btree (title DESC) WHERE (submission_status > 0);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_title_without_lb ON beatmapsets USING btree (title DESC) WHERE (submission_status <= 0);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_topic_id ON beatmapsets USING btree (topic_id);
CREATE INDEX IF NOT EXISTS idx_beatmapsets_video ON beatmapsets USING btree (has_video) WHERE (has_video = true);

-- clients
CREATE INDEX IF NOT EXISTS idx_clients_hwid ON clients USING btree (adapters, unique_id, disk_signature);

-- favourites
CREATE INDEX IF NOT EXISTS idx_favourites_set_id_desc ON favourites USING btree (set_id DESC);

-- forum_posts
CREATE INDEX IF NOT EXISTS forum_posts_last_visible_by_forum ON forum_posts USING btree (forum_id, id) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS idx_forum_posts_forum_hidden_id_desc ON forum_posts USING btree (forum_id, hidden, id DESC);
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic_hidden_id ON forum_posts USING btree (topic_id, hidden, id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic_hidden_id_desc ON forum_posts USING btree (topic_id, hidden, id DESC);
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic_id_asc ON forum_posts USING btree (topic_id, id) WHERE ((deleted = false) AND (draft = false) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic_id_desc ON forum_posts USING btree (topic_id, id DESC) WHERE ((deleted = false) AND (draft = false) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic_id_fastpath ON forum_posts USING btree (topic_id) INCLUDE (id) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic_lookup ON forum_posts USING btree (topic_id, hidden, id) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS idx_forum_posts_user_topic_draft_id_desc ON forum_posts USING btree (user_id, topic_id, draft, id DESC);
CREATE INDEX IF NOT EXISTS idx_posts_deleted ON forum_posts USING btree (deleted);
CREATE INDEX IF NOT EXISTS idx_posts_topic_id_id ON forum_posts USING btree (topic_id, id);
CREATE INDEX IF NOT EXISTS ix_forum_posts_forum_visible_id ON forum_posts USING btree (forum_id, id DESC) WHERE (hidden = false);

-- forum_subscribers
CREATE INDEX IF NOT EXISTS idx_subscribers_topic_user ON forum_subscribers USING btree (topic_id, user_id);

-- forum_topics
CREATE INDEX IF NOT EXISTS idx_forum_topics_announcement ON forum_topics USING btree (announcement, hidden, id DESC) WHERE ((announcement = true) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_forum_topics_forum_announcement_visible ON forum_topics USING btree (forum_id, id DESC) WHERE ((announcement = true) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_forum_topics_forum_hidden_lastpost ON forum_topics USING btree (forum_id, hidden, last_post_at DESC);
CREATE INDEX IF NOT EXISTS idx_forum_topics_forum_pinned_visible ON forum_topics USING btree (forum_id, id DESC) WHERE ((pinned = true) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_topics_announcements ON forum_topics USING btree (announcement);
CREATE INDEX IF NOT EXISTS idx_topics_forum_id_announcement ON forum_topics USING btree (forum_id, announcement) WHERE (announcement = true);
CREATE INDEX IF NOT EXISTS idx_topics_forum_id_last_post ON forum_topics USING btree (forum_id, last_post_at DESC);
CREATE INDEX IF NOT EXISTS idx_topics_forum_id_pinned ON forum_topics USING btree (forum_id, pinned) WHERE (pinned = true);
CREATE INDEX IF NOT EXISTS idx_topics_hidden ON forum_topics USING btree (hidden);
CREATE INDEX IF NOT EXISTS idx_topics_news ON forum_topics USING btree (id DESC) WHERE ((hidden = false) AND (announcement = true));
CREATE INDEX IF NOT EXISTS ix_forum_topics_forum_visible ON forum_topics USING btree (forum_id) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS ix_forum_topics_recent_visible ON forum_topics USING btree (forum_id, last_post_at DESC) WHERE (hidden = false);

-- forums
CREATE INDEX IF NOT EXISTS idx_forums_hidden ON forums USING btree (hidden);
CREATE INDEX IF NOT EXISTS ix_forums_parent_visible ON forums USING btree (parent_id) WHERE (hidden = false);

-- group_permissions
CREATE INDEX IF NOT EXISTS idx_group_permissions_lookup ON group_permissions USING btree (group_id, permission, rejected);

-- groups
CREATE INDEX IF NOT EXISTS idx_groups_hidden ON groups USING btree (hidden);

-- groups_entries
CREATE INDEX IF NOT EXISTS ix_groups_entries_user_id_included ON groups_entries USING btree (user_id) INCLUDE (group_id);

-- infringements
CREATE INDEX IF NOT EXISTS idx_infringements_user_action ON infringements USING btree (user_id, action);
CREATE INDEX IF NOT EXISTS idx_infringements_user_id_id_desc ON infringements USING btree (user_id, id DESC);
CREATE INDEX IF NOT EXISTS idx_infringements_user_id_time_desc ON infringements USING btree (user_id, "time" DESC);

-- messages
CREATE INDEX IF NOT EXISTS idx_messages_osuchat_id_desc ON messages USING btree (id DESC) WHERE ((target)::text = '#osu'::text);
CREATE INDEX IF NOT EXISTS idx_messages_target_id_desc ON messages USING btree (target, id DESC);
CREATE INDEX IF NOT EXISTS idx_messages_target_time_desc ON messages USING btree (target, "time" DESC);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id_desc ON messages USING btree (sender, id DESC);

-- mp_events
CREATE INDEX IF NOT EXISTS idx_events_match_id_time ON mp_events USING btree (match_id, "time" DESC);
CREATE INDEX IF NOT EXISTS idx_events_match_id_type ON mp_events USING btree (match_id, type);
CREATE INDEX IF NOT EXISTS idx_events_time ON mp_events USING btree ("time" DESC);

-- name_history
CREATE INDEX IF NOT EXISTS idx_name_history_trgm ON name_history USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_name_history_user_id ON name_history USING btree (user_id);

-- notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id_time ON notifications USING btree (user_id, "time" DESC);

-- plays
CREATE INDEX IF NOT EXISTS idx_plays_beatmap_count ON plays USING btree (beatmap_id, count);
CREATE INDEX IF NOT EXISTS idx_plays_beatmap_count_desc ON plays USING btree (beatmap_id, count DESC);
CREATE INDEX IF NOT EXISTS idx_plays_beatmap_id_set_id ON plays USING btree (beatmap_id, set_id);
CREATE INDEX IF NOT EXISTS idx_plays_beatmap_id_set_id_count ON plays USING btree (beatmap_id, set_id) INCLUDE (count);
CREATE INDEX IF NOT EXISTS idx_plays_beatmap_user ON plays USING btree (beatmap_id, user_id);
CREATE INDEX IF NOT EXISTS idx_plays_count_desc ON plays USING btree (count DESC);
CREATE INDEX IF NOT EXISTS idx_plays_set_id ON plays USING btree (set_id);
CREATE INDEX IF NOT EXISTS idx_plays_user_count ON plays USING btree (user_id, count DESC);

-- profile_activity
CREATE INDEX IF NOT EXISTS idx_activity_time ON profile_activity USING btree ("time" DESC);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_no_mode ON profile_activity USING btree (user_id) WHERE (mode = NULL::smallint);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_time_catch ON profile_activity USING btree (user_id, "time" DESC) WHERE (mode = 2);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_time_mania ON profile_activity USING btree (user_id, "time" DESC) WHERE (mode = 3);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_time_no_mode ON profile_activity USING btree (user_id, "time" DESC) WHERE (mode = NULL::smallint);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_time_not_hidden ON profile_activity USING btree (user_id, "time" DESC) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_time_osu ON profile_activity USING btree (user_id, "time" DESC) WHERE (mode = 0);
CREATE INDEX IF NOT EXISTS idx_activity_user_id_time_taiko ON profile_activity USING btree (user_id, "time" DESC) WHERE (mode = 1);
CREATE INDEX IF NOT EXISTS idx_profile_activity_hidden ON profile_activity USING btree (hidden);
CREATE INDEX IF NOT EXISTS idx_profile_activity_user_hidden_id_desc ON profile_activity USING btree (user_id, hidden, id DESC);
CREATE INDEX IF NOT EXISTS idx_profile_activity_user_mode_hidden_time ON profile_activity USING btree (user_id, mode, hidden, "time" DESC);
CREATE INDEX IF NOT EXISTS idx_profile_activity_user_mode_time ON profile_activity USING btree (user_id, mode, "time" DESC);
CREATE INDEX IF NOT EXISTS idx_profile_activity_user_time ON profile_activity USING btree (user_id, "time" DESC);
CREATE INDEX IF NOT EXISTS profile_activity_user_mode_visible_id_idx ON profile_activity USING btree (user_id, mode, id DESC) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS profile_activity_user_mode_visible_time_idx ON profile_activity USING btree (user_id, mode, "time" DESC, id DESC) WHERE (hidden = false);

-- profile_badges
CREATE INDEX IF NOT EXISTS idx_badges_user_id ON profile_badges USING btree (user_id);

-- profile_rank_history
CREATE INDEX IF NOT EXISTS idx_peak_global_rank_ctb ON profile_rank_history USING btree (user_id, global_rank DESC) WHERE (mode = 2);
CREATE INDEX IF NOT EXISTS idx_peak_global_rank_mania ON profile_rank_history USING btree (user_id, global_rank DESC) WHERE (mode = 3);
CREATE INDEX IF NOT EXISTS idx_peak_global_rank_std ON profile_rank_history USING btree (user_id, global_rank DESC) WHERE (mode = 0);
CREATE INDEX IF NOT EXISTS idx_peak_global_rank_taiko ON profile_rank_history USING btree (user_id, global_rank DESC) WHERE (mode = 1);
CREATE INDEX IF NOT EXISTS idx_rank_history_user_id_time ON profile_rank_history USING btree (user_id, "time" DESC);

-- ratings
CREATE INDEX IF NOT EXISTS idx_ratings_map_checksum ON ratings USING btree (map_checksum);
CREATE INDEX IF NOT EXISTS idx_ratings_rating_asc ON ratings USING btree (rating);
CREATE INDEX IF NOT EXISTS idx_ratings_set_id_rating ON ratings USING btree (set_id, rating);

-- relationships
CREATE INDEX IF NOT EXISTS idx_relationships_target_id ON relationships USING btree (target_id);
CREATE INDEX IF NOT EXISTS idx_relationships_target_id_blocked ON relationships USING btree (target_id) WHERE (status = 1);
CREATE INDEX IF NOT EXISTS idx_relationships_target_id_friends ON relationships USING btree (target_id) WHERE (status = 0);
CREATE INDEX IF NOT EXISTS idx_relationships_user_id_blocked ON relationships USING btree (user_id) WHERE (status = 1);
CREATE INDEX IF NOT EXISTS idx_relationships_user_id_friends ON relationships USING btree (user_id) WHERE (status = 0);
CREATE INDEX IF NOT EXISTS idx_relationships_user_id_target_id_blocked ON relationships USING btree (user_id, target_id) WHERE (status = 1);
CREATE INDEX IF NOT EXISTS idx_relationships_user_id_target_id_friends ON relationships USING btree (user_id, target_id) WHERE (status = 0);
CREATE INDEX IF NOT EXISTS idx_relationships_user_status_target ON relationships USING btree (user_id, status, target_id);

-- releases_titanic
CREATE INDEX IF NOT EXISTS idx_releases_version ON releases_titanic USING btree (version);

-- resource_mirrors
CREATE INDEX IF NOT EXISTS idx_resource_mirrors_priority ON resource_mirrors USING btree (priority);
CREATE INDEX IF NOT EXISTS idx_resource_mirrors_server_priority ON resource_mirrors USING btree (server, priority);
CREATE INDEX IF NOT EXISTS idx_resource_mirrors_type_server_priority ON resource_mirrors USING btree (type, server, priority);
CREATE INDEX IF NOT EXISTS idx_resource_mirrors_type_server_priority_cover ON resource_mirrors USING btree (type, server, priority) INCLUDE (url);
CREATE INDEX IF NOT EXISTS idx_resource_mirrors_type_url_server_priority ON resource_mirrors USING btree (type, url, server, priority);

-- scores
CREATE INDEX IF NOT EXISTS idx_beatmap_mode_status_score_total_score_id ON scores USING btree (beatmap_id, mode, status_score, total_score DESC, id);
CREATE INDEX IF NOT EXISTS idx_score_user_status_beatmap ON scores USING btree (user_id, status, beatmap_id);
CREATE INDEX IF NOT EXISTS idx_scores_active_beatmap ON scores USING btree (beatmap_id, submitted_at DESC) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS idx_scores_agg ON scores USING btree (beatmap_id, total_score) WHERE ((mode = 0) AND (status_score = 3) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_beatmap_id_user_id_mode_status_score_hidden_grade ON scores USING btree (beatmap_id, user_id, mode, status_score, hidden) INCLUDE (grade);
CREATE INDEX IF NOT EXISTS idx_scores_beatmap_leader_lookup ON scores USING btree (beatmap_id, mode, total_score DESC) WHERE ((hidden = false) AND (status_score = 2));
CREATE INDEX IF NOT EXISTS idx_scores_beatmap_mode_hidden_tscore_id ON scores USING btree (beatmap_id, mode, hidden, total_score DESC, id);
CREATE INDEX IF NOT EXISTS idx_scores_beatmap_mode_status_hidden_tscore ON scores USING btree (beatmap_id, mode, status_score, hidden, total_score DESC, id);
CREATE INDEX IF NOT EXISTS idx_scores_best_for_query ON scores USING btree (mode, beatmap_id, total_score DESC, status_score, hidden, user_id, id DESC);
CREATE INDEX IF NOT EXISTS idx_scores_best_scores_by_pp ON scores USING btree (beatmap_id, pp DESC) WHERE ((mode = 0) AND (status_score = 3) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_best_scores_by_score ON scores USING btree (beatmap_id, total_score DESC) WHERE ((mode = 0) AND (status_score = 3) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_catch_active_pp_desc ON scores USING btree (pp DESC) WHERE ((mode = 2) AND (status > 0) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_mania_active_pp_desc ON scores USING btree (pp DESC) WHERE ((mode = 3) AND (status > 0) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_mode_status_hidden_pp_desc ON scores USING btree (mode, status, hidden, pp DESC);
CREATE INDEX IF NOT EXISTS idx_scores_mode_status_score_hidden_beatmap_id_total_score ON scores USING btree (mode, status_score, hidden, beatmap_id, total_score DESC);
CREATE INDEX IF NOT EXISTS idx_scores_mods_ap ON scores USING btree (((mods & 8192)));
CREATE INDEX IF NOT EXISTS idx_scores_mods_rx ON scores USING btree (((mods & 128)));
CREATE INDEX IF NOT EXISTS idx_scores_mods_status_score ON scores USING btree (mods, status_score);
CREATE INDEX IF NOT EXISTS idx_scores_optimized ON scores USING btree (mode, status_score, hidden, beatmap_id, total_score) INCLUDE (user_id, id);
CREATE INDEX IF NOT EXISTS idx_scores_osu_active_pp_desc ON scores USING btree (pp DESC) WHERE ((mode = 0) AND (status > 0) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_replay_md5 ON scores USING btree (replay_md5);
CREATE INDEX IF NOT EXISTS idx_scores_submitted_at_filtered ON scores USING btree (submitted_at, beatmap_id) WHERE (hidden = false);
CREATE INDEX IF NOT EXISTS idx_scores_submitted_brin ON scores USING brin (submitted_at);
CREATE INDEX IF NOT EXISTS idx_scores_taiko_active_pp_desc ON scores USING btree (pp DESC) WHERE ((mode = 1) AND (status > 0) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_user_id_mode_status_score_hidden_id ON scores USING btree (user_id, mode, status_score, hidden, id DESC);
CREATE INDEX IF NOT EXISTS idx_scores_user_leader_count ON scores USING btree (user_id, mode, beatmap_id, total_score) WHERE ((hidden = false) AND (status_score = 2));
CREATE INDEX IF NOT EXISTS idx_scores_user_mode_grade_active ON scores USING btree (user_id, mode, grade) WHERE ((status_score = 3) AND (hidden IS FALSE) AND ((grade)::text <> 'F'::text));
CREATE INDEX IF NOT EXISTS idx_scores_user_mode_id_desc ON scores USING btree (user_id, mode, id DESC) WHERE ((status >= 0) AND (hidden = false));
CREATE INDEX IF NOT EXISTS idx_scores_user_order ON scores USING btree (user_id, mode, status_score, hidden, id DESC) INCLUDE (beatmap_id, total_score) WHERE ((status_score = 3) AND (hidden = false));
CREATE INDEX IF NOT EXISTS scores_global_std_ranked_visible_idx ON scores USING btree (beatmap_id, total_score DESC) WHERE ((hidden = false) AND (status_score = 3) AND (mode = 0));
CREATE INDEX IF NOT EXISTS scores_user_std_ranked_visible_idx ON scores USING btree (user_id, beatmap_id, total_score DESC) WHERE ((hidden = false) AND (status_score = 3) AND (mode = 0));

-- stats
CREATE INDEX IF NOT EXISTS idx_stats_mode ON stats USING btree (mode);
CREATE INDEX IF NOT EXISTS idx_stats_user_mode_pp ON stats USING btree (id, mode, pp DESC);

-- user_activity
CREATE INDEX IF NOT EXISTS idx_usercount_time ON user_activity USING btree ("time" DESC);

-- user_permissions
CREATE INDEX IF NOT EXISTS idx_user_permissions_lookup ON user_permissions USING btree (user_id, permission, rejected);

-- users
CREATE INDEX IF NOT EXISTS idx_users_activated ON users USING btree (activated);
CREATE INDEX IF NOT EXISTS idx_users_activated_restricted ON users USING btree (activated, restricted) WHERE ((NOT restricted) AND activated);
CREATE INDEX IF NOT EXISTS idx_users_discord_id ON users USING btree (discord_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users USING btree (email);
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON users USING btree (lower((email)::text));
CREATE INDEX IF NOT EXISTS idx_users_id_activated ON users USING btree (id, activated);
CREATE INDEX IF NOT EXISTS idx_users_latest_activity ON users USING btree (latest_activity DESC) WHERE (NOT restricted);
CREATE INDEX IF NOT EXISTS idx_users_name_lower ON users USING btree (lower((name)::text));
CREATE INDEX IF NOT EXISTS idx_users_name_trgm ON users USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_rankings_covering ON users USING btree (id) INCLUDE (name, country);
CREATE INDEX IF NOT EXISTS idx_users_restricted ON users USING btree (restricted);
CREATE INDEX IF NOT EXISTS idx_users_safe_name ON users USING btree (safe_name);

-- verifications
CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications USING btree (user_id);

-- direct_messages
CREATE INDEX IF NOT EXISTS idx_dms_sender_target_id_desc ON direct_messages USING btree (sender_id, target_id, id DESC);
CREATE INDEX IF NOT EXISTS idx_dms_unread_by_target ON direct_messages USING btree (target_id, sender_id) WHERE (read = false);

-- comments
CREATE INDEX IF NOT EXISTS idx_comments_target ON comments USING btree (target_id, target_type, "time");
CREATE INDEX IF NOT EXISTS idx_comments_user_time ON comments USING btree (user_id, "time");
