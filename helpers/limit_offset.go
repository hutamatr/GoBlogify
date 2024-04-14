package helpers

import (
	"net/http"
	"strconv"
)

func GetLimitOffset(request *http.Request) (int, int) {
	limit := request.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	limitVal, err := strconv.Atoi(limit)
	PanicError(err, "failed to convert limit to int")
	offset := request.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	offsetVal, err := strconv.Atoi(offset)
	PanicError(err, "failed to convert offset to int")
	return limitVal, offsetVal
}
