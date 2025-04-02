-- This column is used to determine if the name should be
-- reserved, i.e. if it can be used in registration
ALTER TABLE name_history ADD COLUMN reserved BOOLEAN NOT NULL DEFAULT 'true';