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

create table if not exists posts
(
    id        bigserial primary key,
    parent    bigint                   default 0,
    author    citext not null references users (nickname),
    message   text   not null,
    is_edited bool   not null          default false,
    forum     citext references forums (slug),
    thread    int references threads (id),
    created   timestamp with time zone default now()
);

create table if not exists thread_votes
(
    nickname  citext not null references users (nickname),
    voice     int,
    thread_id int    not null references threads (id),
    unique (thread_id, nickname)
);


CREATE OR REPLACE FUNCTION add_votes() RETURNS TRIGGER AS
$add_votes$
BEGIN
    UPDATE threads SET votes=(votes + NEW.voice) WHERE id = NEW.thread_id;
    return NEW;
end
$add_votes$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$update_votes$
BEGIN
    UPDATE threads SET votes = (votes + NEW.voice * 2) WHERE id = NEW.thread_id;
    return NEW;
end
$update_votes$ LANGUAGE plpgsql;


CREATE TRIGGER add_vote
    BEFORE INSERT
    ON thread_votes
    FOR EACH ROW
EXECUTE PROCEDURE add_votes();

CREATE TRIGGER edit_vote
    BEFORE UPDATE
    ON thread_votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();