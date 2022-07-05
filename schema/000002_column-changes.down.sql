ALTER TABLE users DROP COLUMN phone_number;

ALTER TABLE products
    DROP COLUMN amount,
    DROP COLUMN in_stock;

ALTER TABLE orders
    DROP COLUMN amount,
    DROP COLUMN end_status;