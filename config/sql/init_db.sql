
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;


-- GRANT ALL PRIVILEGES ON database car_life_api TO postgres;
-- ALTER USER postgres WITH PASSWORD 'ysnpkoyapassword';

create table if not exists events
(
    id        bigserial primary key,
    name      text not null,
    club_id   bigint not null,
    creator_id bigint not null,
    description text null,
    event_date timestamp not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    latitude  float                 DEFAULT 55.753808,
    longitude float                 DEFAULT 37.620017,
    avatar    varchar(512) NOT NULL DEFAULT '/img/events/default.jpeg'
);

create table if not exists clubs
(
    id        bigserial primary key,
    name      text not null,
    description text null,
    tags text[],
    events_count int default 0,
    participants_count int default 0,
    created_at timestamp default CURRENT_TIMESTAMP,
    avatar    varchar(512) NOT NULL DEFAULT '/img/clubs/default.jpeg'
    );

create table if not exists tags
(
    id        bigserial primary key,
    name      text not null,
    usage_count int default 0
);

create table if not exists users
(
        vk_id bigint PRIMARY key,
        name text not null,
        surname text not null,
        avatar    varchar(512) NOT NULL,
    tags text[],
    created_at timestamp default CURRENT_TIMESTAMP,
    description text
    );
create table if not exists cars
(
    id bigserial primary key,
    owner_id bigint not null,
    brand text not null,
    model text not null,
    date timestamp not null,
    description text not null,
    avatar    varchar(512) NOT NULL DEFAULT '/img/cars/default.jpeg',
    body text,
    engine text,
    horse_power text,
    name text
);

CREATE TYPE user_club_status AS ENUM ('admin', 'participant', 'participant_request', 'subscriber', 'moderator');

create table if not exists users_clubs
    (
    user_id bigint,
    club_id bigint,
    status user_club_status,
    PRIMARY KEY (user_id, club_id)
    );

CREATE TYPE user_event_status AS ENUM ('admin', 'participant', 'participant_request', 'spectator');

create table if not exists users_events
(
    user_id bigint,
    event_id bigint,
    status user_event_status,
    PRIMARY KEY (user_id, event_id)
    );

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
       ('Usdm', 0)
