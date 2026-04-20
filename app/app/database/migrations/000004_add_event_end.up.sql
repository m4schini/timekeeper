ALTER TABLE raumzeitalpaka.events
    ADD COLUMN "end" DATE NOT NULL;

UPDATE raumzeitalpaka.events
SET "end" = start + (
    SELECT (count(day)-1) * INTERVAL '1 day'
    FROM (SELECT DISTINCT day FROM raumzeitalpaka.timeslots WHERE event = events.id) AS days
);