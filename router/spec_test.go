package router

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed test_specs/openapi.yaml
var specData []byte

func TestNewSpec(t *testing.T) {
	spec := NewSpec()
	assert.NotNil(t, spec)
}

func TestNewSpecFromData(t *testing.T) {
	spec, err := NewSpecFromData(specData)
	require.Nil(t, err)
	assert.NotNil(t, spec)

	_, err = NewSpecFromData(bytes.NewBufferString("}").Bytes())
	require.NotNil(t, err)
}

func TestNewSpecFromFile(t *testing.T) {
	spec, err := NewSpecFromFile("test_specs/openapi.yaml")
	require.Nil(t, err)
	assert.NotNil(t, spec)

	_, err = NewSpecFromFile("test_specs/no_such_file.yaml")
	require.NotNil(t, err)
}

func TestFindSpecOperationByIDFail(t *testing.T) {
	spec, err := NewSpecFromData(bytes.NewBufferString("{}").Bytes())
	require.Nil(t, err)
	_, ok := spec.findSpecOperationByID("no-such-id")
	assert.False(t, ok)
}

func TestFindSpecOperationByIDPass(t *testing.T) {
	spec, err := NewSpecFromData(specData)
	require.Nil(t, err)
	op, ok := spec.findSpecOperationByID("greet")
	assert.True(t, ok)
	require.NotNil(t, op)
	require.Equal(t, "greet", op.OperationID)
	require.Equal(t, "/{version}/greet", op.path)
	require.Equal(t, "POST", op.method)

}
