package custom

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"io"
)

type ErrInvalidUuid struct {
	Field string
}

func (e ErrInvalidUuid) Error() string {
	return fmt.Sprintf("invalid uuid for field %s", e.Field)
}

type UuidOrError struct {
	Uuid  uuid.UUID
	Error error
}

func (u *UuidOrError) UnmarshalGQLContext(ctx context.Context, value interface{}) error {
	rawUuid, ok := value.(string)
	if ok {
		id, err := uuid.FromString(rawUuid)
		if err != nil {
			pc := graphql.GetPathContext(ctx)
			u.Error = ErrInvalidUuid{Field: pc.Path().String()}

			return nil
		}

		u.Uuid = id
	}

	return nil
}

func MarshalUuid(id uuid.UUID) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		_, _ = w.Write([]byte(fmt.Sprintf("%q", id.String())))
		return nil
	})
}

func UnmarshalUuid(ctx context.Context, value interface{}) (uuid.UUID, error) {
	rawUuid, ok := value.(string)
	if ok {
		id, err := uuid.FromString(rawUuid)
		if err == nil {
			return id, nil
		}
	}

	return uuid.Nil, nil
}
