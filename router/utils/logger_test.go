package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Logger act as a regular logger that counts logged errors and warnings.
func TestLogger(t *testing.T) {

	testCases := map[LogLevel]struct {
		levelName string
		expected  string
	}{
		Info: {"info", `[Info] info 1
[Warning] warn 1
[Error] error 1
[Error] error 2
[Info] info 2
[Error] error 3
[Error] error 4
[Error] error 5
[Error] error 6
[Info] info 3
[Warning] warn 2
[Warning] warn 3
[Warning] warn 4
[Warning] warn 5
[Error] new logger error 1
`},
		Warn: {"warn", `[Warning] warn 1
[Error] error 1
[Error] error 2
[Error] error 3
[Error] error 4
[Error] error 5
[Error] error 6
[Warning] warn 2
[Warning] warn 3
[Warning] warn 4
[Warning] warn 5
[Error] new logger error 1
`},
		Error: {"error", `[Error] error 1
[Error] error 2
[Error] error 3
[Error] error 4
[Error] error 5
[Error] error 6
[Error] new logger error 1
`},
		Off: {"off", "[Error] new logger error 1\n"},
	}

	for level, testCase := range testCases {
		t.Run(fmt.Sprintf("logger with level %s", testCase.levelName), func(t *testing.T) {
			l := NewInMemoryLoggerWithLevel(level)
			l.Info("info 1")

			assert.Nil(t, l.MustHaveNoWarnings())
			assert.Nil(t, l.MustHaveNoErrors())
			assert.Nil(t, l.MustHaveNoWarningsf("%s error", "shouldn't"))
			assert.Nil(t, l.MustHaveNoErrorsf("%s error", "shouldn't"))

			l.Warn("warn 1")
			assert.ErrorContains(t, l.MustHaveNoWarnings(), "one warning logged")
			assert.ErrorContains(t, l.MustHaveNoWarningsf("warnings (%d)", l.Warnings()), "warnings (1)")
			assert.Nil(t, l.MustHaveNoErrors())
			assert.Nil(t, l.MustHaveNoErrorsf("%s error", "shouldn't"))

			l.Error("error 1")
			assert.ErrorContains(t, l.MustHaveNoWarnings(), "one error and one warning logged")
			assert.ErrorContains(t, l.MustHaveNoWarningsf("warnings (%d)", l.Warnings()), "warnings (1)")
			assert.ErrorContains(t, l.MustHaveNoErrors(), "one error logged")
			assert.ErrorContains(t, l.MustHaveNoErrorsf("errors (%d)", l.Errors()), "errors (1)")
			l.Error("error 2")
			l.Infof("%s %d", "info", 2)
			l.Error("error 3")
			l.Errorf("error %d", 4)
			l.ErrorIfNotNil(nil)
			l.ErrorIfNotNil("error 5")
			l.ErrorIfNotNilf(nil, "this is not an error")
			l.ErrorIfNotNilf(true, "error %d", 6)
			l.Info("info 3")
			l.Warn("warn 2")
			l.Warnf("warn %d", 3)
			l.WarnIfNotNil(nil)
			l.WarnIfNotNil("warn 4")
			l.WarnIfNotNilf(nil, "this is not a warning")
			l.WarnIfNotNilf(true, "warn %d", 5)

			assert.ErrorContains(t, l.MustHaveNoWarnings(), "6 errors and 5 warnings logged")
			assert.ErrorContains(t, l.MustHaveNoErrors(), "6 errors logged")

			l2 := l.NewCounter()
			l2.WithLevel(Error)
			l2.Info("new logger info 1")
			l2.Warn("new logger warn 1")
			l2.Error("new logger error 1")

			assert.Equal(t, testCase.expected, l.Printed())
			assert.Equal(t, LogCounters{Errors: 6, Warnings: 5}, l.Counters())
			assert.Equal(t, LogCounters{Errors: 1, Warnings: 1}, l2.Counters())
			assert.Equal(t, LogCounters{Errors: 0, Warnings: 0}, l.NewCounter().Counters())
			assert.Equal(t, 5, l.Warnings())
			assert.Equal(t, 6, l.Errors())
			l.AppendCounters(l2.Counters())
			assert.Equal(t, LogCounters{Errors: 7, Warnings: 6}, l.Counters())
		})
	}
}

func TestLogLevelMarshalText(t *testing.T) {
	logLevels := []LogLevel{Info, Warn, Error, Off}
	jsonBytes, err := json.Marshal(logLevels)
	require.NoError(t, err)
	assert.Equal(t, `["info","warn","error","off"]`, string(jsonBytes))

	_, err = LogLevel(10).MarshalText()
	require.Error(t, err)
}

func TestLogLevelUnmarshalText(t *testing.T) {
	var behaviours []LogLevel
	err := json.Unmarshal([]byte(`["info","warn","error","off"]`), &behaviours)
	require.NoError(t, err)

	assert.Equal(t, []LogLevel{Info, Warn, Error, Off}, behaviours)

	var behaviour LogLevel
	err = json.Unmarshal([]byte(`"invalid-log-level"`), &behaviour)
	require.Error(t, err)
}
