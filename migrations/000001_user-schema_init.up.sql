CREATE TABLE roles
(
    id SERIAL PRIMARY KEY,
    role_title TEXT NOT NULL UNIQUE
);

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE ,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    phone_number varchar(20) NOT NULL UNIQUE,
    role_id INT REFERENCES roles(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
);

INSERT INTO roles (role_title) VALUES ('CLIENT');
INSERT INTO roles (role_title) VALUES ('MODERATOR');
INSERT INTO roles (role_title) VALUES ('SUPER_ADMIN');