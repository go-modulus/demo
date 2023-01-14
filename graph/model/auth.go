package model

func NewErrUnauthorized() ErrUnauthorized {
	return ErrUnauthorized{Message: "Unauthorized"}
}
