package router

import (
	"context"
	"net/http"
	"reflect"
)

type operationFunc[B, P, Q, R any] func(Request[B, P, Q]) (Response[R], error)

type Request[B, P, Q any] struct {
	Context     context.Context
	Body        B
	PathParams  P
	QueryParams Q
	Headers     http.Header
}
type Response[R any] struct {
	Status   int
	Headers  http.Header
	Response R
	Bytes    []byte
	// indicate whether the response has been bound by a response binder.
	// A response passing through the middleware chain should be bound only once.
	bound bool
}

func NewOperationHandler[B, P, Q, R any](h operationFunc[B, P, Q, R]) Handler {
	return h
}
func (fn operationFunc[B, P, Q, R]) responseTypes() handlerResponses {
	var r R
	responseType := reflect.TypeOf(r)
	return extractResponses(responseType)
}
func (fn operationFunc[B, P, Q, R]) requestTypes() requestTypes {
	var body B
	var path P
	var query Q
	return requestTypes{
		requestBody: reflect.TypeOf(body),
		pathParams:  reflect.TypeOf(path),
		queryParams: reflect.TypeOf(query),
	}
}

func (fn operationFunc[B, P, Q, R]) sourcePosition() sourcePosition {
	return functionSourcePosition(fn)
}

func (fn operationFunc[B, P, Q, R]) handler(oa openapi, _ groupHandlerFunc[any]) groupHandlerFunc[any] {
	binders := bindersFactory[B, P, Q, R](oa, fn)
	return func(c HandlerContext) (Response[any], error) {
		request, err := binders.requestBinder(c)
		if err != nil {
			return Response[any]{}, err
		}
		response, err := fn(request)
		if err != nil {
			return Response[any]{}, err
		}
		return binders.responseBinder(c, response)
	}
}
