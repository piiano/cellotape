package restcontroller

import (
	"reflect"
)

type Params[B, P, Q any] struct {
	Body    B
	Path    P
	Query   Q
	Headers map[string][]string
}

type Controller[B, P, Q, R any] func(params Params[B, P, Q]) (R, error)

type ControllerTypeInfo struct {
	PathParams             reflect.Type
	QueryParams            reflect.Type
	RequestBody            reflect.Type
	SuccessfulResponseBody reflect.Type
}
