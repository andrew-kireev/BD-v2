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

create table if not exists users_forum
(
    nickname citext NOT NULL,
    slug     citext NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (slug) REFERENCES forums (slug),
    UNIQUE (nickname, slug)
);

CREATE OR REPLACE FUNCTION add_user_to_forum() RETURNS TRIGGER AS
$add_user_to_forum$
BEGIN
    INSERT INTO users_forum (nickname, slug) VALUES (NEW.author, NEW.forum) on conflict do nothing;
    return NEW;
end
$add_user_to_forum$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE add_user_to_forum();

CREATE TRIGGER post_insert
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE add_user_to_forum();

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
    created   timestamp with time zone default now(),
    path      BIGINT[]                 default array []::INTEGER[]
);

CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM posts WHERE id = new.parent INTO parent_path;
        SELECT thread FROM posts WHERE id = parent_path[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION '99' USING ERRCODE = '00409';
        end if;

        NEW.path := NEW.path || parent_path || new.id;
    end if;
    UPDATE forums SET posts=forums.posts + 1 WHERE forums.slug = new.forum;
    RETURN new;
end
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER path_update_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();

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