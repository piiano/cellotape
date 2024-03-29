package router

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestSchema(t *testing.T) {
	schema := jsonschema.Reflect(&Options{})
	bytes, _ := json.MarshalIndent(schema, "", "  ")

	fmt.Println(string(bytes))
	schemaFile, err := os.ReadFile("../options-schema.json")
	assert.NoError(t, err)

	assert.Equal(t, string(schemaFile), string(bytes))
}

func TestBehaviourZeroValue(t *testing.T) {
	var zeroValueBehaviour Behaviour
	assert.Equal(t, PropagateError, zeroValueBehaviour)
	assert.NotEqual(t, PrintWarning, zeroValueBehaviour)
	assert.NotEqual(t, Ignore, zeroValueBehaviour)
}

func TestBehaviourMarshalText(t *testing.T) {
	behaviours := []Behaviour{PropagateError, PrintWarning, Ignore}
	jsonBytes, err := json.Marshal(behaviours)
	require.NoError(t, err)
	assert.Equal(t, `["propagate-error","print-warning","ignore"]`, string(jsonBytes))

	_, err = Behaviour(utils.Info).MarshalText()
	require.Error(t, err)
}

func TestBehaviourUnmarshalText(t *testing.T) {
	var behaviours []Behaviour
	err := json.Unmarshal([]byte(`["propagate-error","print-warning","ignore"]`), &behaviours)
	require.NoError(t, err)

	assert.Equal(t, []Behaviour{PropagateError, PrintWarning, Ignore}, behaviours)

	var behaviour Behaviour
	err = json.Unmarshal([]byte(`"invalid-behavior"`), &behaviour)
	require.Error(t, err)
}

func TestOperationValidationOptions(t *testing.T) {
	options := DefaultOptions()

	operationOptions := options.operationValidationOptions("foo")
	assert.Equal(t, options.DefaultOperationValidation, operationOptions)

	customOperationOptions := OperationValidationOptions{
		ValidateRequestBody: Ignore,
	}
	options.OperationValidations = map[string]OperationValidationOptions{
		"foo": customOperationOptions,
	}
	operationOptions = options.operationValidationOptions("foo")
	assert.Equal(t, customOperationOptions, operationOptions)
}
