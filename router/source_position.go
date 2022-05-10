package router

import (
	"fmt"
	"reflect"
	"runtime"
)

// functionSourcePosition receive a function and try to extract its file and line position in sources.
// If provided input is not a valid function or of fails from any other reason then sourcePosition.ok will be false.
func functionSourcePosition(function any) sourcePosition {
	t := reflect.TypeOf(function)
	if t == nil {
		return sourcePosition{}
	}
	if t.Kind() != reflect.Func {
		return sourcePosition{}
	}
	fn := runtime.FuncForPC(reflect.ValueOf(function).Pointer())
	file, line := fn.FileLine(fn.Entry())
	return sourcePosition{file: file, line: line, ok: true}
}

type sourcePosition struct {
	ok   bool
	file string
	line int
}

// String create a string representation of the position that can be interpreted by standard terminal and tools as links
// to the actual source position of the function.
func (sp sourcePosition) String() string {
	return fmt.Sprintf("%s:%d", sp.file, sp.line)
}
