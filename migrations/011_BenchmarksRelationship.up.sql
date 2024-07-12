ALTER TABLE benchmarks
    ALTER COLUMN user_id TYPE integer;

ALTER TABLE benchmarks
	ADD CONSTRAINT users_id_fkey
	FOREIGN KEY (user_id)
	REFERENCES users (id);