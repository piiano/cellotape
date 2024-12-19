package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/examples/todo_list_app_example/middlewares"
	"github.com/piiano/cellotape/examples/todo_list_app_example/rest"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	"github.com/piiano/cellotape/router"
)

func TestGetAllTasks(t *testing.T) {
	ts := initAPI(t)
	defer ts.Close()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks", ts.URL), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"results": [],
		"page": 0,
		"pageSize": 10,
		"isLast": true
	}`, string(response))
}

func TestCreateNewTaskAndGetIt(t *testing.T) {
	ts := initAPI(t)
	defer ts.Close()
	taskJson := `{
		"summary": "code first approach",
		"description": "add support for code first approach",
		"status": "open"
	}`
	request := bytes.NewBufferString(taskJson)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tasks", ts.URL), request)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	response := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Len(t, response, 1, "expecting one id field in the output")
	id, found := response["id"]
	assert.True(t, found)
	assert.Regexp(t, `[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}`, id)

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/tasks/%s", ts.URL, id), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	resp, err = client.Do(req)
	require.NoError(t, err)

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, taskJson, string(data))
}

func TestRequestQueryParamViolateSchemaValidations(t *testing.T) {
	ts := initAPI(t)
	defer ts.Close()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks?pageSize=30", ts.URL), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t,
		`invalid request query param. parameter "pageSize" in query has an error: Error at "/": number must be at most 20`,
		string(response))
}

func TestRequestPathParamViolateSchemaValidations(t *testing.T) {
	ts := initAPI(t)
	defer ts.Close()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks/123", ts.URL), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t,
		`invalid request path param. parameter "id" in path has an error: Error at "/": minimum string length is 36`,
		string(response))
}

func TestRequestBodyViolateSchemaValidations(t *testing.T) {
	ts := initAPI(t)
	defer ts.Close()
	taskJson := `{
		"summary": "code first approach",
		"description": "add support for code first approach",
		"status": "archived"
	}`
	request := bytes.NewBufferString(taskJson)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tasks", ts.URL), request)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t,
		`invalid request body. request body has an error: doesn't match schema #/components/schemas/Task: Error at "/status": value is not one of the allowed values ["open","in progress","closed"]`,
		string(response))
}

func initAPI(t *testing.T) *httptest.Server {
	spec, err := router.NewSpecFromData(specData)
	require.NoError(t, err)

	tasksService := services.NewTasksService()
	handler, err := router.NewOpenAPIRouter(spec).
		Use(middlewares.LoggerMiddleware, middlewares.AuthMiddleware).
		WithGroup(rest.TasksOperationsGroup(tasksService)).
		AsHandler()
	require.NoError(t, err)

	return httptest.NewServer(handler)
}
