CREATE TABLE carts
(
    id serial primary key,
    user_id int references users (id) on delete cascade not null unique
);

CREATE TABLE categories
(
    id serial primary key,
    category_title varchar(255) not null unique,
    img_url text
);

CREATE TABLE products_types
(
    id serial primary key,
    type_title varchar(64) not null unique
);

CREATE TABLE products
(
    id serial primary key,
    article varchar(32) not null unique,
    category_id int references categories (id) on delete cascade not null,
    product_title varchar(255) not null,
    img_url text,
    type_id int references products_types (id) on delete cascade not null,
    amount_in_stock decimal(12, 2),
    price decimal(12, 2) not null,
    units_in_package int not null,
    packages_in_box int not null
);

CREATE TABLE carts_products
(
    cart_id int references carts (id) on delete cascade not null,
    product_id int references products (id) on delete cascade not null,
    quantity int not null,
    price_for_quantity decimal(12, 2) not null,
    primary key (cart_id, product_id)
);

CREATE TABLE products_info
(
    product_id int references products (id) on delete cascade not null,
    info_title varchar(255) not null,
    description varchar(255) not null,
    primary key (product_id, info_title, description)
);
