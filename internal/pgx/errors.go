package pgx

import "demo/internal/errors"

func NewPgxError(err error) *errors.Error {
	return errors.
		New("pgx.common", "pgx error").
		WithType("PgxError").
		WithCause(err)
}

func IsPgxError(err error) bool {
	return errors.Type(err) == "PgxError"
}
