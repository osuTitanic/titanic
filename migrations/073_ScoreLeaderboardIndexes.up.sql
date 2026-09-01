CREATE INDEX CONCURRENTLY idx_scores_ranked_visible_order
ON scores (
    beatmap_id,
    mode,
    total_score DESC,
    submitted_at ASC,
    id ASC
)
WHERE
    status_score = 3
    AND hidden = FALSE;