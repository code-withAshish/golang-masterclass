create table if not exists users(
    id uuid primary key default gen_random_uuid(),
    name varchar(100) not null,
    email varchar(100) not null unique
    );


create table if not exists notes(
    id uuid primary key default gen_random_uuid(),
    title varchar(100) not null,
    description text not null,
    content text not null,
    user_id uuid not null,
    constraint fk_user foreign key(user_id) references users(id)
    )