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

CREATE TABLE products
(
    id              SERIAL PRIMARY KEY,
    article         TEXT                                                 NOT NULL UNIQUE,
    category_id     INT REFERENCES categories (id) ON DELETE CASCADE     NOT NULL,
    product_title   TEXT                                                 NOT NULL,
    img_url         TEXT,
    amount_in_stock DECIMAL(12, 2) DEFAULT 0,
    price           DECIMAL(12, 2)                                       NOT NULL,
    characteristics JSONB,
    packaging       JSONB,
    created_at      TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3'),
    CONSTRAINT valid_price CHECK ( price > 0 ),
    CONSTRAINT valid_amount_in_stock CHECK ( amount_in_stock >= 0 )
);

CREATE TABLE carts_products
(
    cart_id            INT REFERENCES carts (id) ON DELETE CASCADE    NOT NULL,
    product_id         INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    quantity           INT                                            NOT NULL,
    price_for_quantity DECIMAL(12, 2)                                 NOT NULL,
    PRIMARY KEY (cart_id, product_id),
    CONSTRAINT valid_quantity CHECK ( quantity > 0 )
);
