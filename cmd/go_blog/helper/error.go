package helper

import "errors"

var (
	ValidationError = errors.New("validation error")
	NotFoundError   = errors.New("not found error")
)

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
