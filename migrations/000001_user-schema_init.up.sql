CREATE TABLE roles
(
    id serial primary key,
    role_title varchar(20) not null unique
);

CREATE TABLE users
(
    id serial primary key,
    email varchar(255) not null unique,
    password_hash varchar(255) not null,
    name varchar(255) not null,
    phone_number varchar(20) not null unique,
    role_id int references roles(id) not null,
    created_at timestamptz default (now() at time zone 'utc-3')
);

INSERT INTO roles (role_title) VALUES ('CLIENT');
INSERT INTO roles (role_title) VALUES ('MODERATOR');
INSERT INTO roles (role_title) VALUES ('SUPER_ADMIN');