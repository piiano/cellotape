package router

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
	"github.com/piiano/cellotape/router/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestSchema(t *testing.T) {

	schema := jsonschema.Reflect(&Options{})
	bytes, _ := schema.MarshalJSON()
	log.Println(string(bytes))

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
