package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
