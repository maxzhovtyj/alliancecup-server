CREATE TABLE sessions
(
    id uuid primary key,
    user_id int references users (id) on delete cascade not null,
    refresh_token varchar not null,
    is_blocked boolean default false,
    client_ip varchar not null,
    user_agent varchar not null,
    expires_at timestamptz not null,
    created_at timestamptz not null default (now())
);