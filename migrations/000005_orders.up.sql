CREATE TABLE delivery_types
(
    id                  serial primary key,
    delivery_type_title varchar(128) not null
);

CREATE TABLE payment_types
(
    id                 serial primary key,
    payment_type_title varchar(128) not null
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
    id                uuid primary key default gen_random_uuid(),
    user_id           int references users (id) default NULL,
    user_lastname     varchar(128)                       not null,
    user_firstname    varchar(128)                       not null,
    user_middle_name  varchar(128)                       not null,
    user_phone_number varchar(20)                        not null,
    user_email        varchar(64)                        not null,
    order_status      varchar(64)      default 'IN_PROGRESS',
    order_comment     text             default null,
    order_sum_price   decimal(12, 2)                     not null,
    delivery_type_id  int references delivery_types (id) not null,
    payment_type_id   int references payment_types (id)  not null,
    created_at        timestamptz      default (now() at time zone 'GMT+3'),
    closed_at         timestamptz      default (null at time zone 'GMT+3')
);

CREATE TABLE orders_delivery
(
    order_id             uuid references orders (id) not null,
    delivery_title       varchar(128)                not null,
    delivery_description text                        not null,
    primary key (order_id, delivery_title, delivery_description)
);

CREATE TABLE orders_products
(
    order_id           uuid references orders (id)  not null,
    product_id         int references products (id) not null,
    quantity           int                          not null,
    price_for_quantity decimal(12, 2)               not null,
    primary key (order_id, product_id)
);