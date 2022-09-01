CREATE TABLE supply
(
    id              serial primary key,
    supplier        varchar(256)   not null,
    supply_time     timestamptz  default (now()),
    comment         text,
    created_at      timestamp  default (now() at time zone 'utc-3')
);

CREATE TABLE supply_payment
(
    supply_id int references supply(id) on delete cascade not null,
    payment_account varchar(256) default '-',
    payment_time    timestamptz  default (now()),
    payment_sum     decimal(12, 2) not null
);

CREATE TABLE supply_products
(
    supply_id       int references supply (id) on delete cascade not null,
    product_id      int references products (id)                 not null,
    packaging       varchar(12)                                  not null,
    amount          decimal(12, 2)                               not null,
    price_for_unit  decimal(12, 2)                               not null,
    sum_without_tax decimal(12, 2)                               not null,
    tax             decimal(12, 2),
    total_sum       decimal(12, 2)                               not null
);

CREATE TABLE products_review
(
    id          serial primary key,
    user_id     int references users (id) on delete cascade default null,
    user_name   varchar(64),
    mark        int                                         default 5,
    review_text text not null
);