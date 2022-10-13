CREATE TABLE supply
(
    id          SERIAL PRIMARY KEY,
    supplier    TEXT           NOT NULL,
    supply_time TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3'),
    comment     TEXT,
    sum         DECIMAL(12, 2) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3')
);

CREATE TABLE supply_payment
(
    supply_id       INT REFERENCES supply (id) ON DELETE CASCADE NOT NULL,
    payment_account TEXT DEFAULT '-',
    payment_time    TIMESTAMP    DEFAULT (now() at time zone 'utc-3'),
    payment_sum     DECIMAL(12, 2)                               not null
);

CREATE TABLE supply_products
(
    supply_id       INT REFERENCES supply (id) ON DELETE CASCADE NOT NULL,
    product_id      INT REFERENCES products (id)                 NOT NULL,
    packaging       TEXT                                  NOT NULL,
    amount          DECIMAL(12, 2)                               NOT NULL,
    price_for_unit  DECIMAL(12, 2)                               NOT NULL,
    sum_without_tax DECIMAL(12, 2)                               NOT NULL,
    tax             DECIMAL(12, 2),
    total_sum       DECIMAL(12, 2)                               NOT NULL
);

CREATE TABLE products_review
(
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES users (id) ON DELETE CASCADE DEFAULT NULL,
    user_name   TEXT,
    mark        INT                                         DEFAULT 5,
    review_text TEXT NOT NULL,
    created_at  TIMESTAMPTZ                                 DEFAULT (now() AT TIME ZONE 'utc-3')
);