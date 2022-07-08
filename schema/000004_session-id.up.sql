DROP TABLE sessions;

CREATE TABLE sessions
(
    id serial primary key,
    user_id int references users (id) on delete cascade not null,
    role_id int references roles (id) on delete cascade not null,
    refresh_token varchar not null,
    is_blocked boolean default false not null,
    client_ip varchar not null,
    user_agent varchar not null,
    expires_at timestamptz not null,
    created_at timestamptz not null default (now())
);