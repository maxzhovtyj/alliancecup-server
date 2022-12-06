CREATE TABLE categories_filtration
(
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE DEFAULT NULL,
    img_url TEXT,
    img_uuid UUID,
    search_key TEXT NOT NULL,
    search_characteristic TEXT NOT NULL,
    filtration_title TEXT NOT NULL,
    filtration_description TEXT,
    filtration_list_id INT REFERENCES categories_filtration(id) DEFAULT NULL
);

