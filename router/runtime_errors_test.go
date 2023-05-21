package router

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

type contextModifier func(*Context)

func testContext(modifiers ...contextModifier) *Context {
	ctx := &Context{
		Operation: SpecOperation{
			Operation: openapi3.NewOperation(),
		},
		RawResponse: &RawResponse{},
		Request: &http.Request{
			URL:    &url.URL{},
			Header: http.Header{},
		},
		Writer: &httptest.ResponseRecorder{},
		Params: &httprouter.Params{},
	}
	for _, modifier := range modifiers {
		modifier(ctx)
	}
	return ctx
}

func withURL(t *testing.T, urlString string) contextModifier {
	urlValue, err := url.Parse(urlString)
	require.NoError(t, err)
	return func(ctx *Context) {
		ctx.Request.URL = urlValue
	}
}

func withBody(body string) contextModifier {
	bodyReader := io.NopCloser(bytes.NewBuffer([]byte(body)))

	return withBodyReader(bodyReader)
}

func withBodyReader(bodyReader io.ReadCloser) contextModifier {
	return func(ctx *Context) {
		ctx.Request.Body = bodyReader
	}
}

func withParams(params *httprouter.Params) contextModifier {
	return func(ctx *Context) {
		ctx.Params = params
	}
}

func withHeader(header string, values ...string) contextModifier {
	return func(ctx *Context) {
		for _, value := range values {
			ctx.Request.Header.Add(header, value)
		}
	}
}

func withOperation(operation *openapi3.Operation) contextModifier {
	return func(ctx *Context) {
		ctx.Operation.Operation = operation
	}
}

func withResponseWriter(writer http.ResponseWriter) contextModifier {
	return func(ctx *Context) {
		ctx.Writer = writer
	}
}

func TestNewBadRequestErr(t *testing.T) {
	cause := errors.New("test error")
	textCtx := testContext()
	testErr := newBadRequestErr(textCtx, cause, InBody)
	assert.Equal(t, textCtx, testErr.Context)
	assert.Equal(t, InBody, testErr.In)
	assert.Equal(t, cause, testErr.Err)
	assert.ErrorIs(t, testErr, cause)
	assert.ErrorIs(t, testErr.Err, cause)
	assert.ErrorIs(t, testErr, BadRequestErr{
		In:  InBody,
		Err: cause,
	})
	assert.NotErrorIs(t, testErr, BadRequestErr{
		In:  InPathParams,
		Err: cause,
	})
	assert.NotErrorIs(t, testErr, BadRequestErr{
		In: InBody,
	})
	assert.Equal(t, cause, errors.Unwrap(testErr))

	assert.ErrorContains(t, testErr, "invalid request body.")
	testErr.In = InQueryParams
	assert.ErrorContains(t, testErr, "invalid request query param.")
	testErr.In = InPathParams
	assert.ErrorContains(t, testErr, "invalid request path param.")
}

type ErrorResponse struct {
	OK      string `status:"200"`
	Message string `status:"400"`
}

func TestErrorHandler(t *testing.T) {
	errorHandler := ErrorHandler(func(c *Context, err error) (Response[ErrorResponse], error) {
		return SendText(ErrorResponse{Message: err.Error()}).Status(400), nil
	})
	assert.Equal(t, utils.NilType, errorHandler.requestTypes().requestBody)
	assert.Equal(t, utils.NilType, errorHandler.requestTypes().pathParams)
	assert.Equal(t, utils.NilType, errorHandler.requestTypes().queryParams)
	assert.Equal(t, handlerResponses{
		200: {
			status:       200,
			responseType: reflect.TypeOf(""),
			fieldIndex:   []int{0},
			isNilType:    false,
		},
		400: {
			status:       400,
			responseType: reflect.TypeOf(""),
			fieldIndex:   []int{1},
			isNilType:    false,
		},
	}, errorHandler.responseTypes())

	successResponse := RawResponse{
		Status:      200,
		ContentType: "text/plain",
		Body:        []byte("ok"),
		Headers:     http.Header{},
	}
	handlerFunc := errorHandler.handlerFactory(openapi{contentTypes: ContentTypes{
		"text/plain": PlainTextContentType{},
	}}, func(c *Context) (RawResponse, error) {
		if c.Params.ByName("foo") == "bar" {
			return RawResponse{}, errors.New("foo can not be bar")
		}
		*c.RawResponse = successResponse
		return successResponse, nil
	})
	resp, err := handlerFunc(testContext())
	require.NoError(t, err)
	assert.Equal(t, successResponse, resp)

	testContextWithParam := testContext()
	testContextWithParam.Params = &httprouter.Params{{
		Key: "foo", Value: "bar",
	}}
	resp, err = handlerFunc(testContextWithParam)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Status)
	assert.Equal(t, "foo can not be bar", string(resp.Body))
	assert.Equal(t, "text/plain", resp.ContentType)
}
