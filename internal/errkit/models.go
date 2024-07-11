package errkit

import (
	"fmt"
)

type errorResponse struct {
	Code     int                    `json:"code"`
	Error    string                 `json:"error"`
	MetaInfo map[string]interface{} `json:"meta_info,omitempty"`
}

func newErrorResponse(err error) errorResponse {
	if customErr, ok := err.(*Error); ok {
		return errorResponse{
			Code:     customErr.Code,
			Error:    fmt.Sprint(customErr.ErrorMessage),
			MetaInfo: customErr.MetaInfo,
		}
	} else {
		return errorResponse{
			Code:  500,
			Error: err.Error(),
		}
	}
}
