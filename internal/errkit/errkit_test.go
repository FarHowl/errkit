package errkit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedErrors(t *testing.T) {
	actualErr := NewError("nested message").WithMetaInfo("service", "http://shard-query-storage/api").WithCode(500).WithCode(400).WithCode(404)
	assert.Equal(t, 500, actualErr.Code)
	assert.Equal(t, "nested message", actualErr.ErrorMessage)

	wrapErr := NewError(actualErr).WithMetaInfo("service", "http://shard-query-indexer/api")
	assert.Equal(t, "http://shard-query-storage/api", wrapErr.MetaInfo["service"])
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

func TestMarshalError(t *testing.T) {
	actualErr := NewError(NewError("message").WithCode(400).WithMetaInfo("service", "http://shard-query-storage/api"))
	errData, errCode := Marshal(actualErr)

	actualErr = NewError(errData).WithMetaInfo("service", "http://aggregator/api").WithCode(666)
	fmt.Println(actualErr.Error())
	assert.Equal(t, "[errkit.TestMarshalError]-->(http://aggregator/api): message", actualErr.Error())
	assert.Equal(t, 400, errCode)
	assert.Equal(t, "http://aggregator/api", actualErr.MetaInfo["service"])
}

func TestUnknownMarshal(t *testing.T) {
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

	differentErrData, differentErrCode := Marshal(actualErr)
	actualErr = NewError(differentErrData)
	fmt.Println(actualErr.Error())
	assert.Equal(t, "[errkit.TestUnknownMarshal]: couldn`t parse error form another service", actualErr.Error())
	assert.Equal(t, 520, differentErrCode)

	actualErr = NewError([]byte{})
	assert.Equal(t, "[errkit.TestUnknownMarshal]: expected json message, received - byte[]", actualErr.Error())
}
