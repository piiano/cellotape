package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestDefaultRecoverFromError(t *testing.T) {
	writer := httptest.ResponseRecorder{}
	func() {
		defer defaultRecoverBehaviour(&writer)
		panic("unexpected error")
	}()
	assert.Equal(t, 500, writer.Code)
}

func TestName(t *testing.T) {
	writer := httptest.ResponseRecorder{}
	handlerFunc := asHttpRouterHandler(openapi{}, SpecOperation{}, func(_ *Context) (RawResponse, error) {
		return RawResponse{}, nil
	})
	handlerFunc(&writer, &http.Request{}, httprouter.Params{})

	assert.Equal(t, 500, writer.Code)
}

func TestDefaultOptionsHandler(t *testing.T) {
	writer := httptest.ResponseRecorder{}
	writer.Header().Add("Allow", "GET, POST")
	DefaultOptionsHandler(&writer, &http.Request{})

	assert.Equal(t, 204, writer.Code)
	assert.Equal(t, "GET, POST", writer.Header().Get("Allow"))
	assert.ElementsMatch(t, []string{"Allow"}, utils.Keys(writer.Header()))
}

func TestFailStartOnValidationError(t *testing.T) {
	_, err := createMainRouterHandler(&openapi{
		spec:    NewSpec(),
		options: DefaultOptions(),
		group: group{
			operations: []operation{
				{
					id: "test",
					handler: handler{
						request: requestTypes{
							requestBody: utils.NilType,
							pathParams:  utils.NilType,
							queryParams: utils.NilType,
						},
					},
				},
			},
		},
	})

	require.Error(t, err)
}

func TestOptionsHandler(t *testing.T) {
	router := httprouter.New()
	router.HandleMethodNotAllowed = false
	options := DefaultOptions()

	setGlobalHandlers(router, &openapi{
		spec: OpenAPISpec{
			Paths: map[string]*openapi3.PathItem{
				"/foo/{param1}/bar/{param2}": {
					Get:  &openapi3.Operation{},
					Post: &openapi3.Operation{},
				},
				"/foo/{param1}/bar": {
					Delete: &openapi3.Operation{},
					Patch:  &openapi3.Operation{},
				},
				"/foo/{param1}": {
					Get:     &openapi3.Operation{},
					Options: &openapi3.Operation{},
				},
			},
		},
		options: options,
	})

	router.GET("/foo/:param1/bar/:param2", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.POST("/foo/:param1/bar/:param2", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.DELETE("/foo/:param1/bar", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.PATCH("/foo/:param1/bar", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.GET("/foo/:param1", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.OPTIONS("/foo/:param1", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.WriteHeader(200)
	})

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/foo/1/bar/2", nil))
	response := writer.Result()

	require.Equal(t, 204, response.StatusCode)
	require.ElementsMatch(t, []string{"Allow"}, utils.Keys(response.Header))
	require.ElementsMatch(t, []string{"GET", "POST", "OPTIONS"}, strings.Split(response.Header.Get("Allow"), ", "))

	writer = httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/foo/1/bar", nil))
	response = writer.Result()

	require.Equal(t, 204, response.StatusCode)
	require.ElementsMatch(t, []string{"Allow"}, utils.Keys(response.Header))
	require.ElementsMatch(t, []string{"DELETE", "PATCH", "OPTIONS"}, strings.Split(response.Header.Get("Allow"), ", "))

	writer = httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/foo/1", nil))
	response = writer.Result()

	// Get to the registered handler and not the default one
	require.Equal(t, 200, response.StatusCode)

	writer = httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("GET", "/non-existing-path", nil))
	response = writer.Result()

	require.Equal(t, 404, response.StatusCode)
	require.ElementsMatch(t, []string{}, utils.Keys(response.Header))

	writer = httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/non-existing-path", nil))
	response = writer.Result()

	require.Equal(t, 404, response.StatusCode)
	require.ElementsMatch(t, []string{}, utils.Keys(response.Header))
}

func TestOptionsHandlerIsNil(t *testing.T) {
	router := httprouter.New()
	router.HandleMethodNotAllowed = false
	options := DefaultOptions()
	options.OptionsHandler = nil

	setGlobalHandlers(router, &openapi{
		spec: OpenAPISpec{
			Paths: map[string]*openapi3.PathItem{
				"/foo/{param1}/bar": {
					Get:  &openapi3.Operation{},
					Post: &openapi3.Operation{},
				},
				"/foo/{param1}": {
					Get:     &openapi3.Operation{},
					Options: &openapi3.Operation{},
				},
			},
		},
		options: options,
	})

	router.GET("/foo/:param1/bar", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.POST("/foo/:param1/bar", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.GET("/foo/:param1", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {})
	router.OPTIONS("/foo/:param1", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.WriteHeader(200)
	})

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/foo/1/bar", nil))
	response := writer.Result()

	require.Equal(t, 404, response.StatusCode)
	// No Allow header
	require.ElementsMatch(t, []string{}, utils.Keys(response.Header))

	writer = httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/foo/1", nil))
	response = writer.Result()

	// Get to the registered handler and not the default one
	require.Equal(t, 200, response.StatusCode)

	writer = httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest("OPTIONS", "/foo", nil))
	response = writer.Result()

	require.Equal(t, 404, response.StatusCode)
	require.ElementsMatch(t, []string{}, utils.Keys(response.Header))
}
