package exception

type BadRequestError struct {
	Error string `json:"error"`
}

func NewBadRequestError(err string) BadRequestError {
	return BadRequestError{Error: err}
}
