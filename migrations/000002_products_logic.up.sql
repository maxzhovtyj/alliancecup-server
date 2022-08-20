CREATE TABLE carts
(
    id      serial primary key,
    user_id int references users (id) on delete cascade not null unique
);

INSERT INTO users (email, password_hash, name, phone_number, role_id) VALUES ('zhovtyjshady@gmail.com', '647361646b617364693231323331326d646d61636d787a3030af41f41f071b4374175183d6ffdf93a54bc84daf', 'Maksym Zhovtaniuk', '+380683062975', 3);
INSERT INTO carts (user_id) VALUES (1);

CREATE TABLE categories
(
    id             serial primary key,
    category_title varchar(255) not null unique,
    img_url        text
);

INSERT INTO categories (category_title, img_url) values ('Одноразові стакани', 'https://lh3.googleusercontent.com/IdCynSY2wph7gMjirustYgPT2mHY_wWmcXMgP7jKS0YaX3aYWw2ZryeMGPgF8_qxhg0jx5pumyLYrICq7wU9Jm3ZVhdDPWLHHv843BN62xx1IAwxTjs_0NNQAkC8G8Qzm3JTgIivqdS56h89WUqV91n5Hl3OrEUTd68ZRTpYpnj6MO5OTcQCUCZjpcv8OwI9xlO4UZCawSe8Ru8WtA8JCZVzQ9pXL3D2SA9z-JdozCYmE4jZNYc0bl762jLlCDxkbpi83t9zXt3vDCG6-a0ulNYQKKp3SSLoypWSVOo7qE6B5hWr2Br4zgkQroW7WzyhDEpmRb0geFlLpVrMILG4WNwPW5xlWwZ7V0xMtbVQmDyXvYH8imiuRGZFVndDeXaXVqjvY2TrAPSMsXqCfFMx7Lq_Uegtzwsiwu3e9H8R0urBw1rUMb_EOsfqHU5Q8AHs13cjq_QnF-KwJnFc2J0qHsqDp8FFvmrUmdPFISWV2jKVc1LJJaYmVOhmxqRQEtiTiF2HIAu7sFLVhGtDElkOCl4kIBktpXhOopoeiilipInUqW8JlfNnRKxHi6KPBQR1J2pi5Dw1EA_Ltt_Fd6jZxlhxpXgDQ5zY838R1yVofi1mb-ysCc3qpy00u5987A-53EUscKHUBeIQ9D5P6aoLTVsA0bSfp9ohko4jAgZIeEirlCh3WjKMzc9tyW-RVh30KE3nZ-jtyomz6DN_d6vwyTPD4xHOhw36EDO-oXRV8GgUSjDf5MkKT-bX2VJUN-cYnhzJgSQegMihbFbxdi1aUWD0H6vHTRjx=w325-h231-no?authuser=0');
INSERT INTO categories (category_title, img_url) values ('Купольні стакани', 'https://lh3.googleusercontent.com/IdCynSY2wph7gMjirustYgPT2mHY_wWmcXMgP7jKS0YaX3aYWw2ZryeMGPgF8_qxhg0jx5pumyLYrICq7wU9Jm3ZVhdDPWLHHv843BN62xx1IAwxTjs_0NNQAkC8G8Qzm3JTgIivqdS56h89WUqV91n5Hl3OrEUTd68ZRTpYpnj6MO5OTcQCUCZjpcv8OwI9xlO4UZCawSe8Ru8WtA8JCZVzQ9pXL3D2SA9z-JdozCYmE4jZNYc0bl762jLlCDxkbpi83t9zXt3vDCG6-a0ulNYQKKp3SSLoypWSVOo7qE6B5hWr2Br4zgkQroW7WzyhDEpmRb0geFlLpVrMILG4WNwPW5xlWwZ7V0xMtbVQmDyXvYH8imiuRGZFVndDeXaXVqjvY2TrAPSMsXqCfFMx7Lq_Uegtzwsiwu3e9H8R0urBw1rUMb_EOsfqHU5Q8AHs13cjq_QnF-KwJnFc2J0qHsqDp8FFvmrUmdPFISWV2jKVc1LJJaYmVOhmxqRQEtiTiF2HIAu7sFLVhGtDElkOCl4kIBktpXhOopoeiilipInUqW8JlfNnRKxHi6KPBQR1J2pi5Dw1EA_Ltt_Fd6jZxlhxpXgDQ5zY838R1yVofi1mb-ysCc3qpy00u5987A-53EUscKHUBeIQ9D5P6aoLTVsA0bSfp9ohko4jAgZIeEirlCh3WjKMzc9tyW-RVh30KE3nZ-jtyomz6DN_d6vwyTPD4xHOhw36EDO-oXRV8GgUSjDf5MkKT-bX2VJUN-cYnhzJgSQegMihbFbxdi1aUWD0H6vHTRjx=w325-h231-no?authuser=0');

CREATE TABLE products_types
(
    id         serial primary key,
    type_title varchar(64) not null unique
);

INSERT INTO products_types (type_title) values ('Стакан');

CREATE TABLE products
(
    id               serial primary key,
    article          varchar(32)                                          not null unique,
    category_id      int references categories (id) on delete cascade     not null,
    product_title    varchar(255)                                         not null,
    img_url          text,
    type_id          int references products_types (id) on delete cascade not null,
    amount_in_stock  decimal(12, 2),
    price            decimal(12, 2)                                       not null,
    units_in_package int                                                  not null,
    packages_in_box  int                                                  not null,
    created_at       timestamptz default (now())
);

CREATE TABLE carts_products
(
    cart_id            int references carts (id) on delete cascade    not null,
    product_id         int references products (id) on delete cascade not null,
    quantity           int                                            not null,
    price_for_quantity decimal(12, 2)                                 not null,
    primary key (cart_id, product_id)
);

CREATE TABLE products_info
(
    product_id  int references products (id) on delete cascade not null,
    info_title  varchar(255)                                   not null,
    description varchar(255)                                   not null,
    primary key (product_id, info_title, description)
);
