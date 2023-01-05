package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
