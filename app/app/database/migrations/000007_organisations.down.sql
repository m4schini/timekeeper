ALTER TABLE raumzeitalpaka.organisations
    RENAME TO groups;

ALTER TABLE raumzeitalpaka.organisation_has_user
    RENAME TO group_has_user;

ALTER TABLE raumzeitalpaka.events
    DROP COLUMN IF EXISTS "organisation_id";

ALTER TABLE raumzeitalpaka.locations
    DROP COLUMN IF EXISTS "organisation_id";