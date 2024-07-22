package errkit

// These dictionaries allow to parse error received via JSON and inherit its fields
var (
	codeFieldDictionary  = []string{"status", "code"}
	errorFieldDictionary = []string{"error", "message"}
	metaFieldDictionary  = []string{"meta", "meta_info"}
)

// Allows to choose code dictionary that will be used when calling NewError() on received JSON error
func SetCodeFieldDictionary(newDictionary []string) {
	codeFieldDictionary = newDictionary
}

// Allows to choose error dictionary that will be used when calling NewError() on received JSON error
func SetErrorFieldDictionary(newDictionary []string) {
	errorFieldDictionary = newDictionary
}

// Allows to choose meta dictionary that will be used when calling NewError() on received JSON error
func SetMetaFieldDictionary(newDictionary []string) {
	metaFieldDictionary = newDictionary
}

func isJsonErr(errBytes []byte) (map[string]interface{}, bool) {
	var parsedError map[string]interface{}
	err := jsonDecoder(errBytes, &parsedError)
	if err != nil {
		return nil, false
	}
	return parsedError, true
}

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
