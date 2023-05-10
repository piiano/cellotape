package router

import (
	"net/http"

	"github.com/piiano/cellotape/router/utils"
)

type Nil = utils.Nil
type MultiType[T any] utils.MultiType[T]

// Send constructs a new Response.
func Send[R any](response R, headers ...http.Header) Response[R] {
	aggregatedHeaders := make(http.Header, 0)
	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				aggregatedHeaders.Add(key, value)
			}
		}
	}
	return Response[R]{
		response: response,
		headers:  aggregatedHeaders,
	}
}

// SendBytes constructs a new Response with an "application/octet-stream" content-type.
func SendBytes[R any](response R, headers ...http.Header) Response[R] {
	return Send(response, headers...).ContentType(OctetStreamContentType{}.Mime())
}

// SendText constructs a new Response with a "plain/text" content-type.
func SendText[R any](response R, headers ...http.Header) Response[R] {
	return Send(response, headers...).ContentType(PlainTextContentType{}.Mime())
}

// SendJSON constructs a new Response with an "application/json" content-type.
func SendJSON[R any](response R, headers ...http.Header) Response[R] {
	return Send(response, headers...).ContentType(JSONContentType{}.Mime())
}

// SendOK constructs a new Response with an OK 200 status code.
func SendOK[R any](response R, headers ...http.Header) Response[R] {
	return Send(response, headers...).Status(http.StatusOK)
}

// SendOKBytes constructs a new Response with an OK 200 status code and with an "application/octet-stream" content-type.
func SendOKBytes[R any](response R, headers ...http.Header) Response[R] {
	return SendBytes(response, headers...).Status(http.StatusOK)
}

// SendOKText constructs a new Response with an OK 200 status code and with a "plain/text" content-type.
func SendOKText[R any](response R, headers ...http.Header) Response[R] {
	return SendText(response, headers...).Status(http.StatusOK)
}

// SendOKJSON constructs a new Response with an OK 200 status code and with an "application/json" content-type.
func SendOKJSON[R any](response R, headers ...http.Header) Response[R] {
	return SendJSON(response, headers...).Status(http.StatusOK)
}

// Error is a helper function for constructing a HandlerFunc error response.
func Error[R any](err error) (Response[R], error) {
	return Response[R]{}, err
}

// RawHandler adds a handler that doesn't define any type information.
func RawHandler(f func(c *Context) error) Handler {
	return NewHandler(func(c *Context, _ Request[utils.Nil, utils.Nil, utils.Nil]) (Response[any], error) {
		return Response[any]{}, f(c)
	})
}
