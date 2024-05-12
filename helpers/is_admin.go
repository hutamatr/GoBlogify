package helpers

import (
	"net/http"
	"strconv"
)

func IsAdmin(request *http.Request) bool {
	isAdminString := request.Header.Get("isAdmin")
	isAdmin, err := strconv.ParseBool(isAdminString)
	PanicError(err, "failed to convert isAdmin to bool")
	return isAdmin
}
