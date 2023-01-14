package custom

import (
	"context"
	"demo/internal/validator"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"io"
	"strconv"
	"time"
)

func MarshalMilliTimestamp(time time.Time) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		_, _ = w.Write([]byte(strconv.FormatInt(time.UnixMilli(), 10)))
		return nil
	})
}

func UnmarshalMilliTimestamp(ctx context.Context, value interface{}) (time.Time, error) {
	rawTimestamp, ok := value.(string)
	if ok {
		timestamp, err := strconv.ParseInt(rawTimestamp, 10, 64)
		if err == nil {
			return time.UnixMilli(timestamp), nil
		}
	}

	pathCtx := graphql.GetPathContext(ctx)
	return time.Time{}, validator.NewValidationError([]validator.FieldValidationError{
		{
			Field:   pathCtx.Path().String(),
			Code:    is.ErrUUID.Code(),
			Message: is.ErrUUID.Message(),
		},
	})
}
