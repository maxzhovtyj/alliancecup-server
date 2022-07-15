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
    phone_number varchar(20) not null,
    role_id int references roles(id) not null
);