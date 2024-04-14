CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users
(
    id UUID NOT NULL DEFAULT gen_random_uuid(), 
    login text NOT NULL UNIQUE,
    passhash bytea NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS records
(
    id UUID NOT NULL DEFAULT gen_random_uuid(), 
	user_id UUID NOT NULL REFERENCES users(id),
    info bytea NOT NULL,
    data_type text NOT NULL,
    meta text NOT NULL,
    version int NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT records_pkey PRIMARY KEY (id)
);
