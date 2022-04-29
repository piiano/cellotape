package router

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type OperationHandler interface {
	requestTypes() (RequestTypes, error)
	responseTypes() (ResponseType, error)
	types() (OperationTypes, error)
	asGinHandler(oa OpenAPI) gin.HandlerFunc
}

type OperationTypes struct {
	RequestTypes
	ResponseType
}

type operationFunc[B, P, Q, R any] func(Request[B, P, Q]) (R, error)

func OperationFunc[B, P, Q, R any](opFunc func(Request[B, P, Q]) (R, error)) OperationHandler {
	return operationFunc[B, P, Q, R](opFunc)
}

func (fn operationFunc[B, P, Q, R]) requestTypes() (RequestTypes, error) {
	return requestTypes(fn)
}

func (fn operationFunc[B, P, Q, R]) responseTypes() (ResponseType, error) {
	fnType := reflect.TypeOf(fn)
	successType := fnType.Out(0)
	return ResponseType(successType), nil
}

func (fn operationFunc[B, P, Q, R]) types() (OperationTypes, error) {
	operationTypes := OperationTypes{}
	var err error
	if operationTypes.RequestTypes, err = fn.requestTypes(); err != nil {
		return operationTypes, err
	}
	if operationTypes.ResponseType, err = fn.responseTypes(); err != nil {
		return operationTypes, err
	}
	return operationTypes, nil
}

type HandlerContentTypes struct {
	request  ContentType
	response ContentType
}

func (fn operationFunc[B, P, Q, R]) asGinHandler(oa OpenAPI) gin.HandlerFunc {
	binders := bindersFactory[B, P, Q, R](oa, fn)
	return func(c *gin.Context) {
		var request = Request[B, P, Q]{
			Context: c,
			Headers: c.Request.Header,
		}
		if err := binders.requestBodyBinder(c, &(request.Body)); err != nil {
			binders.errorResponseBinder(c, &err)
			return
		}
		if err := binders.pathParamsBinder(c, &(request.PathParams)); err != nil {
			binders.errorResponseBinder(c, &err)
			return
		}
		if err := binders.queryParamsBinder(c, &(request.QueryParams)); err != nil {
			binders.errorResponseBinder(c, &err)
			return
		}
		response, err := fn(request)
		if err != nil {
			binders.errorResponseBinder(c, &err)
			return
		}
		binders.successfulResponseBinder(c, &response)
	}
}
