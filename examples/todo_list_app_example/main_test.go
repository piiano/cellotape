package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/middlewares"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/rest"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	"github.com/piiano/restcontroller/router"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTasks(t *testing.T) {
	ts, err := initAPI(t)
	defer ts.Close()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks", ts.URL), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	client := http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.JSONEq(t, `{
		"results": [],
		"page": 0,
		"pageSize": 0,
		"isLast": true
	}`, string(response))
}

func TestCreateNewTaskAndGetIt(t *testing.T) {
	ts, err := initAPI(t)
	defer ts.Close()
	taskJson := `{
		"summary": "code first approach",
		"description": "add support for code first approach",
		"status": "open"
	}`
	request := bytes.NewBufferString(taskJson)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tasks", ts.URL), request)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	//req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	response := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.Nil(t, err)
	assert.Len(t, response, 1, "expecting one id field in the output")
	id, found := response["id"]
	assert.True(t, found)
	assert.Regexp(t, `[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}`, id)

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/tasks/%s", ts.URL, id), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", "Bearer secret")
	resp, err = client.Do(req)
	require.Nil(t, err)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.JSONEq(t, taskJson, string(data))
}

func initAPI(t *testing.T) (*httptest.Server, error) {
	spec, err := router.NewSpecFromData(specData)
	require.Nil(t, err)

	tasksService := services.NewTasksService()
	handler, err := router.NewOpenAPIRouter(spec).
		Use(middlewares.LoggerMiddleware, middlewares.AuthMiddleware).
		WithGroup(rest.TasksOperationsGroup(tasksService)).
		AsHandler()
	require.Nil(t, err)

	ts := httptest.NewServer(handler)
	return ts, err
}
