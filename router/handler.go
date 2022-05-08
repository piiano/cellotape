package router

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
)

type Handler interface {
	requestTypes() requestTypes
	responseTypes() handlerResponses
	sourcePosition() sourcePosition
	handlerFactory(oa openapi, next BoundHandlerFunc) BoundHandlerFunc
}

type HandlerFunc[B, P, Q, R any] func(Context, Request[B, P, Q]) (Response[R], error)

type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Params      *httprouter.Params
	NextFunc    BoundHandlerFunc
	RawResponse RawResponse
}

func (c Context) Next() (RawResponse, error) {
	return c.NextFunc(c)
}

type Request[B, P, Q any] struct {
	Context     context.Context
	Method      string
	URL         *url.URL
	Body        B
	PathParams  P
	QueryParams Q
	Headers     http.Header
}
type Response[R any] struct {
	Status   int
	Headers  http.Header
	Response R
}

type BoundHandlerFunc func(Context) (RawResponse, error)

type RawResponse struct {
	// Status written with WriteHeader
	Status int
	// buffered Body bytes written by calls to Write
	Body []byte
	// response Headers
	Headers http.Header
}

func NewHandler[B, P, Q, R any](h HandlerFunc[B, P, Q, R]) Handler {
	return h
}

// requestTypes extracts the request types defined by the HandlerFunc
func (h HandlerFunc[B, P, Q, R]) requestTypes() requestTypes {
	return requestTypes{
		requestBody: getType[B](),
		pathParams:  getType[P](),
		queryParams: getType[Q](),
	}
}

// responseTypes extracts the responses defined by the HandlerFunc and returns handlerResponses
func (h HandlerFunc[B, P, Q, R]) responseTypes() handlerResponses {
	return extractResponses(getType[R]())
}

// sourcePosition finds the sourcePosition of the HandlerFunc function for printing meaningful messages during validations
func (h HandlerFunc[B, P, Q, R]) sourcePosition() sourcePosition {
	return functionSourcePosition(h)
}

func (h HandlerFunc[B, P, Q, R]) handlerFactory(oa openapi, next BoundHandlerFunc) BoundHandlerFunc {
	bindRequest := requestBinderFactory[B, P, Q](oa, h.requestTypes())
	bindResponse := responseBinderFactory[R](h.responseTypes(), oa.contentTypes)
	return func(context Context) (RawResponse, error) {
		// when handler will be called, set the next to next
		context.NextFunc = next
		request, err := bindRequest(context.Request, context.Params)
		if err != nil {
			return RawResponse{}, err
		}
		// call current handler and all nested handlers in the chain
		response, err := h(context, request)
		// bind the response
		return bindResponse(context.Writer, context.Request, response)
	}
}

// Send is a helper function for constructing a HandlerFunc response.
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
