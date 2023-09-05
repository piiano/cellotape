package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/piiano/cellotape/router/utils"
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
type HandlerFunc[B, P, Q, R any] func(*Context, Request[B, P, Q]) (Response[R], error)

// BoundHandlerFunc is an untyped wrapper to HandlerFunc that bound internally calls to request binding, response
// binding and propagate next handler in the Context.
type BoundHandlerFunc func(*Context) (RawResponse, error)

// Context carries the original http.Request and http.ResponseWriter and additional important parameters through the
// handlers chain.
type Context struct {
	Operation   SpecOperation
	Writer      http.ResponseWriter
	Request     *http.Request
	Params      *httprouter.Params
	RawResponse *RawResponse
	NextFunc    BoundHandlerFunc
}

func (c *Context) Next() (RawResponse, error) {
	return c.NextFunc(c)
}

type Request[B, P, Q any] struct {
	Body        B
	PathParams  P
	QueryParams Q
	Headers     http.Header
}

type Response[R any] struct {
	status      int
	contentType string
	headers     http.Header
	response    R
}

// Status set the response status code.
func (r Response[R]) Status(status int) Response[R] {
	r.status = status
	return r
}

// ContentType set the response content type.
func (r Response[R]) ContentType(contentType string) Response[R] {
	r.contentType = contentType
	return r
}

// AddHeader add a response header. It appends to any existing values associated with key.
func (r Response[R]) AddHeader(key, value string) Response[R] {
	r.headers.Add(key, value)
	return r
}

// SetHeader set a response header. It replaces any existing values associated with key
func (r Response[R]) SetHeader(key, value string) Response[R] {
	r.headers.Set(key, value)
	return r
}

type RawResponse struct {
	// Status written with WriteHeader
	Status int
	// ContentType is the content type used to write the response
	ContentType string
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
		requestBody: utils.GetType[B](),
		pathParams:  utils.GetType[P](),
		queryParams: utils.GetType[Q](),
	}
}

// responseTypes extracts the responses defined by the HandlerFunc and returns handlerResponses
func (h HandlerFunc[B, P, Q, R]) responseTypes() handlerResponses {
	return extractResponses(utils.GetType[R]())
}

// sourcePosition finds the sourcePosition of the HandlerFunc function for printing meaningful messages during validations
func (h HandlerFunc[B, P, Q, R]) sourcePosition() sourcePosition {
	return functionSourcePosition(h)
}

func (h HandlerFunc[B, P, Q, R]) handlerFactory(oa openapi, next BoundHandlerFunc) BoundHandlerFunc {
	bindRequest := requestBinderFactory[B, P, Q](oa, h.requestTypes())
	bindResponse := responseBinderFactory[R](h.responseTypes(), oa.contentTypes, oa.options.DefaultOperationValidation.RuntimeValidateResponses)
	return func(context *Context) (RawResponse, error) {
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
