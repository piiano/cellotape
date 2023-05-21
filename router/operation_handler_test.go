package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestHandlerFuncTypeExtraction(t *testing.T) {
	fn := HandlerFunc[utils.Nil, utils.Nil, utils.Nil, utils.Nil](func(*Context, Request[utils.Nil, utils.Nil, utils.Nil]) (Response[utils.Nil], error) {
		return Response[utils.Nil]{}, nil
	})
	types := fn.requestTypes()
	assert.Equal(t, types.requestBody, utils.NilType)
	assert.Equal(t, types.pathParams, utils.NilType)
	assert.Equal(t, types.queryParams, utils.NilType)
}

func TestRouterAsHandler(t *testing.T) {
	type responses struct {
		Answer int `status:"200"`
	}
	fn := HandlerFunc[utils.Nil, utils.Nil, utils.Nil, responses](func(*Context, Request[utils.Nil, utils.Nil, utils.Nil]) (Response[responses], error) {
		return SendOKJSON(responses{Answer: 42}), nil
	})
	spec, err := NewSpecFromData([]byte(`
  { "openapi": "3.0.3", "info": { "title": "test", "version": "1.0.0" }, "paths": { "/abc": { "post": {
    "operationId": "id",
    "summary": "the ultimate answer to life the universe and everything",
    "responses":{ "200": { "description": "ok", "content": { "application/json": { "schema": { "type": "integer" } } } } }
  } } } }`))
	require.NoError(t, err)
	router := NewOpenAPIRouter(spec).WithOperation("id", fn)
	h, err := router.AsHandler()
	require.NoError(t, err)
	ts := httptest.NewServer(h)
	defer ts.Close()

	// test valid request
	resp, err := http.Post(fmt.Sprintf("%s/abc", ts.URL), "application/json", nil)
	require.NoError(t, err)
	res, err := toString(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "42", res)
	assert.Equal(t, 200, resp.StatusCode)

	// TODO: Add support for runtime validation based on spec and restore this test
	// test bad request
	//resp, err = http.Post(fmt.Sprintf("%s/abc", ts.URL), "application/json", bytes.NewBufferString("{}"))
	//require.NoError(t, err)
	//assert.Equal(t, 400, resp.StatusCode)
	//res, err = toString(resp.Body)
	//require.NoError(t, err)
	//assert.Equal(t, `{"error":"expected request with no body payload"}`, res)
}

func toString(reader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
