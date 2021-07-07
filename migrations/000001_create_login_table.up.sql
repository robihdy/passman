CREATE TABLE IF NOT EXISTS logins (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    username text NOT NULL,
    password bytea NOT NULL,
    website text,
    version integer NOT NULL DEFAULT 1
);