package exception

type NotFoundError struct {
	Error string `json:"error"`
}

func NewNotFoundError(err string) NotFoundError {
	return NotFoundError{Error: err}
}
