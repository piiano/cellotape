package router

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

type OKResponse[R any] struct {
	OK R `status:"200"`
}

func TestSimpleSend(t *testing.T) {
	response := Send(OKResponse[string]{OK: "ok"})
	assert.Equal(t, "ok", response.response.OK)
	assert.Zero(t, response.status)
	assert.Zero(t, response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSendOK(t *testing.T) {
	response := SendOK(OKResponse[string]{OK: "ok"})
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, 200, response.status)
	assert.Zero(t, response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSendOKJSON(t *testing.T) {
	response := SendOKJSON(OKResponse[string]{OK: "ok"})
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, 200, response.status)
	assert.Equal(t, "application/json", response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSendOKBytes(t *testing.T) {
	response := SendOKBytes(OKResponse[string]{OK: "ok"})
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, 200, response.status)
	assert.Equal(t, "application/octet-stream", response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSendOKText(t *testing.T) {
	response := SendOKText(OKResponse[string]{OK: "ok"})
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, 200, response.status)
	assert.Equal(t, "text/plain", response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSimpleSendWithContentType(t *testing.T) {
	response := Send(OKResponse[string]{OK: "ok"}).ContentType("example")
	assert.Equal(t, "ok", response.response.OK)
	assert.Zero(t, response.status)
	assert.Equal(t, "example", response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSimpleSendWithStatus(t *testing.T) {
	response := Send(OKResponse[string]{OK: "ok"}).Status(200)
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, 200, response.status)
	assert.Zero(t, response.contentType)
	assert.Equal(t, http.Header{}, response.headers)

	response2 := response.Status(400)
	// non mutating
	assert.Equal(t, 200, response.status)
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, 400, response2.status)
	assert.Zero(t, response.contentType)
	assert.Equal(t, http.Header{}, response.headers)
}

func TestSimpleSendAddHeaders(t *testing.T) {
	response := Send(OKResponse[string]{OK: "ok"}).
		AddHeader("x-foo", "bar").
		AddHeader("x-foo", "baz")
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, http.Header{"X-Foo": {"bar", "baz"}}, response.headers)
	assert.Zero(t, response.status)
	assert.Zero(t, response.contentType)
}

func TestSimpleSendSetHeader(t *testing.T) {
	response := Send(OKResponse[string]{OK: "ok"}).
		AddHeader("x-foo", "bar").
		SetHeader("x-foo", "baz")
	assert.Equal(t, "ok", response.response.OK)
	assert.Equal(t, http.Header{"X-Foo": {"baz"}}, response.headers)
	assert.Zero(t, response.status)
	assert.Zero(t, response.contentType)
}

func TestError(t *testing.T) {
	response, err := Error[OKResponse[string]](nil)
	assert.Zero(t, response)
	assert.Nil(t, err)
}

func TestRawHandler(t *testing.T) {
	rawHandler := RawHandler(func(c *Context) error {
		assert.Zero(t, *c.RawResponse)
		response, err := c.Next()
		assert.NoError(t, err)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "text/plain", response.ContentType)
		assert.Equal(t, []byte("test"), response.Body)
		return nil
	})
	assert.Equal(t, utils.NilType, rawHandler.requestTypes().requestBody)
	assert.Equal(t, utils.NilType, rawHandler.requestTypes().pathParams)
	assert.Equal(t, utils.NilType, rawHandler.requestTypes().queryParams)
	assert.Len(t, rawHandler.responseTypes(), 0)
	rawResponse := RawResponse{
		Status:      200,
		ContentType: "text/plain",
		Body:        []byte("test"),
		Headers:     nil,
	}
	handlerFunc := rawHandler.handlerFactory(openapi{}, func(c *Context) (RawResponse, error) {
		return rawResponse, nil
	})
	resp, err := handlerFunc(&Context{Request: &http.Request{}, RawResponse: &RawResponse{}})
	require.ErrorIs(t, err, UnsupportedResponseStatusErr)
	assert.Zero(t, resp)

}
