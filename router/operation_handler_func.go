package router

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type OperationHandler interface {
	types() operationTypes
	asGinHandler(oa OpenAPI) gin.HandlerFunc
}

type operationTypes struct {
	RequestTypes
	responsesType reflect.Type
}

type operationFunc[B, P, Q, R any] func(Request[B, P, Q], Send[R])

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
		RequestTypes: RequestTypes{
			requestBody: reflect.TypeOf(body),
			pathParams:  reflect.TypeOf(path),
			queryParams: reflect.TypeOf(query),
		},
		responsesType: reflect.TypeOf(responses),
	}
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
			binders.sendError(c, err)
			return
		}
		if err := binders.pathParamsBinder(c, &(request.PathParams)); err != nil {
			binders.sendError(c, err)
			return
		}
		if err := binders.queryParamsBinder(c, &(request.QueryParams)); err != nil {
			binders.sendError(c, err)
			return
		}
		sendFunc := binders.sendFactory(c)
		fn(request, sendFunc)
	}
}
