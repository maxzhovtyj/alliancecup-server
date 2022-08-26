ALTER TABLE categories ADD COLUMN category_description text;

CREATE TABLE categories_filtration
(
    id serial primary key,
    category_id int references categories(id) on delete cascade default NULL,
    img_url text,
    info_description varchar(255) not null,
    filtration_title varchar(255) not null,
    filtration_description text,
    filtration_list_id int references categories_filtration(id) default NULL
);

