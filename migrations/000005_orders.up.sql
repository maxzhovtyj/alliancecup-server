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
    user_id           INT REFERENCES users (id) ON DELETE CASCADE DEFAULT NULL,
    user_lastname     TEXT                               NOT NULL,
    user_firstname    TEXT                               NOT NULL,
    user_middle_name  TEXT                               NOT NULL,
    user_phone_number TEXT                               NOT NULL,
    user_email        TEXT                               NOT NULL,
    order_status      TEXT                                        DEFAULT 'IN_PROGRESS',
    order_comment     TEXT                                        DEFAULT NULL,
    order_sum_price   DECIMAL(12, 2)                     NOT NULL,
    delivery_type_id  INT REFERENCES delivery_types (id) NOT NULL,
    payment_type_id   INT REFERENCES payment_types (id)  NOT NULL,
    created_at        TIMESTAMPTZ                                 DEFAULT (now() AT TIME ZONE 'utc-3'),
    closed_at         TIMESTAMPTZ                                 DEFAULT (NULL AT TIME ZONE 'utc-3')
);

CREATE TABLE orders_delivery
(
    order_id             INT REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    delivery_title       TEXT                                         NOT NULL,
    delivery_description TEXT                                         NOT NULL,
    PRIMARY KEY (order_id, delivery_title, delivery_description)
);

CREATE TABLE orders_products
(
    order_id           INT REFERENCES orders (id) ON DELETE CASCADE   NOT NULL,
    product_id         INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    quantity           INT                                            NOT NULL,
    price_for_quantity DECIMAL(12, 2)                                 NOT NULL,
    PRIMARY KEY (order_id, product_id)
);