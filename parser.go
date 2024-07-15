package errkit

import (
	"encoding/json"
)

// Dictionaries of JSON fields`s names. Are used to parse different keys from the received JSON error
var (
	codeFieldDictionary = []string{
		"status", "code",
	}
	errorFieldDictionary = []string{
		"error", "message",
	}
	metaFieldDictionary = []string{
		"meta", "meta_info",
	}
)

func isJsonErr(errBytes []byte) (map[string]interface{}, bool) {
	var parsedError map[string]interface{}
	err := json.Unmarshal(errBytes, &parsedError)
	if err != nil {
		return nil, false
	}
	return parsedError, true
}

// Tries to retrieve `code` and `error`. Uses customizable dictionaries
func parseJsonErr(jsonErr map[string]interface{}) (int, string, map[string]interface{}) {
	status, exists := getJsonCode(jsonErr, codeFieldDictionary...)
	if !exists {
		status = 520
	}
	errMessage, exists := getJsonErrMessage(jsonErr, errorFieldDictionary...)
	if !exists {
		errMessage = "couldn`t parse error form another service"
	}
	metaInfo, exists := getJsonMetaInfo(jsonErr, metaFieldDictionary...)
	if !exists {
		metaInfo = make(map[string]interface{})
	}
	return status, errMessage, metaInfo
}

// Tries to get HTTP-code from received JSON error
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

// Tries to get message from received JSON error
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

// Tries to get meta information from received JSON error
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
