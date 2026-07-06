CREATE TABLE moods (
    id          BIGSERIAL       PRIMARY KEY,
    user_id     BIGINT          REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    emoji       VARCHAR(10)     NOT NULL,
    date        DATE            NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, date)
);