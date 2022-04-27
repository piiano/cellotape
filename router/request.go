package router

import (
	"context"
	"net/http"
	"reflect"
)

type Request[B, P, Q any] struct {
	Context         context.Context
	Body            B
	PathParameters  P
	QueryParameters Q
	Headers         http.Header
}

type Nil uintptr

var nilValue Nil
var nilType = reflect.TypeOf(nilValue)
