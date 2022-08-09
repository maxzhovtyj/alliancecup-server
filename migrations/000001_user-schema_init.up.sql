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
    role_id int references roles(id) not null
);

INSERT INTO roles (role_title) VALUES ('CLIENT');
INSERT INTO roles (role_title) VALUES ('MODERATOR');
INSERT INTO roles (role_title) VALUES ('SUPER_ADMIN');

INSERT INTO users (email, password_hash, name, phone_number, role_id) VALUES ('zhovtyjshady@gmail.com', '647361646b617364693231323331326d646d61636d787a3030af41f41f071b4374175183d6ffdf93a54bc84daf', 'Maksym Zhovtaniuk', '+380683062975', 3);