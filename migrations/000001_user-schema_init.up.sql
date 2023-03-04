CREATE TABLE roles
(
    id    SERIAL PRIMARY KEY,
    title TEXT    NOT NULL UNIQUE,
    code  CHAR(4) NOT NULL UNIQUE
);

CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email         TEXT                      NOT NULL UNIQUE,
    password_hash TEXT                      NOT NULL,
    lastname      TEXT                      NOT NULL,
    firstname     TEXT                      NOT NULL,
    middle_name   TEXT                      NOT NULL,
    phone_number  varchar(20)               NOT NULL UNIQUE,
    role_id       INT REFERENCES roles (id) NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
);

INSERT INTO roles (title, code)
VALUES ('CLIENT', '0001');
INSERT INTO roles (title, code)
VALUES ('MODERATOR', '1002');
INSERT INTO roles (title, code)
VALUES ('SUPERADMIN', '2003');

INSERT INTO users
(email, password_hash, lastname, firstname, middle_name, phone_number, role_id)
VALUES ('zhovtyjshady@gmail.com',
        '647361646b617364693231323331326d646d61636d787a3030af41f41f071b4374175183d6ffdf93a54bc84daf',
        'Жовтанюк',
        'Максим',
        'В''ячеславович',
        '+380 (68) 306-29-75',
        3);