
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;


-- GRANT ALL PRIVILEGES ON database car_life_api TO postgres;
-- ALTER USER postgres WITH PASSWORD 'ysnpkoyapassword';

create table if not exists events
(
    id        bigserial primary key,
    name      text not null,
    club_id   bigint not null,
    description text null,
    event_date timestamp not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    latitude  float                 DEFAULT 55.753808,
    longitude float                 DEFAULT 37.620017,
    avatar    varchar(512) NOT NULL DEFAULT '/static/events/default.jpeg'
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
    avatar    varchar(512) NOT NULL DEFAULT '/static/clubs/default.jpeg'
    );

create table if not exists tags
(
    id        bigserial primary key,
    name      text not null,
    usage_count int default 0
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
