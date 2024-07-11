CREATE TABLE IF NOT EXISTS releases
(
    name character varying(64) NOT NULL PRIMARY KEY,
    version int NOT NULL,
    description text NOT NULL,
    known_bugs text,
    supported boolean NOT NULL DEFAULT true,
    recommended boolean NOT NULL DEFAULT false,
    preview boolean NOT NULL DEFAULT false,
    downloads varchar[] NOT NULL DEFAULT '{}',
    hashes jsonb NOT NULL DEFAULT '[]',
    screenshots jsonb NOT NULL DEFAULT '[]',
    actions jsonb NOT NULL DEFAULT '[]',
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS clients
(
    user_id int NOT NULL REFERENCES users (id),
    executable character(32) NOT NULL,
    adapters character(32) NOT NULL,
    unique_id character(32) NOT NULL,
    disk_signature character(32) NOT NULL,
    banned boolean NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, executable, adapters, unique_id, disk_signature)
);

-- Table for verified hardware ids, to
-- bypass multiaccounting checks.
-- Types:
-- 0: Adapters
-- 1: Unique Id
-- 2: Disk Signature
CREATE TABLE IF NOT EXISTS clients_verified
(
    "type" smallint NOT NULL,
    "hash" character(32) NOT NULL,
    PRIMARY KEY ("type", "hash")
);

INSERT INTO clients_verified ("type", "hash")
VALUES (0, 'b4ec3c4334a0249dae95c284ec5983df'), -- "runningunderwine"
       (0, '74be16979710d4c4e7c6647856088456'), -- ""
       (0, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (1, 'ad921d60486366258809553a3db49a4a'), -- "unknown"
       (1, '74be16979710d4c4e7c6647856088456'), -- ""
       (1, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (2, 'ad921d60486366258809553a3db49a4a'), -- "unknown"
       (2, 'dcfcd07e645d245babe887e5e2daa016'), -- "0"
       (2, '28c8edde3d61a0411511d3b1866f0636'), -- "1"
       (2, '74be16979710d4c4e7c6647856088456'), -- ""
       (2, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (2, 'd1c651c36f499849f1c9a5843567e686'); -- toshiba hdd