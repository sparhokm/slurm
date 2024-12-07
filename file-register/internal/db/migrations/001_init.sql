-- Init

CREATE TABLE files
(
    id           UUID PRIMARY KEY,
    owner_id     INTEGER      NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    size         BIGINT       NOT NULL,
    filepath     VARCHAR(255) NOT NULL,
    version      SMALLINT     NOT NULL,
    created_at   TIMESTAMP    NOT NULL,
    updated_at   TIMESTAMP
);
CREATE UNIQUE INDEX idx_owner_id_filepath ON files (owner_id, filepath);

CREATE TABLE outbox
(
    id           SERIAL PRIMARY KEY,
    event_type   SMALLINT  NOT NULL,
    request_id   VARCHAR(255),
    trace_id     VARCHAR(255),
    span_id      VARCHAR(255),
    payload      JSONB     NOT NULL,
    created_at   TIMESTAMP NOT NULL,
    processed_at TIMESTAMP
);

---- create above / drop below ----
DROP TABLE files, outbox;