package exception

type UnauthorizedError struct {
	Error string `json:"error"`
}

func NewUnauthorizedError(error string) UnauthorizedError {
	return UnauthorizedError{Error: error}
}
