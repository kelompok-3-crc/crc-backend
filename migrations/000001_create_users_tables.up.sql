CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(60) NOT NULL UNIQUE,
    nip VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_users_nama UNIQUE (nama),
    CONSTRAINT uq_users_nip UNIQUE (nip)
);
