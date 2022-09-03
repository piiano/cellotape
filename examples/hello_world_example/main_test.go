package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/piiano/cellotape/examples/hello_world_example/api"
	"github.com/piiano/cellotape/router"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorldExample(t *testing.T) {
	spec, err := router.NewSpecFromData(specData)
	require.NoError(t, err)

	handler, err := router.NewOpenAPIRouter(spec).
		WithOperation("greet", api.GreetOperationHandler).
		AsHandler()
	fmt.Println(err)
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
