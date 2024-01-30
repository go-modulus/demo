-- migrate:up

CREATE FUNCTION fill_translation_pseudo_id()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS $$
BEGIN
    new.pseudo_id := new.path || '.' || new.key || '.' || new.locale;
    RETURN new;
END;
$$;

ALTER TABLE translation.translation ADD COLUMN pseudo_id text not null default gen_random_uuid()::text;

ALTER TABLE translation.translation ADD CONSTRAINT idx_translation_pseudo_id_uniq UNIQUE (pseudo_id);

CREATE TRIGGER tg_fill_translation_pseudo_id
    BEFORE INSERT OR UPDATE ON translation.translation
    FOR EACH ROW EXECUTE PROCEDURE fill_translation_pseudo_id();

-- migrate:down
ALTER TABLE translation.translation DROP CONSTRAINT idx_translation_pseudo_id_uniq;
ALTER TABLE translation.translation DROP COLUMN pseudo_id;
DROP TRIGGER tg_fill_translation_pseudo_id ON translation.translation;
DROP FUNCTION fill_translation_pseudo_id();



