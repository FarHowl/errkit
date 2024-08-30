package errkit

// This method encapsulates error initiator and marshalls error to JSON
func MarshalJSON(err error) ([]byte, int) {
	customErr := NewError(err)
	errResponse := newErrorResponse(customErr)
	errResponseData, _ := jsonEncoder(errResponse)
	return errResponseData, customErr.Code
}
