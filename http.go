package errkit

import (
	"encoding/json"
)

// This method encapsulates error initiator and marshalls error data to JSON.
//
// Accepts any error and translates it to a errkit.Error
func MarshalJSON(err error) ([]byte, int) {
	customErr := NewError(err)
	errResponse := newErrorResponse(customErr)
	errResponseData, _ := json.Marshal(errResponse)
	return errResponseData, customErr.Code
}
