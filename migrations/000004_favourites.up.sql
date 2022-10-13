CREATE TABLE favourites
(
    user_id    INT REFERENCES users (id) ON DELETE CASCADE    NOT NULL,
    product_id INT REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (user_id, product_id)
);