package errkit

import (
	"fmt"
	"strings"
)

var errorLog func(e *Error) string = func(e *Error) string {
	metaInfo := []string{}
	for key, value := range e.Meta {
		metaInfo = append(metaInfo, fmt.Sprintf("%s: %v", key, value))
	}
	metaInfoStr := ""
	if len(metaInfo) > 0 {
		metaInfoStr = strings.Join(metaInfo, "|")
	}
	if len(metaInfoStr) == 0 {
		return fmt.Sprintf("[%s]: %v",
			e.InitiatorFunction,
			e.ErrorMessage,
		)
	} else {
		return fmt.Sprintf(
			"[%s]: %v |---Meta-Info-->| %s",
			e.InitiatorFunction,
			e.ErrorMessage,
			metaInfoStr,
		)
	}
}

// Allows to adjust your own log format
func SetLog(newErrorLog func(e *Error) string) {
	errorLog = newErrorLog
}
