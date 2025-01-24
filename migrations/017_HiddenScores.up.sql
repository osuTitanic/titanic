-- This migration aims to deprecate the previous way of how scores were hidden, which was
-- setting the score status to '-1'. Instead we directly add a new hidden column to the scores table.
ALTER TABLE scores ADD COLUMN hidden BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE scores SET hidden = TRUE WHERE status = -1;
