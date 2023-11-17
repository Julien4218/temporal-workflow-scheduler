package utils

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

func MapToLowerKeys(f interface{}) interface{} {
	switch f := f.(type) {
	case map[string]interface{}:
		lf := make(map[string]interface{}, len(f))
		for k, v := range f {
			lf[LowerCaseInitial(k)] = MapToLowerKeys(v)
		}
		return lf
	default:
		return f
	}
}

func LowerCaseInitial(str string) string {
	if len(str) == 2 {
		return strings.ToLower(str)
	}

	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func StructToMap(f interface{}) map[string]interface{} {
	var resultMap map[string]interface{}

	jsonMap, _ := json.Marshal(f)
	err := json.Unmarshal(jsonMap, &resultMap)
	if err != nil {
		return nil
	}

	return resultMap
}

func GetLevelFromString(logLevel string, defaultLevel logrus.Level) logrus.Level {
	switch level := strings.ToUpper(logLevel); level {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	default:
		logrus.Debugf("log level %s could not be parsed, using default level %s", level, defaultLevel)
		return defaultLevel
	}
}
