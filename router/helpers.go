package router

import "reflect"

func getType[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

type Nil *uintptr

var nilType = getType[Nil]()
