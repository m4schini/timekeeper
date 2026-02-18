ALTER TABLE raumzeitalpaka.users
    ADD COLUMN last_login TIMESTAMP,
    ADD COLUMN display_name VARCHAR(255);

UPDATE raumzeitalpaka.users
SET display_name = login_name;