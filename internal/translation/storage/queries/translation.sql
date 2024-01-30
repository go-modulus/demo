-- name: SaveTranslation :one
insert into "translation"."translation"
    (key, path, value, locale)
values (@key, @path, @value, @locale)
ON CONFLICT (key, path, locale) DO UPDATE
    SET value      = excluded.value,
        updated_at = now()
RETURNING *;

-- name: FindTranslationValue :one
select value
from "translation"."translation"
where key = @key
  and path = @path
  and locale = @locale;

-- name: FindTranslations :many
SELECT t.*
FROM jsonb_to_recordset(@keys::jsonb) AS input_ids(key text, path translation.path, locale translation.locale)
JOIN "translation"."translation" as t
     on input_ids.key = t.key and input_ids.path = t.path and input_ids.locale = t.locale;

-- name: DeleteTransactionsByKey :exec
delete from "translation"."translation"
where key = @key;