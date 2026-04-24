-- +goose Up

-- Replace per-(user,flower,date) uniqueness with per-(flower,date) uniqueness.
-- A flower can only be watered once per day regardless of who does it.
ALTER TABLE waterings
    DROP CONSTRAINT waterings_user_flower_id_watered_by_user_id_watered_date_key,
    ADD  CONSTRAINT waterings_user_flower_id_watered_date_key UNIQUE (user_flower_id, watered_date);

-- +goose Down

ALTER TABLE waterings
    DROP CONSTRAINT waterings_user_flower_id_watered_date_key,
    ADD  CONSTRAINT waterings_user_flower_id_watered_by_user_id_watered_date_key
         UNIQUE (user_flower_id, watered_by_user_id, watered_date);
