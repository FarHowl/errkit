package errkit

import "encoding/json"

// By default `errkit` library uses `encoding/json`
var (
	jsonEncoder func(v any) ([]byte, error)    = json.Marshal
	jsonDecoder func(data []byte, v any) error = json.Unmarshal
)

// Allows to choose your own JSONEncoder
func SetJSONEncoder(newEncoder func(v any) ([]byte, error)) {
	jsonEncoder = newEncoder
}

// Allows to choose your own JSONDecoder
func SetJSONDecoder(newDecoder func(data []byte, v any) error) {
	jsonDecoder = newDecoder
}
