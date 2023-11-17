package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLowerCaseFirstChar(t *testing.T) {
	testStr := "CliVersion"
	expectedStr := "cliVersion"
	actualStr := LowerCaseInitial(testStr)

	require.Equal(t, expectedStr, actualStr)
}

func TestConvertMapKeysToLower(t *testing.T) {
	testMap := map[string]interface{}{
		"OS":       "windows",
		"TestKey1": "test value 1",
		"TestKey2": "test value 2",
	}

	actualMap := MapToLowerKeys(testMap)
	expectedMap := map[string]interface{}{
		"os":       "windows",
		"testKey1": "test value 1",
		"testKey2": "test value 2",
	}

	require.Equal(t, expectedMap, actualMap)
}

func TestStructToMap(t *testing.T) {
	type myStruct struct {
		Test string
	}

	actualResult := StructToMap(&myStruct{Test: "test"})
	expectedResult := map[string]interface{}{"Test": "test"}

	nilResult := StructToMap(nil)

	require.Equal(t, expectedResult, actualResult)
	require.Nil(t, nil, nilResult)
}

func TestShouldParseTraceLogLevel(t *testing.T) {
	level := GetLevelFromString("TRACE", logrus.DebugLevel)
	require.Equal(t, logrus.TraceLevel, level)
}

func TestShouldParseDebugLogLevel(t *testing.T) {
	level := GetLevelFromString("DEBUG", logrus.InfoLevel)
	require.Equal(t, logrus.DebugLevel, level)
}

func TestShouldParseInfoLogLevel(t *testing.T) {
	level := GetLevelFromString("INFO", logrus.DebugLevel)
	require.Equal(t, logrus.InfoLevel, level)
}

func TestShouldParseInfoLogLevelCaseInsensitive(t *testing.T) {
	level := GetLevelFromString("Info", logrus.DebugLevel)
	require.Equal(t, logrus.InfoLevel, level)
}

func TestShouldParseWarnLogLevel(t *testing.T) {
	level := GetLevelFromString("WARN", logrus.DebugLevel)
	require.Equal(t, logrus.WarnLevel, level)
}

func TestShouldParseErrorLogLevel(t *testing.T) {
	level := GetLevelFromString("ERROR", logrus.DebugLevel)
	require.Equal(t, logrus.ErrorLevel, level)
}

func TestShouldNotParseLogLevel(t *testing.T) {
	level := GetLevelFromString("abcd", logrus.DebugLevel)
	require.Equal(t, logrus.DebugLevel, level)
}
