-- +goose Up

CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name    TEXT NOT NULL DEFAULT '',
    telegram_id   TEXT UNIQUE,
    fd_balance    INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE flowers (
    id         SERIAL PRIMARY KEY,
    season     INTEGER NOT NULL DEFAULT 1,
    image_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_flowers (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    flower_id       INTEGER NOT NULL REFERENCES flowers(id),
    day             INTEGER NOT NULL DEFAULT 1,
    last_watered_at DATE,
    is_dried        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE waterings (
    id                   SERIAL PRIMARY KEY,
    user_flower_id       INTEGER NOT NULL REFERENCES user_flowers(id),
    watered_by_user_id   INTEGER NOT NULL REFERENCES users(id),
    watered_date         DATE NOT NULL,
    UNIQUE (user_flower_id, watered_by_user_id, watered_date)
);

CREATE TABLE seeds (
    id        SERIAL PRIMARY KEY,
    user_id   INTEGER NOT NULL REFERENCES users(id),
    flower_id INTEGER NOT NULL REFERENCES flowers(id),
    quantity  INTEGER NOT NULL DEFAULT 0,
    UNIQUE (user_id, flower_id)
);

CREATE TABLE notifications (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id),
    type       TEXT NOT NULL,
    payload    JSONB NOT NULL DEFAULT '{}',
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_flowers_user_id ON user_flowers(user_id);
CREATE INDEX idx_waterings_user_flower_date ON waterings(user_flower_id, watered_date);
CREATE INDEX idx_notifications_user_id ON notifications(user_id, is_read);

-- +goose Down

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS seeds;
DROP TABLE IF EXISTS waterings;
DROP TABLE IF EXISTS user_flowers;
DROP TABLE IF EXISTS flowers;
DROP TABLE IF EXISTS users;
