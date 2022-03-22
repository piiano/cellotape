package restcontroller

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Controller interface {
	typeInfo() controllerTypeInfo
	GinHandler() gin.HandlerFunc
	OpenAPIOperation(ID string, options *OperationOptions) (*openapi3.Operation, error)
}

// ControllerFn A REST controller function with generic parameters and response type
type ControllerFn[B, P, Q, R any] func(params Params[B, P, Q]) (R, error)

// NewController A wrapper for declaring controller without duplicate declaration of all generic parameters
func NewController[B, P, Q, R any](f ControllerFn[B, P, Q, R]) Controller {
	return f
}

// Params Generic REST Parameters of controller function
type Params[B, P, Q any] struct {
	Body    B
	Path    P
	Query   Q
	Headers map[string][]string
}

// ControllerTypeInfo Describe type information of a controller
type controllerTypeInfo struct {
	PathParams             reflect.Type
	QueryParams            reflect.Type
	RequestBody            reflect.Type
	SuccessfulResponseBody reflect.Type
}

// TypeInfo Extract type information from a controller using reflections.
func (fn ControllerFn[B, P, Q, R]) typeInfo() controllerTypeInfo {
	controllerType := reflect.TypeOf(fn)
	paramsType := controllerType.In(0)
	return controllerTypeInfo{
		RequestBody:            paramsType.Field(0).Type,
		PathParams:             paramsType.Field(1).Type,
		QueryParams:            paramsType.Field(2).Type,
		SuccessfulResponseBody: controllerType.Out(0),
	}
}
