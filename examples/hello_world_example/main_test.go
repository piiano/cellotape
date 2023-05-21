package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/examples/hello_world_example/api"
	"github.com/piiano/cellotape/router"
)

func TestHelloWorldExample(t *testing.T) {
	spec, err := router.NewSpecFromData(specData)
	require.NoError(t, err)

	handler, err := router.NewOpenAPIRouter(spec).
		WithOperation("greet", api.GreetOperationHandler).
		AsHandler()
	require.NoError(t, err)

	ts := httptest.NewServer(handler)
	defer ts.Close()
	request := bytes.NewBufferString(`{ "name": "Ori" }`)
	resp, err := http.Post(fmt.Sprintf("%s/v1/greet", ts.URL), "application/json", request)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	response := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"greeting": "Hello Ori!"}, response)
}
