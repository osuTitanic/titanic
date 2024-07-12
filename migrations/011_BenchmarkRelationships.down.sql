ALTER TABLE benchmarks
    DROP CONSTRAINT users_id_fkey;

ALTER TABLE benchmarks
    ALTER COLUMN user_id TYPE bigint;