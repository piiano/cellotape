package router

import (
	"errors"
	"fmt"
)

// Runtime errors causes
var (
	UnsupportedRequestContentTypeErr  = errors.New("unsupported request content type for operation")
	UnsupportedResponseContentTypeErr = errors.New("unsupported response content type for operation")
	UnsupportedResponseStatusErr      = errors.New("unsupported response status for operation")
)

// In is a location of a request binding error
type In int

const (
	InBody In = iota
	InPathParams
	InQueryParams
)

func (in In) String() string {
	inString := ""
	switch in {
	case InBody:
		inString = "body"
	case InPathParams:
		inString = "path param"
	case InQueryParams:
		inString = "query param"
	}
	return inString
}

// BadRequestErr is the error returned when there is an error binding the request.
// You can handle this request using an ErrorHandler middleware to return a custom HTTP response.
type BadRequestErr struct {
	Err     error
	In      In
	Context Context
}

// newBadRequestErr returns a new BadRequestErr.
// BadRequestErr is the error returned when there is an error binding the request.
// You can handle this request using an ErrorHandler middleware to return a custom HTTP response.
func newBadRequestErr(ctx Context, err error, in In) BadRequestErr {
	return BadRequestErr{
		Err:     err,
		In:      in,
		Context: ctx,
	}
}

func (e BadRequestErr) Error() string {
	return fmt.Sprintf("invalid request %s. %s", e.In, e.Err)
}

func (e BadRequestErr) Is(err error) bool {
	if badRequestErr, ok := err.(BadRequestErr); ok {
		return e.In == badRequestErr.In && errors.Is(e.Err, badRequestErr.Err)
	}
	return errors.Is(e.Err, err)
}

func (e BadRequestErr) Unwrap() error {
	return e.Err
}

// ErrorHandler allows providing a handler function that can handle errors occurred in the handlers chain.
// This type of handler is particularly useful for handling BadRequestErr caused by a request binding errors and
// translate it to an HTTP response.
func ErrorHandler[R any](errHandler func(c Context, err error) (Response[R], error)) HandlerFunc[Nil, Nil, Nil, R] {
	return func(c Context, _ Request[Nil, Nil, Nil]) (Response[R], error) {
		_, err := c.Next()
		if err != nil {
			return errHandler(c, err)
		}
		return Error[R](err)
	}
}
