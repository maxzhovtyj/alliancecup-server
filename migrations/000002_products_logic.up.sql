CREATE TABLE carts
(
    id      SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL UNIQUE
);

CREATE TABLE categories
(
    id             SERIAL PRIMARY KEY,
    category_title TEXT NOT NULL UNIQUE,
    img_url        TEXT
);

CREATE TABLE products_types
(
    id         SERIAL PRIMARY KEY,
    type_title TEXT NOT NULL UNIQUE
);

INSERT INTO products_types (type_title)
values ('Стакан');

CREATE TABLE products
(
    id               SERIAL PRIMARY KEY,
    article          TEXT                                                 NOT NULL UNIQUE,
    category_id      INT REFERENCES categories (id) ON DELETE CASCADE     NOT NULL,
    product_title    TEXT                                                 NOT NULL,
    img_url          TEXT,
    type_id          INT REFERENCES products_types (id) ON DELETE CASCADE NOT NULL,
    amount_in_stock  DECIMAL(12, 2),
    price            DECIMAL(12, 2)                                       NOT NULL,
    units_in_package INT                                                  NOT NULL,
    packages_in_box  INT                                                  NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
);

CREATE TABLE carts_products
(
    cart_id            INT REFERENCES carts (id) ON DELETE CASCADE    NOT NULL,
    product_id         INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    quantity           INT                                            NOT NULL,
    price_for_quantity DECIMAL(12, 2)                                 NOT NULL,
    PRIMARY KEY (cart_id, product_id)
);

CREATE TABLE products_info
(
    product_id  INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    info_title  TEXT                                           NOT NULL,
    description TEXT                                           NOT NULL,
    PRIMARY KEY (product_id, info_title)
);
