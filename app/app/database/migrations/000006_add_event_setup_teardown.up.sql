ALTER TABLE raumzeitalpaka.events
    ADD COLUMN "setup" INT NOT NULL DEFAULT 0;

ALTER TABLE raumzeitalpaka.events
    ADD COLUMN "teardown" INT NOT NULL DEFAULT 0;