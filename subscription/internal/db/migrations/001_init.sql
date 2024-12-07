-- Init

CREATE TABLE subscriptions
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER      NOT NULL,
    prefix         VARCHAR(255) NOT NULL,
    files_owner_id INTEGER,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE EXTENSION pg_trgm;
CREATE INDEX idx_prefix ON subscriptions USING GIN (prefix gin_trgm_ops);
CREATE UNIQUE INDEX idx_user_id_prefix ON subscriptions (user_id, prefix);

---- create above / drop below ----
DROP TABLE subscriptions;
DROP EXTENSION pg_trgm;