package router

import (
	"reflect"
)

type Typed interface {
	getType() reflect.Type
}
