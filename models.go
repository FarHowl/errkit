package errkit

import (
	"fmt"
)

type errorResponse struct {
	Code  int                    `json:"code"`
	Error string                 `json:"error"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

// This structure is used to marshal JSON and send it as a response
func newErrorResponse(customErr *Error) errorResponse {
	return errorResponse{
		Code:  customErr.Code,
		Error: fmt.Sprint(customErr.ErrorMessage),
		Meta:  customErr.Meta,
	}
}
