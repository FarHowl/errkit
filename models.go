package errkit

import (
	"fmt"
)

type errorResponse struct {
	Code  int                    `json:"code"`
	Error string                 `json:"error"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

func newErrorResponse(err error) errorResponse {
	if customErr, ok := err.(*Error); ok {
		return errorResponse{
			Code:  customErr.Code,
			Error: fmt.Sprint(customErr.ErrorMessage),
			Meta:  customErr.Meta,
		}
	} else {
		return errorResponse{
			Code:  500,
			Error: err.Error(),
		}
	}
}
