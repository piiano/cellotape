package router

import (
	"fmt"
	"reflect"
	"runtime"
)

func functionSourcePosition(i any) sourcePosition {
	t := reflect.TypeOf(i)
	if t == nil {
		return sourcePosition{}
	}
	if t.Kind() == reflect.Pointer {
		return functionSourcePosition(reflect.ValueOf(t).Elem().Interface())
	}
	if t.Kind() != reflect.Func {
		return sourcePosition{}
	}
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	file, line := fn.FileLine(fn.Entry())
	return sourcePosition{file: file, line: line, ok: true}
}

type sourcePosition struct {
	ok   bool
	file string
	line int
}

func (sp sourcePosition) String() string {
	return fmt.Sprintf(".%s:%d", sp.file, sp.line)
}
