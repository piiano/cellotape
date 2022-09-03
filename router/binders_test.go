package router

import (
	"bytes"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestResponseContentType(t *testing.T) {
	_, err := responseContentType("no-such-content-type", http.Header{}, DefaultContentTypes(), JSONContentType{})
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
	requestURL, err := url.Parse("http:0.0.0.0:90/abc?Foo=42")
	require.NoError(t, err)
	err = queryBinder(&http.Request{
		URL: requestURL,
	}, &params)
	require.NoError(t, err)
	assert.Equal(t, StructType{Foo: 42}, params)
}

func TestQueryBinderFactoryError(t *testing.T) {
	queryBinder := queryBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	requestURL, err := url.Parse("http:0.0.0.0:90/abc?Foo=abc")
	require.NoError(t, err)
	err = queryBinder(&http.Request{
		URL: requestURL,
	}, &params)

	require.Error(t, err)
}

func TestPathBinderFactory(t *testing.T) {
	pathBinder := pathBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := pathBinder(&httprouter.Params{{
		Key:   "Foo",
		Value: "42",
	}}, &params)
	require.NoError(t, err)
	assert.Equal(t, StructType{Foo: 42}, params)
}

func TestPathBinderFactoryError(t *testing.T) {
	pathBinder := pathBinderFactory[StructType](reflect.TypeOf(StructType{}))
	var params StructType
	err := pathBinder(&httprouter.Params{{
		Key:   "Foo",
		Value: "bar",
	}}, &params)
	require.Error(t, err)
}

func TestRequestBodyBinderFactory(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(&http.Request{
		Body: io.NopCloser(bytes.NewBuffer([]byte("42"))),
	}, &param)
	require.NoError(t, err)
	assert.Equal(t, 42, param)
}

func TestRequestBodyBinderFactoryError(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(&http.Request{
		Body: io.NopCloser(bytes.NewBuffer([]byte(`"foo"`))),
	}, &param)
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
	err := requestBodyBinder(&http.Request{
		Body: io.NopCloser(readerWithError(`42`)),
	}, &param)
	require.Error(t, err)
}

func TestRequestBodyBinderFactoryContentTypeError(t *testing.T) {
	requestBodyBinder := requestBodyBinderFactory[int](reflect.TypeOf(0), DefaultContentTypes())
	var param int
	err := requestBodyBinder(&http.Request{
		Header: http.Header{"Content-Type": {"no-such-content-type"}},
		Body:   io.NopCloser(bytes.NewBuffer([]byte(`42`))),
	}, &param)
	require.Error(t, err)
}
