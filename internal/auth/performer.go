package auth

import "github.com/gofrs/uuid"

type NullPerformer struct {
	Value Performer
	Valid bool
}

type Performer struct {
	Id uuid.UUID
}
