CREATE TABLE inventory
(
    id         SERIAL PRIMARY KEY,
    created_at timestamptz default (now() AT TIME ZONE 'utc-3')
);

CREATE TABLE inventory_products
(
    inventory_id     INT REFERENCES inventory (id) ON DELETE CASCADE NOT NULL,
    product_id       INT REFERENCES products (id) ON DELETE CASCADE  NOT NULL,
    last_inventory   timestamptz DEFAULT NULL,
    initial_amount   DECIMAL(12, 2),
    supply           DECIMAL(12, 2),
    spend            DECIMAL(12, 2),
    write_off        DECIMAL(12, 2),
    write_off_price  DECIMAL(12, 2),
    planned_amount   DECIMAL(12, 2),
    difference       DECIMAL(12, 2),
    difference_price DECIMAL(12, 2),
    PRIMARY KEY (inventory_id, product_id)
);

ALTER TABLE products
    ADD COLUMN current_supply DECIMAL(12, 2) DEFAULT 0;
ALTER TABLE products
    ADD COLUMN current_spend DECIMAL(12, 2) DEFAULT 0;
ALTER TABLE products
    ADD COLUMN current_write_off DECIMAL(12, 2) DEFAULT 0;
ALTER TABLE products
    ADD COLUMN last_inventory timestamptz default null