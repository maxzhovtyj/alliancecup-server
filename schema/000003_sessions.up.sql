CREATE TABLE sessions
(
    id            serial primary key,
    user_id       int references users (id) not null,
    role_id       int references roles (id) not null,
    refresh_token varchar(255)              not null,
    client_ip     varchar(255)              not null,
    user_agent    varchar(255)              not null,
    is_blocked    boolean     default false,
    expires_at    timestamptz               not null,
    created_at    timestamptz default (now())
);