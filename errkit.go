package errkit

// This package provides a custom error wrapper, which stores HTTP-code, Initiator (where the error was thrown from),
// error message and customizable meta information
import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Code int
	// Shows the function and package, where the error was thrown from. Example: db_handler.SelectQueries
	InitiatorFunction string
	// Supported types: `error`, `string`, `JSON`
	ErrorMessage interface{}
	// Multiple keys/values can be set. This field will be also shown in MarshalJSON
	Meta map[string]interface{}
}

// Logs the error message and initiator. Can be changed for different purposes.
func (e *Error) Error() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s]: %v", e.InitiatorFunction, e.ErrorMessage))
	if len(e.Meta) > 0 {
		builder.WriteString(" |---Meta-Info-->| ")
		metaInfo := make([]string, 0, len(e.Meta))
		for key, value := range e.Meta {
			metaInfo = append(metaInfo, fmt.Sprintf("%s: %v", key, value))
		}
		builder.WriteString(strings.Join(metaInfo, "|"))
	}
	return builder.String()
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

// Throws the customizable error. Can inherit wrapped error fields
//
// It`s safe to call NewError(NewError(NewError()))
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
			customErr := errorConstructor(jsonErrMessage).WithCode(jsonCode)
			for key, value := range jsonMetaInfo {
				customErr.WithMetaInfo(key, value)
			}
			return customErr
		}
		// It`s a byte[] array - we don`t know, how to parse it, so just return new custom error
		return errorConstructor(receivedErr).WithCode(500)
	default:
		return errorConstructor(receivedErr)
	}
}

// Simple Error errorConstructor
func errorConstructor(errMessage interface{}) *Error {
	return &Error{
		InitiatorFunction: initiator(),
		ErrorMessage:      errMessage,
		Meta:              make(map[string]interface{}),
	}
}

// Retrieves the package and function, where the error was thrown from
func initiator() string {
	pc, _, _, _ := runtime.Caller(3)
	fn := runtime.FuncForPC(pc)
	fullFuncName := fn.Name()
	compactName := compactInitiatorName(fullFuncName)
	return compactName
}

// Deletes excess information
func compactInitiatorName(initiator string) string {
	split := strings.Split(initiator, "/")
	return split[len(split)-1]
}
