package router

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
)

type operation struct {
	id               string
	handlers         []Handler
	operationHandler Handler
	responseTypes    handlerResponseTypes
}

type Handler interface {
	requestTypes() handlerRequestTypes
	responseTypes() handlerResponseTypes
	handler(oa openapi, next func(HandlerContext) (Response[any], error)) func(HandlerContext) (Response[any], error)
}

type HandlerContext struct {
	Request *http.Request
	Params  httprouter.Params
	Next    func(HandlerContext) (Response[any], error)
}

type handlerRequestTypes struct {
	requestBody reflect.Type
	pathParams  reflect.Type
	queryParams reflect.Type
}

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
	// indicate whether the response has been bound to by a response binder.
	// A response passing through the middleware chain should be bound only once.
	bound bool
}

// NextResponse calls next and return its response and error as if they were from the current handler response type
func NextResponse[R any](c HandlerContext) (Response[R], error) {
	resp, err := c.Next(c)
	return Response[R]{
		Status:  resp.Status,
		Headers: resp.Headers,
		Bytes:   resp.Bytes,
		bound:   resp.bound,
	}, err
}
func Send[R any](status int, response R, headers ...http.Header) (Response[R], error) {
	aggregatedHeaders := make(http.Header, 0)
	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				aggregatedHeaders.Add(key, value)
			}
		}
	}
	return Response[R]{
		Status:   status,
		Response: response,
		Headers:  aggregatedHeaders,
	}, nil
}
func Error[R any](err error) (Response[R], error) {
	return Response[R]{}, err
}

type Nil *uintptr

var nilValue Nil
var nilType = reflect.TypeOf(nilValue)

func NewOperationHandler[B, P, Q, R any](h operationFunc[B, P, Q, R]) Handler {
	return h
}
func NewHandler[R any](h handlerFunc[R]) Handler {
	return h
}

func (fn operationFunc[B, P, Q, R]) responseTypes() handlerResponseTypes {
	var r R
	responseType := reflect.TypeOf(r)
	return handlerResponseTypes{
		originalType:      responseType,
		declaredResponses: extractResponses(responseType),
		allowAny:          false,
	}
}
func (fn operationFunc[B, P, Q, R]) requestTypes() handlerRequestTypes {
	var body B
	var path P
	var query Q
	return handlerRequestTypes{
		requestBody: reflect.TypeOf(body),
		pathParams:  reflect.TypeOf(path),
		queryParams: reflect.TypeOf(query),
	}
}

func (fn operationFunc[B, P, Q, R]) handler(oa openapi, _ func(HandlerContext) (Response[any], error)) func(HandlerContext) (Response[any], error) {
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

type handlerFunc[R any] func(HandlerContext) (Response[R], error)

func (h handlerFunc[R]) requestTypes() handlerRequestTypes {
	return handlerRequestTypes{}
}

func (h handlerFunc[R]) responseTypes() handlerResponseTypes {
	var r R
	responseType := reflect.TypeOf(r)
	return handlerResponseTypes{
		originalType:      responseType,
		declaredResponses: extractResponses(responseType),
		allowAny:          true,
	}
}

func (h handlerFunc[R]) handler(oa openapi, next func(HandlerContext) (Response[any], error)) func(HandlerContext) (Response[any], error) {
	responseBinder := responseBinderFactory[R](h.responseTypes(), oa.contentTypes)
	return func(c HandlerContext) (Response[any], error) {
		c.Next = next
		resp, err := h(c)
		if err != nil {
			return Response[any]{}, err
		}
		return responseBinder(c, resp)
	}
}
