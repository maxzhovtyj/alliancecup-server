ALTER TABLE products
    DROP COLUMN current_supply;

ALTER TABLE products
    DROP COLUMN current_spend;

ALTER TABLE products
    DROP COLUMN current_write_off;

DROP TABLE inventory_products;
DROP TABLE inventory;