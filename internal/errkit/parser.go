package errkit

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func isGrpcErr(errBytes []byte) (*status.Status, bool) {
	st := &status.Status{}
	err := proto.Unmarshal(errBytes, st.Proto())
	if err != nil {
		return nil, false
	}
	return st, true
}

func parseGrpcError(grpcErr *status.Status) (int, string, map[string]interface{}, error) {
	metaInfo := make(map[string]interface{})
	details := grpcErr.Details()
	if len(details) > 0 {
		for i, detail := range details {
			metaInfo[fmt.Sprintf("detail_%d", i)] = detail
		}
	}

	return int(grpcErr.Code()), grpcErr.Message(), metaInfo, nil
}

func isJsonErr(errBytes []byte) (map[string]interface{}, bool) {
	var parsedError map[string]interface{}
	err := json.Unmarshal(errBytes, &parsedError)
	if err != nil {
		return nil, false
	}
	return parsedError, true
}

func parseJsonErr(jsonErr map[string]interface{}) (int, string, map[string]interface{}, error) {
	status, exists := getJsonCode(jsonErr, "status", "code")
	if !exists {
		status = 520
	}
	errMessage, exists := getJsonErrMessage(jsonErr, "error", "message")
	if !exists {
		errMessage = "couldn`t parse error form another service"
	}
	metaInfo, exists := getJsonMetaInfo(jsonErr, "meta", "meta_info")
	if !exists {
		metaInfo = make(map[string]interface{})
	}

	return status, errMessage, metaInfo, nil
}

func getJsonCode(parsedMessage map[string]interface{}, keys ...string) (int, bool) {
	for _, key := range keys {
		if value, exists := parsedMessage[key]; exists {
			if status, ok := value.(float64); ok {
				return int(status), true
			} else if status, ok := value.(float32); ok {
				return int(status), true
			} else if status, ok := value.(int64); ok {
				return int(status), true
			} else if status, ok := value.(int32); ok {
				return int(status), true
			}
		}
	}
	return 0, false
}

func getJsonErrMessage(parsedMessage map[string]interface{}, keys ...string) (string, bool) {
	for _, key := range keys {
		if value, exists := parsedMessage[key]; exists {
			if errorMessage, ok := value.(string); ok {
				return errorMessage, true
			}
		}
	}
	return "", false
}

func getJsonMetaInfo(parsedMessage map[string]interface{}, keys ...string) (map[string]interface{}, bool) {
	for _, key := range keys {
		if value, exists := parsedMessage[key]; exists {
			if metaInfo, ok := value.(map[string]interface{}); ok {
				return metaInfo, true
			}
		}
	}
	return nil, false
}
