package errors

type UserError struct {
	Code       string
	Message    string
	DontHandle bool
	Extra      map[string]interface{}
}

type UserErrorProvider interface {
	ToUserError() *UserError
}
