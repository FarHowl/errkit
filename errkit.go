// This package provides a custom `Error` type, which can holds error`s code, message function that threw it
// and some customizable meta information.
//
// This custom `Error` type acts like a default `error` and can be safely as default one.
package errkit

import (
	"runtime"
	"strings"
)

type Error struct {
	Code              int
	InitiatorFunction string      // Shows the function, where the error was thrown from
	ErrorMessage      interface{} // Supported types: `error`, `string`, `JSON`
	Meta              map[string]interface{}
}

// Logs the error message and initiator. Can be changed for different purposes.
func (e *Error) Error() string {
	return errorLog(e)
}

// Sets HTTP status code
func (e *Error) WithCode(code int) *Error {
	e.Code = code
	return e
}

// Sets any meta information
func (e *Error) WithMetaInfo(key string, value interface{}) *Error {
	e.Meta[key] = value
	return e
}

// It`s safe to call NewError(NewError()). The method tries to capture receivedErr`s states.
// Allows to wrap a error received via HTTP-protocol and tries to parse the JSON into its own fields.
func NewError(receivedErr interface{}) *Error {
	switch receivedErr := receivedErr.(type) {
	case *Error:
		return &Error{
			Code:              receivedErr.Code,
			InitiatorFunction: receivedErr.InitiatorFunction,
			ErrorMessage:      receivedErr.ErrorMessage,
			Meta:              receivedErr.Meta,
		}
	case []byte:
		// Check if it`s a marshalled JSON
		jsonErr, isJsonErr := isJsonErr(receivedErr)
		if isJsonErr {
			jsonCode, jsonErrMessage, jsonMetaInfo := parseJsonErr(jsonErr)
			// It`s JSON
			customErr := buildError(jsonErrMessage).WithCode(jsonCode)
			for key, value := range jsonMetaInfo {
				customErr.WithMetaInfo(key, value)
			}
			return customErr
		}
		// It`s not a JSON
		return buildError(string(receivedErr)).WithCode(500)
	default:
		return buildError(receivedErr)
	}
}

func buildError(errMessage interface{}) *Error {
	return &Error{
		InitiatorFunction: initiator(),
		ErrorMessage:      errMessage,
		Meta:              make(map[string]interface{}),
	}
}

func initiator() string {
	pc, _, _, _ := runtime.Caller(3)
	fn := runtime.FuncForPC(pc)
	fullFuncName := fn.Name()
	compactName := compactInitiatorName(fullFuncName)
	return compactName
}

func compactInitiatorName(initiator string) string {
	split := strings.Split(initiator, "/")
	return split[len(split)-1]
}
