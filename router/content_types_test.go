package router

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestContentTypeMime(t *testing.T) {
	testCases := []struct {
		contentType ContentType
		mime        string
		value       any
		bytes       []byte
	}{
		{
			contentType: OctetStreamContentType{},
			mime:        "application/octet-stream",
			value:       []byte("foo"),
			bytes:       []byte("foo"),
		},
		{
			contentType: OctetStreamContentType{},
			mime:        "application/octet-stream",
			value:       nil,
			bytes:       []byte{},
		},
		{
			contentType: PlainTextContentType{},
			mime:        "text/plain",
			value:       nil,
			bytes:       []byte{},
		},
		{
			contentType: PlainTextContentType{},
			mime:        "text/plain",
			value:       "foo",
			bytes:       []byte("foo"),
		},
		{
			contentType: JSONContentType{},
			mime:        "application/json",
			value:       nil,
			bytes:       []byte(`null`),
		},
		{
			contentType: JSONContentType{},
			mime:        "application/json",
			value:       "foo",
			bytes:       []byte(`"foo"`),
		},
	}
	for _, test := range testCases {
		assert.Equal(t, test.mime, test.contentType.Mime())
		bytes, err := test.contentType.Encode(test.value)
		require.NoError(t, err)
		assert.Equal(t, test.bytes, bytes)
		value := test.value
		err = test.contentType.Decode(test.bytes, &value)
		require.NoError(t, err)
		assert.Equal(t, test.value, value)
	}
}

func TestOctetStreamContentTypeBytesSlice(t *testing.T) {
	encodedBytes, err := OctetStreamContentType{}.Encode([]byte("foo"))
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), encodedBytes)
	var decodedBytes []byte
	err = OctetStreamContentType{}.Decode([]byte("foo"), &decodedBytes)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), decodedBytes)
}

func TestOctetStreamContentTypeError(t *testing.T) {
	_, err := OctetStreamContentType{}.Encode("foo")
	require.Error(t, err)
	var value string
	err = OctetStreamContentType{}.Decode([]byte("foo"), &value)
	require.Error(t, err)
}

func TestOctetStreamContentTypeSchemaCompatability(t *testing.T) {
	l := utils.NewInMemoryLogger()
	err := OctetStreamContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf([]byte{}),
		*openapi3.NewStringSchema().WithFormat("binary"))
	require.NoError(t, err)
	assert.Equal(t, 0, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)

	l = utils.NewInMemoryLogger()
	err = OctetStreamContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf(""),
		*openapi3.NewStringSchema().WithFormat("binary"))
	require.Error(t, err)
	assert.Equal(t, 1, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)

	l = utils.NewInMemoryLogger()
	err = OctetStreamContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf([]byte{}),
		*openapi3.NewStringSchema().WithFormat("base64"))
	require.Error(t, err)
	assert.Equal(t, 1, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)

	l = utils.NewInMemoryLogger()
	err = OctetStreamContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf([]byte{}),
		*openapi3.NewIntegerSchema())
	require.Error(t, err)
	assert.Equal(t, 1, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)
}

func TestPlainTextContentTypeString(t *testing.T) {
	encodedString, err := PlainTextContentType{}.Encode("foo")
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), encodedString)
	var decodedString string
	err = PlainTextContentType{}.Decode([]byte("foo"), &decodedString)
	require.NoError(t, err)
	assert.Equal(t, "foo", decodedString)
}

func TestPlainTextContentTypeError(t *testing.T) {
	_, err := PlainTextContentType{}.Encode(5)
	require.Error(t, err)
	var value int
	err = PlainTextContentType{}.Decode([]byte("foo"), &value)
	require.Error(t, err)
}

func TestPlainTextContentTypeSchemaCompatability(t *testing.T) {
	l := utils.NewInMemoryLogger()
	err := PlainTextContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf(""),
		*openapi3.NewStringSchema())
	require.NoError(t, err)
	assert.Equal(t, 0, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)

	l = utils.NewInMemoryLogger()
	err = PlainTextContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.PointerTo(reflect.TypeOf("")),
		*openapi3.NewStringSchema())
	require.NoError(t, err)
	assert.Equal(t, 0, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)

	l = utils.NewInMemoryLogger()
	err = PlainTextContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf([]byte{}),
		*openapi3.NewStringSchema())
	require.Error(t, err)
	assert.Equal(t, 1, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)

	l = utils.NewInMemoryLogger()
	err = PlainTextContentType{}.ValidateTypeSchema(
		l, utils.Error,
		reflect.TypeOf(""),
		*openapi3.NewIntegerSchema())
	require.Error(t, err)
	assert.Equal(t, 1, l.Counters().Errors)
	assert.Equal(t, 0, l.Counters().Warnings)
}
