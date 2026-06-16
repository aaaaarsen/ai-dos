CREATE TABLE chats (
    id          BIGSERIAL                           PRIMARY KEY,
    user_id     BIGINT      REFERENCES users(id)    ON DELETE CASCADE    NOT NULL,
    title       VARCHAR(255),
    created_at  TIMESTAMPTZ                         NOT NULL    DEFAULT NOW()
);

