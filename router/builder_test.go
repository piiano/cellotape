package router

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestNewOpenAPIRouter(t *testing.T) {
	spec := NewSpec()
	oa := NewOpenAPIRouter(spec)
	assert.NotNil(t, oa)
}

func TestNewOpenAPIRouterWithOptions(t *testing.T) {
	spec := NewSpec()
	oa := NewOpenAPIRouterWithOptions(spec, Options{})
	assert.NotNil(t, oa)
}

type nilContentType struct{}

func (m nilContentType) Mime() string                 { return "nil" }
func (m nilContentType) Encode(_ any) ([]byte, error) { return nil, nil }
func (m nilContentType) Decode(_ []byte, _ any) error { return nil }
func (m nilContentType) ValidateTypeSchema(utils.Logger, utils.LogLevel, reflect.Type, openapi3.Schema) error {
	return nil
}

func TestOpenAPIRouterWithContentType(t *testing.T) {
	spec := NewSpec()
	router := NewOpenAPIRouter(spec)
	oa := router.(*openapi)

	require.Equal(t, DefaultContentTypes(), oa.contentTypes)

	router.WithContentType(nilContentType{})
	contentTypes := DefaultContentTypes()
	contentTypes["nil"] = nilContentType{}
	require.Equal(t, contentTypes, oa.contentTypes)

}

func TestNewOpenAPIRouterReturnSpec(t *testing.T) {
	spec := NewSpec()
	oa := NewOpenAPIRouter(spec)
	require.NotNil(t, oa)

	routerSpec := oa.Spec()
	assert.Equal(t, spec, routerSpec)
}
