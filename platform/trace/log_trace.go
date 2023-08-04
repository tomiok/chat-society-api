package trace

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func getStackTrace() (uintptr, string, int, error) {
	pc, file, line, ok := runtime.Caller(2)

	if !ok {
		return 0, "", 0, errors.New("cannot get stack trace")
	}

	return pc, file, line, nil
}

// Trace is a helper for logs format. Prints information about the file, line and function calling.
// Expected result: log_trace_test.go -> TestTrace:17
func Trace() string {
	pc, path, line, err := getStackTrace()

	if err != nil {
		return ""
	}

	funcCall := runtime.FuncForPC(pc).Name()

	return fmt.Sprintf("%s -> %s:%d", getBase(path), getFuncName(getBase(funcCall)), line)
}

func getFuncName(funcBase string) string {
	funcName := strings.Split(funcBase, ".")
	l := len(funcName)

	if l == 1 {
		return funcName[0]
	}

	return funcName[1]
}

func getBase(path string) string {
	return filepath.Base(path)
}
