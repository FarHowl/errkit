package errkit

import (
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc/status"
)

type Error struct {
	Code              int
	InitiatorFunction string      // Shows the function, where the error was thrown from
	ErrorMessage      interface{} // Supported types: `error`, `string`, `JSON`
	MetaInfo          map[string]interface{}
}

// Logs the error message and initiator. Can be changed for different purposes.
func (e *Error) Error() string {
	metaInfo := []string{}
	for key, value := range e.MetaInfo {
		metaInfo = append(metaInfo, fmt.Sprintf("%s: %v", key, value))
	}
	if len(metaInfo) > 0 {
		return fmt.Sprintf("[%s]: %v ||| %s", e.InitiatorFunction, e.ErrorMessage, strings.Join(metaInfo, " "))
	} else {
		return fmt.Sprintf("[%s]: %v", e.InitiatorFunction, e.ErrorMessage)
	}
}

// Sets HTTP status code. Code field shows only the deepest error`s status code in the chain.
func (e *Error) WithCode(code int) *Error {
	if e.Code == 0 {
		e.Code = code
	}
	return e
}

func (e *Error) WithMetaInfo(key string, value interface{}) *Error {
	if _, keyExists := e.MetaInfo[key]; !keyExists {
		e.MetaInfo[key] = value
	}
	return e
}

// It`s safe to call NewError(NewError()). Nested errors are more valuable.
func NewError(receivedErr interface{}) *Error {
	switch receivedErr := receivedErr.(type) {
	case *Error:
		return &Error{
			Code:              receivedErr.Code,
			InitiatorFunction: receivedErr.InitiatorFunction,
			ErrorMessage:      receivedErr.ErrorMessage,
			MetaInfo:          receivedErr.MetaInfo,
		}
	case []byte:
		// Check if it`s a marshalled JSON
		jsonErr, isJsonErr := isJsonErr(receivedErr)
		if isJsonErr {
			return processJsonErr(jsonErr)
		}
		// Check if it`s a gRPC
		grpcErr, isGrpcErr := isGrpcErr(receivedErr)
		if isGrpcErr {
			return processGrpcErr(grpcErr)
		}
		return buildError(receivedErr).WithCode(500)
	default:
		return buildError(receivedErr)
	}
}

func processJsonErr(jsonErr map[string]interface{}) *Error {
	jsonCode, jsonErrMessage, jsonMetaInfo, err := parseJsonErr(jsonErr)
	if err != nil {
		// It`s not a JSON
		return buildError(err.Error()).WithCode(jsonCode)
	}
	// It`s JSON
	customErr := buildError(jsonErrMessage).WithCode(jsonCode)
	for key, value := range jsonMetaInfo {
		customErr.WithMetaInfo(key, value)
	}
	return customErr
}

func processGrpcErr(grpcErr *status.Status) *Error {
	grpcCode, grpcErrMessage, grpcMetaInfo, err := parseGrpcError(grpcErr)
	if err != nil {
		// It`s not a gRPC
		return buildError(err.Error()).WithCode(grpcCode)
	}
	// It`s gRPC
	customErr := buildError(grpcErrMessage).WithCode(grpcCode)
	for key, value := range grpcMetaInfo {
		customErr.WithMetaInfo(key, value)
	}
	return customErr
}

func buildError(errMessage interface{}) *Error {
	return &Error{
		InitiatorFunction: initiator(),
		ErrorMessage:      errMessage,
		MetaInfo:          make(map[string]interface{}),
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
