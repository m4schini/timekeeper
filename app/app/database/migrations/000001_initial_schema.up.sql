CREATE SCHEMA raumzeitalpaka;

CREATE TABLE raumzeitalpaka.users
(
    id         SERIAL PRIMARY KEY,
    login_name VARCHAR NOT NULL UNIQUE,
    password   VARCHAR NOT NULL
);

CREATE TYPE EVENT_ROLE AS ENUM ('Organizer', 'Mentor', 'Participant');

CREATE TABLE raumzeitalpaka.events
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR NOT NULL,
    start DATE    NOT NULL,
--     owner INT NOT NULL REFERENCES raumzeitalpaka.users(id),
    slug  VARCHAR NOT NULL UNIQUE,
    guid  uuid    NOT NULL UNIQUE DEFAULT gen_random_uuid()
);

CREATE TABLE raumzeitalpaka.locations
(
    id     SERIAL PRIMARY KEY,
    name   VARCHAR NOT NULL,
    file   VARCHAR,
    osm_id VARCHAR,
    guid   uuid    NOT NULL UNIQUE DEFAULT gen_random_uuid()
);
--  https://nominatim.openstreetmap.org/lookup?osm_ids=W286396721
--  https://nominatim.openstreetmap.org/lookup?osm_ids=N290381165


CREATE TABLE raumzeitalpaka.event_has_location
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    event    INT     NOT NULL REFERENCES raumzeitalpaka.events (id) ON DELETE CASCADE,
    location INT     NOT NULL REFERENCES raumzeitalpaka.locations (id) ON DELETE CASCADE,
    note     VARCHAR NOT NULL,
    visible  BOOL    NOT NULL DEFAULT true
);

CREATE TABLE raumzeitalpaka.rooms
(
    id          SERIAL PRIMARY KEY,
    location    INT     NOT NULL REFERENCES raumzeitalpaka.locations (id) ON DELETE CASCADE,
    name        VARCHAR NOT NULL,
    location_x  INT     NOT NULL,
    location_y  INT     NOT NULL,
    location_w  INT     NOT NULL,
    location_h  INT     NOT NULL,
    description VARCHAR NOT NULL,
    guid        uuid    NOT NULL UNIQUE DEFAULT gen_random_uuid()
);

CREATE TABLE raumzeitalpaka.timeslots
(
    id        SERIAL PRIMARY KEY,
    parent_id INT REFERENCES raumzeitalpaka.timeslots (id) ON DELETE CASCADE,
    event     INT                NOT NULL REFERENCES raumzeitalpaka.events (id) ON DELETE CASCADE,
    title     VARCHAR            NOT NULL,
    note      VARCHAR            NOT NULL,
    day       INT                NOT NULL,
    start     TIME               NOT NULL,
    room      INT REFERENCES raumzeitalpaka.rooms (id),
    role      EVENT_ROLE         NOT NULL        DEFAULT 'Organizer',
    duration  INTERVAL SECOND(0) NOT NULL,
    guid      uuid               NOT NULL UNIQUE DEFAULT gen_random_uuid()
);