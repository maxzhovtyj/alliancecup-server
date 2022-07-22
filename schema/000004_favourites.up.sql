CREATE TABLE favourites
(
    user_id    int references users (id) on delete cascade    not null,
    product_id int references products (id) on delete cascade not null,
    primary key (user_id, product_id)
);