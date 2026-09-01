CREATE INDEX CONCURRENTLY idx_scores_user_visible_best_order
ON scores (user_id, mode, id DESC)
INCLUDE (beatmap_id, total_score, submitted_at)
WHERE status_score = 3 AND hidden = FALSE;