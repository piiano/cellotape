package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
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
