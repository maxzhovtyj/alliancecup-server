CREATE TABLE sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    role_id       INT REFERENCES roles (id)                   NOT NULL,
    refresh_token TEXT                                        NOT NULL,
    client_ip     TEXT                                        NOT NULL,
    user_agent    TEXT                                        NOT NULL,
    is_blocked    BOOLEAN     DEFAULT FALSE,
    expires_at    TIMESTAMPTZ                                 NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
);