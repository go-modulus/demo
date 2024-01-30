-- migrate:up
CREATE TABLE translation.translation (
    -- it can be UUID for any table or any other unique key of a row or entity
    key text NOT NULL,
    -- it can be a field name in a table or a path to json or whatever to identify a translatable field of an entity
    -- for example admin.config.tnc to save a translation of terms and conditions and
    -- admin.config.description to save a translation of the collection's description.
    -- in both cases the key will be the same, but the path will be different
    path translation.path NOT NULL,
    locale translation."locale" NOT NULL,
    value text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (key, path, locale)
);


-- migrate:down
DROP TABLE translation.translation;
