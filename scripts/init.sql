CREATE EXTENSION IF NOT EXISTS citext;

create table if not exists users
(
    nickname citext primary key,
    fullname text   not null,
    about    text,
    email    citext not null unique
);

create table if not exists forums
(
    title   text   not null,
    useers     citext not null references users (nickname),
    slug    citext not null primary key,
    posts   int default 0,
    threads int default 0
);