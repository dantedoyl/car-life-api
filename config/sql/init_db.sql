CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;


-- GRANT ALL PRIVILEGES ON database car_life_api TO postgres;
-- ALTER USER postgres WITH PASSWORD 'ysnpkoyapassword';

create table if not exists events
(
    id          bigserial primary key,
    name        text         not null,
    club_id     bigint       not null,
    creator_id  bigint       not null,
    description text         null,
    event_date  timestamp    not null,
    created_at  timestamp             default CURRENT_TIMESTAMP,
    latitude    float                 DEFAULT 55.753808,
    longitude   float                 DEFAULT 37.620017,
    avatar      varchar(512) NOT NULL DEFAULT '/img/events/default.jpeg',
    chat_id     bigint,
    spectators_count int default 0,
    participants_count int default 0
);

create table if not exists clubs
(
    id                 bigserial primary key,
    name               text         not null,
    description        text         null,
    tags               text[],
    events_count       int                   default 0,
    participants_count int                   default 0,
    subscribers_count int default 0,
    created_at         timestamp             default CURRENT_TIMESTAMP,
    avatar             varchar(512) NOT NULL DEFAULT '/img/clubs/default.jpeg',
    chat_id            bigint
);

create table if not exists tags
(
    id          bigserial primary key,
    name        text not null,
    usage_count int default 0
);

create table if not exists users
(
    vk_id       bigint PRIMARY key,
    name        text         not null,
    surname     text         not null,
    avatar      varchar(512) NOT NULL,
    tags        text[],
    created_at  timestamp default CURRENT_TIMESTAMP,
    description text
);
create table if not exists cars
(
    id          bigserial primary key,
    owner_id    bigint       not null,
    brand       text         not null,
    model       text         not null,
    date        timestamp    not null,
    description text         not null,
    avatar      varchar(512) NOT NULL DEFAULT '/img/cars/default.jpeg',
    body        text,
    engine      text,
    horse_power text,
    name        text
);

CREATE TYPE user_club_status AS ENUM ('admin', 'participant', 'participant_request', 'subscriber', 'moderator');

create table if not exists users_clubs
(
    user_id bigint,
    club_id bigint,
    status  user_club_status,
    PRIMARY KEY (user_id, club_id)
);

CREATE TYPE user_event_status AS ENUM ('admin', 'participant', 'participant_request', 'spectator');

create table if not exists users_events
(
    user_id  bigint,
    event_id bigint,
    status   user_event_status,
    PRIMARY KEY (user_id, event_id)
);

create table if not exists mini_event_type
(
    id                 bigserial primary key,
    public_name        text,
    public_description text
);

create table if not exists mini_events
(
    id          bigserial primary key,
    type_id     bigint not null,
    user_id     bigint not null,
    description text,
    created_at  timestamp,
    ended_at    timestamp,
    latitude    float DEFAULT 55.753808,
    longitude   float DEFAULT 37.620017
);

create table if not exists events_posts
(
    id          bigserial primary key,
    text        text null,
    user_id     bigint not null,
    event_id    bigint not null,
    created_at  timestamp default CURRENT_TIMESTAMP
);

create table if not exists events_posts_attachments
(
    id          bigserial primary key,
    url         text null,
    post_id     bigint not null
);

CREATE TYPE target_type AS ENUM ('club', 'event', 'post', 'car', 'user');

create table if not exists complaints
(
    id bigserial primary key,
    target_type   target_type not null,
    target_id bigint not null,
    user_id bigint not null,
    text text null
);

insert into mini_event_type (public_name, public_description)
VALUES ('Помощь', 'Нужна помощь'),
       ('Мини-сходка', 'Организуй встречу с друзьями'),
       ('Проишествие', 'Оставь, если что-то произошло на дороге');

INSERT INTO tags (name, usage_count)
VALUES ('jdm', 0),
       ('vintage', 0),
       ('4x4', 0),
       ('racing', 0),
       ('cars&coffee', 0),
       ('exclusive car clubs', 0),
       ('photography', 0),
       ('motosports', 0),
       ('supercar', 0),
       ('local', 0),
       ('drift', 0),
       ('brand specific', 0),
       ('Янгтаймер', 0),
       ('Олдтаймер', 0),
       ('Edm', 0),
       ('Stance', 0),
       ('Racecar', 0),
       ('Offroad', 0),
       ('Traveler', 0),
       ('Street warrior', 0),
       ('Coupe', 0),
       ('Suv', 0),
       ('Usdm', 0);
