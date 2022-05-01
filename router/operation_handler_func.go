package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type operation struct {
	id               string
	handlers         []http.HandlerFunc
	operationHandler OperationHandler
}

type OperationHandler interface {
	types() operationTypes
	asGinHandler(oa openapi) gin.HandlerFunc
}

type operationTypes struct {
	requestBody   reflect.Type
	pathParams    reflect.Type
	queryParams   reflect.Type
	responsesType reflect.Type
}

type operationFunc[B, P, Q, R any] func(Request[B, P, Q]) (int, R)

type Request[B, P, Q any] struct {
	Context     context.Context
	Body        B
	PathParams  P
	QueryParams Q
	Headers     http.Header
}

type Nil *uintptr

var nilValue Nil
var nilType = reflect.TypeOf(nilValue)

func OperationFunc[B, P, Q, R any](opFunc operationFunc[B, P, Q, R]) OperationHandler {
	return opFunc
}

func (fn operationFunc[B, P, Q, R]) types() operationTypes {
	var (
		body      B
		path      P
		query     Q
		responses R
	)
	return operationTypes{
		requestBody:   reflect.TypeOf(body),
		pathParams:    reflect.TypeOf(path),
		queryParams:   reflect.TypeOf(query),
		responsesType: reflect.TypeOf(responses),
	}
}

func (fn operationFunc[B, P, Q, R]) asGinHandler(oa openapi) gin.HandlerFunc {
	binders := bindersFactory[B, P, Q, R](oa, fn)
	return func(c *gin.Context) {
		request, err := binders.requestBinder(c)
		if err != nil {
			binders.errorBinder(c, err)
			return
		}
		status, responses := fn(request)
		binders.responseBinder(c, status, responses)
	}
}
