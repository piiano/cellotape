package router

import (
	"net/http"
	"reflect"
)

type Handler interface {
	requestTypes() requestTypes
	responseTypes() handlerResponses
	sourcePosition() sourcePosition
	handler(oa openapi, next groupHandlerFunc[any]) groupHandlerFunc[any]
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
