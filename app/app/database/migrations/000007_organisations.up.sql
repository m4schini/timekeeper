ALTER TABLE raumzeitalpaka.groups
    RENAME TO organisations;

ALTER TABLE raumzeitalpaka.group_has_user
    RENAME TO organisation_has_user;

ALTER TABLE raumzeitalpaka.organisation_has_user
    RENAME COLUMN "group_id" TO organisation_id;

ALTER TABLE raumzeitalpaka.events
    ADD COLUMN "organisation_id" INT NOT NULL DEFAULT 1;

ALTER TABLE raumzeitalpaka.locations
    ADD COLUMN "organisation_id" INT NOT NULL DEFAULT 1;