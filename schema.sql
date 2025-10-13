
CREATE SCHEMA timekeeper;

CREATE TABLE timekeeper.users (
   id SERIAL PRIMARY KEY ,
   login_name VARCHAR NOT NULL ,
   password DATE NOT NULL
);

CREATE TYPE EVENT_ROLE AS ENUM ('Organizer', 'Mentor', 'Participant');

CREATE TABLE timekeeper.events (
    id SERIAL PRIMARY KEY ,
    name VARCHAR NOT NULL ,
    start DATE NOT NULL,
    owner INT NOT NULL REFERENCES timekeeper.users(id)
);

CREATE TABLE timekeeper.locations (
    id SERIAL PRIMARY KEY ,
    name VARCHAR NOT NULL ,
    file VARCHAR,
    osm_id VARCHAR
);
--  https://nominatim.openstreetmap.org/lookup?osm_ids=W286396721
--  https://nominatim.openstreetmap.org/lookup?osm_ids=N290381165

CREATE TABLE timekeeper.rooms (
  id SERIAL PRIMARY KEY ,
  location INT NOT NULL REFERENCES timekeeper.locations(id) ON DELETE CASCADE ,
  name VARCHAR NOT NULL,
  location_x INT NOT NULL ,
  location_y INT NOT NULL ,
  location_w INT NOT NULL ,
  location_h INT NOT NULL
);

CREATE TABLE timekeeper.timeslots (
    id SERIAL PRIMARY KEY ,
    event INT NOT NULL REFERENCES timekeeper.events(id) ON DELETE CASCADE ,
    title VARCHAR NOT NULL ,
    note VARCHAR NOT NULL ,
    day INT NOT NULL ,
    start TIME NOT NULL,
    room INT REFERENCES timekeeper.rooms(id),
    role EVENT_ROLE NOT NULL DEFAULT 'Organizer'
)