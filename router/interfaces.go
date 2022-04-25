package router

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

const (
	BodyFieldName            = "Body"
	PathParametersFieldName  = "PathParameters"
	QueryParametersFieldName = "QueryParameters"
)

type Request[B, P, Q any] struct {
	Context         context.Context
	Body            B
	PathParameters  P
	QueryParameters Q
	Headers         http.Header
}

type OperationHandler interface {
	requestTypes() (map[string]reflect.Type, error)
	responseBodyType() (reflect.Type, error)
	handler(openapi, ...http.Handler) http.Handler
}

func NewOperationHandler[B, P, Q, R any](handlerFunc OperationHandlerFunc[B, P, Q, R]) OperationHandler {
	return handlerFunc
}

type OperationHandlerFunc[B, P, Q, R any] func(request Request[B, P, Q]) (R, error)

func (fn OperationHandlerFunc[B, P, Q, R]) requestTypes() (map[string]reflect.Type, error) {
	fields := []string{BodyFieldName, PathParametersFieldName, QueryParametersFieldName}
	requestTypes := make(map[string]reflect.Type, len(fields))
	fnType := reflect.TypeOf(fn)
	requestStructType := fnType.In(0)
	for _, fieldName := range []string{BodyFieldName, PathParametersFieldName, QueryParametersFieldName} {
		field, ok := requestStructType.FieldByName(fieldName)
		if !ok {
			return requestTypes, fmt.Errorf("missing %q field in %s", fieldName, requestStructType)
		}
		requestTypes[fieldName] = field.Type
	}
	return requestTypes, nil
}

//SuccessfulResponseBody: controllerType.Out(0),
func (fn OperationHandlerFunc[B, P, Q, R]) responseBodyType() (reflect.Type, error) {
	fnType := reflect.TypeOf(fn)
	responseType := fnType.Out(0)
	return responseType, nil
}

func (fn OperationHandlerFunc[B, P, Q, R]) handler(oa openapi, handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		for _, h := range handlers {
			h.ServeHTTP(responseWriter, request)
		}
		request.Header.Get("Accept")
		request.Header.Get("Content-Type")
		//r := Request[B, P, Q]{
		//	Context: request.Context(),
		//	Headers: request.Header,
		//}
		//requestBytes, _ := io.ReadAll(request.Body)
		//r.Body, _ = oa.contentTypes[""].Unmarshal(requestByte)

		//response, err := fn(r)
		//if err != nil {
		//	return
		//}
		//oa.contentTypes[""].Marshal(response)
	})
}

type ContentType interface {
	Mime() string
	Marshal(value any) ([]byte, error)
	Unmarshal([]byte) (any, error)
}
