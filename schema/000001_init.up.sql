CREATE TABLE roles
(
    id   serial       primary key,
    role varchar(255) not null
);

CREATE TABLE users
(
    id            serial                                      primary key,
    role_id       int references roles (id) on delete cascade not null,
    email         varchar(255)                                not null unique,
    password_hash varchar(255)                                not null,
    name          varchar(255)                                not null
);

CREATE TABLE categories
(
    id    serial       primary key,
    title varchar(255) not null unique
);

CREATE TABLE products
(
    id             serial                                           primary key,
    category_id    int references categories (id) on delete cascade not null,
    title          varchar(255)                                     not null,
    price          decimal(12, 2)                                   not null,
    size           int                                              not null,
    characteristic varchar(255)                                     not null,
    description    varchar(255)                                     not null
);

CREATE TABLE orders
(
    id         serial                                         primary key,
    product_id int references products (id) on delete cascade not null,
    user_id    int references users (id) on delete cascade    not null,
    order_date date                                           not null
);