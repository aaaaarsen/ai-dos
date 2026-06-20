CREATE TABLE summaries (
    id          BIGSERIAL       PRIMARY KEY,
    chat_id     BIGINT          REFERENCES chats(id) ON DELETE CASCADE  NOT NULL,
    content     TEXT            NOT NULL,
    created_at  TIMESTAMPTZ     NOT NULL    DEFAULT NOW()
);