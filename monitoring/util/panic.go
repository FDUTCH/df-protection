package util

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// PanicMessage returns similar to original panic stack trace.
// You may call it inside your Recovery function to get panic stack trace.
func PanicMessage(err error) string {
	data := fmt.Sprintf("panic: %v\n\n", err.Error())
	stack := debug.Stack()
	arr := strings.Split(string(stack), "\n")

	for i, line := range arr {
		// ensuring this is the original panic() function.
		if strings.HasPrefix(line, "panic(") {
			return data + strings.Join(arr[i+2:], "\n")
		}
	}
	return "panic did not happen"
}
