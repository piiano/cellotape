package router

import "reflect"

func getType[T any]() reflect.Type {
	return reflect.TypeOf(*new(T))
}

type Nil *uintptr

var nilType = getType[Nil]()
