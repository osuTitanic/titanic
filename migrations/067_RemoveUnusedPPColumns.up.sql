ALTER TABLE stats
    DROP COLUMN IF EXISTS pp_vn,
    DROP COLUMN IF EXISTS pp_rx,
    DROP COLUMN IF EXISTS pp_ap;

ALTER TABLE profile_rank_history
    DROP COLUMN IF EXISTS pp_vn,
    DROP COLUMN IF EXISTS pp_rx,
    DROP COLUMN IF EXISTS pp_ap,
    DROP COLUMN IF EXISTS pp_vn_rank,
    DROP COLUMN IF EXISTS pp_rx_rank,
    DROP COLUMN IF EXISTS pp_ap_rank;
