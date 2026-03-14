package apperror

type AppError struct {
	Message string
	Status  int
}

// Implementasi interface error
func (e AppError) Error() string {
	return e.Message
}

func (e AppError) HTTPStatus() int {
	return e.Status
}

func NewNotFound(msg string) error {
	return &AppError{Message: msg, Status: 404}
}

func NewInternal(msg string) error {
	return &AppError{Message: msg, Status: 500}
}

func NewBadRequest(msg string) error {
	return &AppError{Message: msg, Status: 400}
}

func NewUnauthorized(msg string) error {
	return &AppError{Message: msg, Status: 401}
}

func NewForbidden(msg string) error {
	return &AppError{Message: msg, Status: 403}
}
