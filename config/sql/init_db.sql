
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
    tags []text,
    events_count int,
    participants_count int,
    created_at timestamp default CURRENT_TIMESTAMP,
    avatar    varchar(512) NOT NULL DEFAULT '/static/events/default.jpeg'
    );

create table if not exists tags
(
    id        bigserial primary key,
    name      text not null,
    usage_count int default 0
);