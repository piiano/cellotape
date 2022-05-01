package main

import (
	"bytes"
	"fmt"
	models "github.com/piiano/restcontroller/examples/todo_list_app_example/rest"
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
	err, ts := initAPI(t)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/tasks", ts.URL))
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

func TestCreateNewTasks(t *testing.T) {
	err, ts := initAPI(t)
	defer ts.Close()
	request := bytes.NewBufferString(`{
		"summary": "code first approach",
		"description": "add support for code first approach",
		"status": "open"
	}`)
	resp, err := http.Post(fmt.Sprintf("%s/tasks", ts.URL), "application/json", request)
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.Regexp(t, `\{"id":"[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}"}`, string(response))
}

func initAPI(t *testing.T) (error, *httptest.Server) {
	spec, err := router.NewSpecFromData(specData)
	require.Nil(t, err)

	tasks := services.NewTasksService()
	handler, err := router.NewOpenAPIRouter(spec).
		WithGroup(models.TasksOperationsGroup(tasks)).
		AsHandler()
	require.Nil(t, err)

	ts := httptest.NewServer(handler)
	return err, ts
}
