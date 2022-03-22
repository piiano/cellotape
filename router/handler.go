package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
)

// Handler described the HandlerFunc in a non parametrized way.
// The router "magic" of being a collection of many HandlerFunc with each having a different generic parameter is
// possible by describing each HandlerFunc with this unified Handler interface.
// Eventually a Handler is anything that can report its requestTypes, handlerResponses, sourcePosition and can become a
// BoundHandlerFunc by calling the handlerFactory method.
type Handler interface {
	requestTypes() requestTypes
	responseTypes() handlerResponses
	sourcePosition() sourcePosition
	handlerFactory(oa openapi, next BoundHandlerFunc) BoundHandlerFunc
}

// HandlerFunc is the typed handler that declare explicitly all types of request and responses.
// Check the repo examples for seeing how the HandlerFunc can be used to define Handler for operations and middlewares.
type HandlerFunc[B, P, Q, R any] func(Context, Request[B, P, Q]) (Response[R], error)

// BoundHandlerFunc is an untyped wrapper to HandlerFunc that bound internally calls to request binding, response
// binding and propagate next handler in the Context.
type BoundHandlerFunc func(Context) (RawResponse, error)

// Context carries the original http.Request and http.ResponseWriter and additional important parameters through the
// handlers chain.
type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Params      *httprouter.Params
	RawResponse *RawResponse
	NextFunc    BoundHandlerFunc
	keyValues   map[string]any
}

func (c Context) Next() (RawResponse, error) {
	return c.NextFunc(c)
}
func (c Context) Get(key string) (any, bool) {
	value, ok := c.keyValues[key]
	return value, ok
}
func (c *Context) Set(key string, value any) {
	c.keyValues[key] = value
}

type Request[B, P, Q any] struct {
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

// getType returns reflect.Type of the generic parameter it receives.
func getType[T any]() reflect.Type { return reflect.TypeOf(new(T)).Elem() }

// Nil represents an empty type.
// You can use it with the HandlerFunc generic parameters to declare no Request with no request body, no path or query
// params, or responses with no response body.
type Nil *uintptr

// nilType represent the type of Nil.
var nilType = getType[Nil]()

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
		request, err := bindRequest(context)
		if err != nil {
			return RawResponse{}, err
		}
		// call current handler and all nested handlers in the chain
		response, err := h(context, request)
		if err != nil {
			return *context.RawResponse, err
		}
		// bind the response
		return bindResponse(context, response)
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

// Error is a helper function for constructing a HandlerFunc error response.
func Error[R any](err error) (Response[R], error) {
	return Response[R]{}, err
}
