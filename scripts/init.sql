CREATE EXTENSION IF NOT EXISTS citext;

create table if not exists users
(
    nickname citext primary key,
    fullname text   not null,
    about    text,
    email    citext not null unique
);