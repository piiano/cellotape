package router

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

//go:embed test_specs/openapi.yaml
var specData []byte

func TestNewSpec(t *testing.T) {
	spec := NewSpec()
	assert.NotNil(t, spec)
}

func TestNewSpecFromData(t *testing.T) {
	spec, err := NewSpecFromData(specData)
	require.NoError(t, err)
	assert.NotNil(t, spec)

	_, err = NewSpecFromData(bytes.NewBufferString("}").Bytes())
	require.NotNil(t, err)
}

func TestNewSpecFromFile(t *testing.T) {
	spec, err := NewSpecFromFile("test_specs/openapi.yaml")
	require.NoError(t, err)
	assert.NotNil(t, spec)

	_, err = NewSpecFromFile("test_specs/no_such_file.yaml")
	require.NotNil(t, err)
}

func TestFindSpecOperationByIDFail(t *testing.T) {
	spec, err := NewSpecFromData(bytes.NewBufferString("{}").Bytes())
	require.NoError(t, err)
	_, ok := spec.findSpecOperationByID("no-such-id")
	assert.False(t, ok)
}

func TestFindSpecOperationByIDPass(t *testing.T) {
	spec, err := NewSpecFromData(specData)
	require.NoError(t, err)
	op, ok := spec.findSpecOperationByID("greet")
	assert.True(t, ok)
	require.NotNil(t, op)
	require.Equal(t, "greet", op.OperationID)
	require.Equal(t, "/{version}/greet", op.Path)
	require.Equal(t, "POST", op.Method)
}

func TestFindSpecContentTypes(t *testing.T) {
	spec := OpenAPISpec(openapi3.T{
		Paths: openapi3.NewPaths(
			openapi3.WithPath("/1", &openapi3.PathItem{
				Get: &openapi3.Operation{
					RequestBody: &openapi3.RequestBodyRef{
						Value: openapi3.NewRequestBody().WithSchema(openapi3.NewSchema(), []string{"text/plain"}),
					},
				},
			}),
			openapi3.WithPath("/2", &openapi3.PathItem{
				Get: &openapi3.Operation{
					Responses: openapi3.NewResponses(
						openapi3.WithStatus(200, &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithJSONSchema(openapi3.NewSchema()),
						}),
						openapi3.WithStatus(500, &openapi3.ResponseRef{}),
					),
				},
			}),
		),
	})

	contentTypes := spec.findSpecContentTypes(utils.NewSet[string]())
	require.ElementsMatch(t, []string{"application/json", "text/plain"}, contentTypes)
}
