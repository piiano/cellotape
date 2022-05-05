package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
)

type groupHandlerFunc[R any] func(HandlerContext) (Response[R], error)

func NewHandler[R any](h groupHandlerFunc[R]) Handler {
	return h
}

// Currently groupHandlerFunc don't have support for typed input
// TODO: should groupHandlerFunc be able to handle cross cutting requests pathParams?
func (h groupHandlerFunc[R]) requestTypes() requestTypes { return requestTypes{} }

// responseTypes Extract the responses defined by the groupHandlerFunc and returns handlerResponses
func (h groupHandlerFunc[R]) responseTypes() handlerResponses {
	var r R
	responseType := reflect.TypeOf(r)
	return extractResponses(responseType)
}

func (h groupHandlerFunc[R]) sourcePosition() sourcePosition {
	return functionSourcePosition(h)
}

func (h groupHandlerFunc[R]) handler(oa openapi, next groupHandlerFunc[any]) groupHandlerFunc[any] {
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

type HandlerContext struct {
	Request *http.Request
	Params  httprouter.Params
	Next    groupHandlerFunc[any]
}

//type handlerFunc[B, P, Q, R any] func(Request[B, P, Q], HandlerContext2) (Response[R], error)
//
//func (fn handlerFunc[B, P, Q, R]) boundHandler(originalRequest originalRequest, next handlerFunc[any, any, any, any]) handlerFunc[B, P, Q, R] {
//	var binders handlerBinders2[B, P, Q, R]
//	return func(r Request[B, P, Q], c HandlerContext2) (Response[R], error) {
//		c.next = next
//		c.httpRequest = originalRequest.Request
//		c.pathParams = originalRequest.Params
//		request, err := binders.requestBinder(c)
//		if err != nil {
//			return Response[R]{}, err
//		}
//		response, err := fn(request)
//		if err != nil {
//			return Response[R]{}, err
//		}
//		return binders.responseBinder(c, response)
//	}
//}
//
//type originalRequest struct {
//	Request *http.Request
//	Params  httprouter.Params
//}
//
//type HandlerContext2[B, P, Q, R any] struct {
//	httpRequest *http.Request
//	pathParams  httprouter.Params
//	next        handlerFunc[B, P, Q, R]
//}
//
//type handlerBinders2[B, P, Q, R any] struct {
//	requestBinder  func(HandlerContext2) (Request[B, P, Q], error)
//	responseBinder func(HandlerContext2, Response[R]) (Response[R], error)
//}
