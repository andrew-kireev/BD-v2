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
    useers  citext not null references users (nickname),
    slug    citext not null primary key,
    posts   int default 0,
    threads int default 0
);

create table if not exists threads
(
    id      serial primary key,
    title   text   not null,
    author  citext not null references users (nickname),
    forum   citext not null references forums (slug),
    message text   not null,
    votes   int                      default 0,
    slug    citext unique,
    created timestamp with time zone default now()
);