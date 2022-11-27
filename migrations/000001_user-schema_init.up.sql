CREATE TABLE roles
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    code CHAR(4) NOT NULL UNIQUE
);

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE ,
    password_hash TEXT NOT NULL,
    lastname TEXT NOT NULL,
    firstname TEXT NOT NULL,
    middle_name TEXT NOT NULL,
    phone_number varchar(20) NOT NULL UNIQUE,
    role_id INT REFERENCES roles(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
);

INSERT INTO roles (title, code) VALUES ('CLIENT', '0001');
INSERT INTO roles (title, code) VALUES ('MODERATOR', '1002');
INSERT INTO roles (title, code) VALUES ('SUPERADMIN', '2003');