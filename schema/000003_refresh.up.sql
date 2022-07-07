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

ALTER TABLE users ALTER COLUMN phone_number SET NOT NULL;

ALTER TABLE products DROP COLUMN in_stock;

ALTER TABLE orders ALTER COLUMN amount SET NOT NULL;
ALTER TABLE orders ALTER COLUMN end_status SET DEFAULT FALSE;
