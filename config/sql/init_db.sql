create database car-life-api
	with owner postgres
	encoding 'utf8'
	LC_COLLATE = 'ru_RU.UTF-8'
    LC_CTYPE = 'ru_RU.UTF-8'
    TABLESPACE = pg_default
    TEMPLATE template0
	;

CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;


GRANT ALL PRIVILEGES ON database car-life-api TO postgres;
ALTER USER postgres WITH PASSWORD 'ysnpkoyapassword';

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
    avatar    varchar(512) NOT NULL DEFAULT ''
);