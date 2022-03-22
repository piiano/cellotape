package restcontroller

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Params[B, P, Q any] struct {
	Body    B
	Path    P
	Query   Q
	Headers map[string][]string
}

type ControllerFn[B, P, Q, R any] func(params Params[B, P, Q]) (R, error)

type Controller interface {
	TypeInfo() ControllerTypeInfo
	GinHandler() gin.HandlerFunc
	OpenAPIOperation(ID string, options *OperationOptions) (*openapi3.Operation, error)
}

type ControllerTypeInfo struct {
	PathParams             reflect.Type
	QueryParams            reflect.Type
	RequestBody            reflect.Type
	SuccessfulResponseBody reflect.Type
}

func (fn ControllerFn[B, P, Q, R]) TypeInfo() ControllerTypeInfo {
	controllerType := reflect.TypeOf(fn)
	paramsType := controllerType.In(0)
	return ControllerTypeInfo{
		RequestBody:            paramsType.Field(0).Type,
		PathParams:             paramsType.Field(1).Type,
		QueryParams:            paramsType.Field(2).Type,
		SuccessfulResponseBody: controllerType.Out(0),
	}
}
