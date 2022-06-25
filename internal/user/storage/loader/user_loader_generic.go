package loader

import (
	"boilerplate/internal/user/storage"
	"context"
	"github.com/google/uuid"
	"github.com/vikstrous/dataloadgen"
	"time"
)

func NewUserLoaderGeneric(finder *storage.Queries) *dataloadgen.Loader[uuid.UUID, storage.User] {
	loader := dataloadgen.NewLoader(
		func(keys []uuid.UUID) (ret []storage.User, errs []error) {
			users := make([]storage.User, len(keys))
			errors := make([]error, len(keys))

			usersMap, err := finder.GetUsersMap(context.Background(), keys)
			if err != nil {
				for i := range keys {
					errors[i] = err
				}
			}

			for i, key := range keys {
				if user, ok := usersMap[key.String()]; ok {
					users[i] = user
				}

			}
			return users, errors
		},
		dataloadgen.WithBatchCapacity(100),
		dataloadgen.WithWait(2*time.Millisecond),
	)
	return loader
}
