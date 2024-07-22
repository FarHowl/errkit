package errkit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedErrors(t *testing.T) {
	actualErr := NewError("nested message").WithMetaInfo("service", "http://storage/api").WithCode(500).WithCode(400).WithCode(404)
	assert.Equal(t, 404, actualErr.Code)
	assert.Equal(t, "nested message", actualErr.ErrorMessage)

	wrapErr := NewError(actualErr).WithMetaInfo("service", "http://indexer/api")
	assert.Equal(t, "http://indexer/api", wrapErr.Meta["service"])
}

func TestNewError(t *testing.T) {
	actualErr := NewError("message").WithCode(500)
	assert.Equal(t, 500, actualErr.Code)
	assert.Equal(t, "message", actualErr.ErrorMessage)
}

func TestInitiator(t *testing.T) {
	actualErr := NewError("message").WithCode(500)
	assert.Equal(t, "errkit.TestInitiator", actualErr.InitiatorFunction)
	assert.Equal(t, "[errkit.TestInitiator]: message", actualErr.Error())
}

func TestMarshalJSON(t *testing.T) {
	actualErr := NewError(NewError("message").WithCode(400).WithMetaInfo("service", "http://storage/api"))
	errData, errCode := MarshalJSON(actualErr)

	actualErr = NewError(errData)
	fmt.Println(actualErr.Error())
	assert.Equal(t, "[errkit.TestMarshalJSON]: message |---Meta-Info-->| service: http://storage/api", actualErr.Error())
	assert.Equal(t, 400, errCode)
	assert.Equal(t, "http://storage/api", actualErr.Meta["service"])
}

func TestUnknownMarshalJSON(t *testing.T) {
	unknownResponse := struct {
		Status   int                    `json:"status-code"`
		Message  string                 `json:"error-message"`
		MetaInfo map[string]interface{} `json:"meta"`
	}{
		Status:  400,
		Message: "some error from another service",
		MetaInfo: map[string]interface{}{
			"service": "localhost:8080",
		},
	}

	unknownResponseData, _ := json.Marshal(unknownResponse)
	actualErr := NewError(unknownResponseData)

	differentErrData, differentErrCode := MarshalJSON(actualErr)
	actualErr = NewError(differentErrData)
	fmt.Println(actualErr.Error())
	assert.Equal(t, "[errkit.TestUnknownMarshalJSON]: couldn`t parse error form another service |---Meta-Info-->| service: localhost:8080", actualErr.Error())
	assert.Equal(t, 520, differentErrCode)
}

func TestByteMessage(t *testing.T) {
	responseErr := "some error"
	responseErrBytes := make([]byte, 0)
	for _, r := range responseErr {
		responseErrBytes = append(responseErrBytes, byte(r))
	}
	actualErr := NewError(responseErrBytes)
	assert.Equal(t, responseErr, actualErr.ErrorMessage)
}
