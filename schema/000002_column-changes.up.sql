ALTER TABLE users
    ADD COLUMN phone_number varchar(20);

ALTER TABLE products
    ADD COLUMN amount int,
    ADD COLUMN in_stock boolean;

ALTER TABLE orders
    ADD COLUMN amount     int,
    ADD COLUMN end_status boolean;