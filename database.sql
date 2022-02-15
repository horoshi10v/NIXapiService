create table ingredient
(
    id   int auto_increment
        primary key,
    name varchar(64) not null,
    constraint ingredient_name_uindex
        unique (name)
);

create table product
(
    id    int         not null
        primary key,
    name  varchar(64) not null,
    price float       null,
    image text        null,
    type  varchar(16) null,
    constraint product_id_uindex
        unique (id),
    constraint product_name_uindex
        unique (name)
);

create table product_ingredient
(
    product_id    int null,
    ingredient_id int null,
    constraint ingredient_id
        foreign key (ingredient_id) references ingredient (id),
    constraint product_id
        foreign key (product_id) references product (id)
            on update cascade on delete cascade
);

create table restaurant
(
    id       int auto_increment
        primary key,
    name     varchar(64) not null,
    image    text        not null,
    type     varchar(16) null,
    close_at varchar(8)  null,
    open_at  varchar(8)  null
);

create table menu_products
(
    product_id int not null,
    rest_id    int not null,
    price      int not null,
    constraint menu_products_product_id_uindex
        unique (product_id),
    constraint products_id
        foreign key (product_id) references product (id)
            on update cascade on delete cascade,
    constraint rests_id
        foreign key (rest_id) references restaurant (id)
            on update cascade on delete cascade
);

