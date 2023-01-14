CREATE TABLE carts
(
    id      SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL UNIQUE
);

CREATE TABLE categories
(
    id             SERIAL PRIMARY KEY,
    category_title TEXT NOT NULL UNIQUE,
    img_url        TEXT,
    img_uuid       UUID,
    description    TEXT
);

CREATE TABLE products
(
    id              SERIAL PRIMARY KEY,
    article         VARCHAR(6)     NOT NULL UNIQUE,
    category_id     INT            REFERENCES categories (id) ON DELETE SET NULL,
    product_title   TEXT           NOT NULL,
    img_url         TEXT,
    img_uuid        UUID,
    amount_in_stock DECIMAL(12, 2) DEFAULT 0,
    price           DECIMAL(12, 2) NOT NULL,
    characteristics JSONB,
    packaging       JSONB,
    is_active       BOOL           DEFAULT TRUE,
    created_at      TIMESTAMPTZ    DEFAULT (now() AT TIME ZONE 'utc-3'),
    CONSTRAINT valid_price CHECK ( price > 0 ),
    CONSTRAINT valid_amount_in_stock CHECK ( amount_in_stock >= 0 )
);

CREATE TABLE carts_products
(
    cart_id    INT REFERENCES carts (id) ON DELETE CASCADE    NOT NULL,
    product_id INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    quantity   INT                                            NOT NULL,
    PRIMARY KEY (cart_id, product_id),
    CONSTRAINT valid_quantity CHECK ( quantity > 0 )
);

-- CREATE TABLE shopping_sessions
-- (
--     id INT PRIMARY KEY,
--     session UUID PRIMARY KEY,
--     created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
-- );
--
-- CREATE TABLE shopping_sessions_products
-- (
--     shopping_session_id INT REFERENCES shopping_sessions (id) ON DELETE CASCADE NOT NULL,
--     product_id INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
--     PRIMARY KEY (shopping_session_id, product_id)
-- );

-- When we create a new session for user in cache we add to a database too.
-- If admin/moderator want to edit or delete a product we just search in table shopping_sessions_products
-- for product that has been changed and clear the session in redis with joined session uuid from shopping_sessions
-- if session UUID wasn't found in cache we delete a current shopping session from database