package helper

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONFromRequest(request *http.Request, result interface{}) {
	err := json.NewDecoder(request.Body).Decode(result)
	PanicError(err)
}

func EncodeJSONFromResponse(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(response)
	PanicError(err)
}
