CREATE TABLE raumzeitalpaka.organisation_roles
(
    organisation INT NOT NULL,
    id      VARCHAR NOT NULL ,
    name    VARCHAR NOT NULL,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY(organisation, id)
);

ALTER TABLE raumzeitalpaka.organisation_has_user
    DROP COLUMN IF EXISTS role;

ALTER TABLE raumzeitalpaka.organisation_has_user
    ADD COLUMN role VARCHAR;