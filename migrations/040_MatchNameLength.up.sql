-- Allow mp matches to have longer names
ALTER TABLE mp_matches ALTER COLUMN name TYPE varchar(128);