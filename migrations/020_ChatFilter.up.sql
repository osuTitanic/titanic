-- Implementation for the chat filter feature, similar to how bancho does it.
-- You are able to directly block messages or replace them with something else.
-- e.g. "fuck" -> "firetruck".
-- Additionally you are able to timeout the user through the `block_timeout_duration` column.
CREATE TABLE filters (
    name text NOT NULL PRIMARY KEY,
    pattern text NOT NULL,
    replacement text DEFAULT NULL,
    block boolean NOT NULL DEFAULT false,
    block_timeout_duration int DEFAULT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);