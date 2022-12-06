CREATE TABLE delivery_types
(
    id                  SERIAL PRIMARY KEY,
    delivery_type_title TEXT NOT NULL
);

CREATE TABLE payment_types
(
    id                 SERIAL PRIMARY KEY,
    payment_type_title TEXT NOT NULL
);

INSERT INTO delivery_types (delivery_type_title)
values ('Нова Пошта');
INSERT INTO delivery_types (delivery_type_title)
values ('Доставка AllianceCup по м. Рівне');
INSERT INTO delivery_types (delivery_type_title)
values ('Самовивіз');

INSERT INTO payment_types (payment_type_title)
values ('Оплата при отриманні');
INSERT INTO payment_types (payment_type_title)
values ('Переказ на карту');

CREATE TABLE orders
(
    id                SERIAL PRIMARY KEY,
    executed_by       INT                                REFERENCES users (id) ON DELETE SET NULL,
    user_id           INT                                REFERENCES users (id) ON DELETE SET NULL DEFAULT NULL,
    user_lastname     TEXT                               NOT NULL,
    user_firstname    TEXT                               NOT NULL,
    user_middle_name  TEXT                               NOT NULL,
    user_phone_number TEXT                               NOT NULL,
    user_email        TEXT                               NOT NULL,
    status            TEXT                                                                        DEFAULT 'IN_PROGRESS',
    comment           TEXT,
    delivery_type_id  INT REFERENCES delivery_types (id) NOT NULL,
    payment_type_id   INT REFERENCES payment_types (id)  NOT NULL,
    delivery_info     JSONB,
    created_at        TIMESTAMPTZ                                                                 DEFAULT (now() AT TIME ZONE 'utc-3'),
    closed_at         TIMESTAMPTZ                                                                 DEFAULT (NULL AT TIME ZONE 'utc-3')
);

CREATE TABLE orders_products
(
    order_id   INT REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    product_id INT                                          REFERENCES products (id) ON DELETE SET NULL NOT NULL,
    price      DECIMAL(12, 2)                               NOT NULL,
    quantity   INT                                          NOT NULL,
    PRIMARY KEY (order_id, product_id),
    CONSTRAINT valid_quantity CHECK ( quantity > 0 )
);