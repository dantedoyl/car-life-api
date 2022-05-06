CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;

-- GRANT ALL PRIVILEGES ON database car_life_api TO postgres;
-- ALTER USER postgres WITH PASSWORD 'ysnpkoyapassword';


CREATE TABLE IF NOT EXISTS clubs
(
    id                 BIGSERIAL PRIMARY KEY,
    name               TEXT         NOT NULL,
    description        TEXT         NULL,
    tags               TEXT[],
    created_at         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    avatar             VARCHAR(512) NOT NULL DEFAULT '/img/clubs/default.webp',
    chat_id            BIGINT,
    events_count       INT                   DEFAULT 0,
    participants_count INT                   DEFAULT 0,
    subscribers_count  INT                   DEFAULT 0
);

CREATE TABLE IF NOT EXISTS events
(
    id                 BIGSERIAL PRIMARY KEY,
    name               TEXT         NOT NULL,
    description        TEXT         NULL,
    club_id            BIGINT       NOT NULL,
    creator_id         BIGINT       NOT NULL,
    event_date         TIMESTAMP    NOT NULL,
    created_at         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    latitude           FLOAT                 DEFAULT 55.753808,
    longitude          FLOAT                 DEFAULT 37.620017,
    avatar             VARCHAR(512) NOT NULL DEFAULT '/img/events/default.webp',
    chat_id            BIGINT,
    spectators_count   INT                   DEFAULT 0,
    participants_count INT                   DEFAULT 0,

    FOREIGN KEY (club_id) REFERENCES clubs (id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users (vk_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tags
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    usage_count INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users
(
    vk_id       BIGINT PRIMARY KEY,
    name        TEXT         NOT NULL,
    surname     TEXT         NOT NULL,
    avatar      VARCHAR(512) NOT NULL,
    tags        TEXT[],
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cars
(
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT       NOT NULL,
    name        TEXT,
    description TEXT         NOT NULL,
    brand       TEXT         NOT NULL,
    model       TEXT         NOT NULL,
    date        TIMESTAMP    NOT NULL,
    body        TEXT,
    engine      TEXT,
    horse_power TEXT,
    avatar      VARCHAR(512) NOT NULL DEFAULT '/img/cars/default.webp',

    FOREIGN KEY (owner_id) REFERENCES users (vk_id) ON DELETE CASCADE
);

CREATE TYPE user_club_status AS ENUM ('admin', 'participant', 'participant_request', 'subscriber', 'moderator');
CREATE TABLE IF NOT EXISTS users_clubs
(
    user_id BIGINT,
    club_id BIGINT,
    status  user_club_status,

    PRIMARY KEY (user_id, club_id),
    FOREIGN KEY (user_id) REFERENCES users (vk_id) ON DELETE CASCADE,
    FOREIGN KEY (club_id) REFERENCES clubs (id) ON DELETE CASCADE
);

CREATE TYPE user_event_status AS ENUM ('admin', 'participant', 'participant_request', 'spectator');
CREATE TABLE IF NOT EXISTS users_events
(
    user_id  BIGINT,
    event_id BIGINT,
    status   user_event_status,

    PRIMARY KEY (user_id, event_id),
    FOREIGN KEY (user_id) REFERENCES users (vk_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mini_event_type
(
    id                 BIGSERIAL PRIMARY KEY,
    public_name        TEXT,
    public_description TEXT
);

CREATE TABLE IF NOT EXISTS mini_events
(
    id          BIGSERIAL PRIMARY KEY,
    type_id     BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    description TEXT,
    created_at  TIMESTAMP,
    ended_at    TIMESTAMP,
    latitude    FLOAT DEFAULT 55.753808,
    longitude   FLOAT DEFAULT 37.620017,

    FOREIGN KEY (type_id) REFERENCES mini_event_type (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (vk_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS events_posts
(
    id         BIGSERIAL PRIMARY KEY,
    text       TEXT   NULL,
    user_id    BIGINT NOT NULL,
    event_id   BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (vk_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS events_posts_attachments
(
    id      BIGSERIAL PRIMARY KEY,
    post_id BIGINT       NOT NULL,
    url     VARCHAR(512) NOT NULL DEFAULT '/img/events-posts/default.webp',

    FOREIGN KEY (post_id) REFERENCES events_posts (id) ON DELETE CASCADE
);

CREATE TYPE target_type AS ENUM ('club', 'event', 'post', 'car', 'user');
CREATE TABLE IF NOT EXISTS complaints
(
    id          BIGSERIAL PRIMARY KEY,
    target_type target_type NOT NULL,
    target_id   BIGINT      NOT NULL,
    user_id     BIGINT      NOT NULL,
    text        TEXT        NULL,

    FOREIGN KEY (user_id) REFERENCES users (vk_id) ON DELETE CASCADE
);

INSERT INTO mini_event_type (public_name, public_description)
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
