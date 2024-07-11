package errkit

import (
	"encoding/json"
)

// This method encapsulates error initiator and marshalls error to JSON
func Marshal(err error) ([]byte, int) {
	customErr := NewError(err)
	errResponse := newErrorResponse(customErr)
	errResponseData, _ := json.Marshal(errResponse)
	return errResponseData, customErr.Code
}
