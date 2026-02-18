ALTER TABLE raumzeitalpaka.users
    DROP COLUMN IF EXISTS last_login,
    DROP COLUMN IF EXISTS display_name;