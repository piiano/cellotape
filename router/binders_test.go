package router

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestResponseContentType(t *testing.T) {
	_, err := responseContentType("no-such-content-type", DefaultContentTypes(), JSONContentType{})
	require.ErrorIs(t, err, UnsupportedResponseContentTypeErr)
}

func TestRequestContentType(t *testing.T) {
	_, err := requestContentType(&http.Request{
		Header: http.Header{"Content-Type": {"no-such-content-type"}},
	}, DefaultContentTypes(), JSONContentType{})
	require.ErrorIs(t, err, UnsupportedRequestContentTypeErr)
}

type StructType struct {
	Foo int
}

func TestQueryBinderFactory(t *testing.T) {
	queryBinder := queryBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := queryBinder(testContext(withURL(t, "http:0.0.0.0:90/abc?Foo=42")), &params)
	require.NoError(t, err)
	assert.Equal(t, StructType{Foo: 42}, params)
}

type StructWithArrayType struct {
	Foo []int
}

func TestQueryBinderFactoryWithArrayType(t *testing.T) {
	queryBinder := queryBinderFactory[StructWithArrayType](reflect.TypeOf(StructWithArrayType{}))
	var params StructWithArrayType
	err := queryBinder(testContext(withURL(t, "http:0.0.0.0:90/abc?Foo=42&Foo=6&Foo=7")), &params)
	require.NoError(t, err)
	assert.Equal(t, StructWithArrayType{Foo: []int{42, 6, 7}}, params)
}

func TestQueryBinderFactoryMultipleParamToNonArrayError(t *testing.T) {
	queryBinder := queryBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := queryBinder(testContext(withURL(t, "http:0.0.0.0:90/abc?Foo=42&Foo=6&Foo=7")), &params)
	require.Error(t, err)
}

func TestQueryBinderFactoryError(t *testing.T) {
	queryBinder := queryBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := queryBinder(testContext(withURL(t, "http:0.0.0.0:90/abc?Foo=abc")), &params)
	require.Error(t, err)
}

func TestPathBinderFactory(t *testing.T) {
	pathBinder := pathBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := pathBinder(testContext(withParams(&httprouter.Params{{
		Key:   "Foo",
		Value: "42",
	}})), &params)
	require.NoError(t, err)
	assert.Equal(t, StructType{Foo: 42}, params)
}

func TestPathBinderFactoryError(t *testing.T) {
	pathBinder := pathBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := pathBinder(testContext(withParams(&httprouter.Params{{
		Key:   "Foo",
		Value: "bar",
	}})), &params)
	require.Error(t, err)
}

func TestRequestBodyBinderFactory(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(testContext(withBody("42")), &param)
	require.NoError(t, err)
	assert.Equal(t, 42, param)
}

func TestRequestBodyBinderFactoryWithSchema(t *testing.T) {
	testOp := openapi3.NewOperation()
	testOp.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewIntegerSchema()),
	}
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(testContext(
		withBody("42"),
		withHeader("Content-Type", "application/json"),
		withOperation(testOp)), &param)
	require.NoError(t, err)
	assert.Equal(t, 42, param)
}

func TestRequestBodyBinderFactoryError(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int

	err := requestBodyBinder(testContext(withBody(`"foo"`)), &param)
	require.Error(t, err)
}

type readerWithError []byte

var readerError = errors.New("error")

func (r readerWithError) Read(_ []byte) (int, error) {
	return 0, readerError
}

func TestRequestBodyBinderFactoryReaderError(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(testContext(
		withBodyReader(io.NopCloser(readerWithError(`42`)))), &param)
	require.Error(t, err)
}

func TestRequestBodyBinderFactoryContentTypeError(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int

	err := requestBodyBinder(testContext(
		withBody("42"),
		withHeader("Content-Type", "no-such-content-type")), &param)
	require.Error(t, err)
}

func TestRequestBodyBinderFactoryContentTypeWithCharset(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(testContext(
		withBody("42"),
		withHeader("Content-Type", "application/json; charset=utf-8")), &param)
	require.NoError(t, err)
	assert.Equal(t, 42, param)
}

func TestRequestBodyBinderFactoryInvalidContentType(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(testContext(
		withBody("42"),
		withHeader("Content-Type", "invalid content type")), &param)
	require.Error(t, err)
}

func TestRequestBodyBinderFactoryContentTypeAnyWithCharset(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(testContext(
		withBody("42"),
		withHeader("Content-Type", "*/*; charset=utf-8")), &param)
	require.NoError(t, err)
	assert.Equal(t, 42, param)
}

type CollidingFieldsParam1 struct {
	Value string `form:"param1"`
}
type CollidingFieldsParam2 struct {
	Value string `form:"param2"`
}
type CollidingFieldsParams struct {
	CollidingFieldsParam1
	CollidingFieldsParam2
}

func TestBindingEmbeddedQueryParamsCollidingFields(t *testing.T) {
	requestBodyBinder := queryBinderFactory[CollidingFieldsParams](reflect.TypeOf(CollidingFieldsParams{}))
	var param CollidingFieldsParams

	ctx := testContext(withURL(t, "http://http:0.0.0.0:8080/path?param1=foo&param2=bar"))

	err := requestBodyBinder(ctx, &param)
	require.NoError(t, err)
	require.Equal(t, "foo", param.CollidingFieldsParam1.Value)
	require.Equal(t, "bar", param.CollidingFieldsParam2.Value)
}

type CollidingParamString struct {
	Value string `form:"param1"`
}
type CollidingParamInt struct {
	Value int `form:"param1"`
}
type CollidingParams struct {
	CollidingParamString
	CollidingParamInt
}

func TestBindingEmbeddedQueryParamsCollidingParams(t *testing.T) {
	requestBodyBinder := queryBinderFactory[CollidingParams](reflect.TypeOf(CollidingParams{}))

	var param CollidingParams
	err := requestBodyBinder(testContext(
		withURL(t, "http://http:0.0.0.0:8080/path?param1=42")), &param)
	require.NoError(t, err)
	require.Equal(t, "42", param.CollidingParamString.Value)
	require.Equal(t, 42, param.CollidingParamInt.Value)
}

type errWriter struct{}

func (e errWriter) Header() http.Header { return http.Header{} }
func (e errWriter) WriteHeader(int)     {}
func (e errWriter) Write(i []byte) (int, error) {
	return 0, errors.New("error")
}

func TestErrOnWriterError(t *testing.T) {
	type R = OKResponse[string]
	responses := extractResponses(utils.GetType[R]())
	binder := responseBinderFactory[R](responses, DefaultContentTypes())
	response := SendOK(R{OK: "foo"}).ContentType("unknown")

	testCases := []struct {
		name      string
		writer    http.ResponseWriter
		assertion func(require.TestingT, error, ...any)
	}{
		{
			name:      "writer error",
			writer:    errWriter{},
			assertion: require.Error,
		},
		{
			name:      "proper writer",
			writer:    httptest.NewRecorder(),
			assertion: require.NoError,
		},
	}
	testOp := openapi3.NewOperation()
	testOp.AddResponse(200, openapi3.NewResponse().WithJSONSchema(openapi3.NewStringSchema()))

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctx := testContext(
				withOperation(testOp),
				withResponseWriter(test.writer),
			)
			_, err := binder(ctx, response)
			test.assertion(t, err)
		})
	}
}
