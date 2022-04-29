package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"testing"
)

func TestBooleanSchemaValidatorPassForBoolType(t *testing.T) {
	booleanSchema := openapi3.NewBoolSchema()
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *booleanSchema, Options{})
	if err := validator.WithType(boolType).Validate(); err != nil {
		expectTypeToBeCompatible(t, validator, boolType, "expect boolean schema to be compatible with %s type", boolType)
	}
}

// according to the spec the boolean validation properties should apply only when the type is set to boolean
func TestBoolSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchema := openapi3.NewSchema()
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *untypedSchema, Options{})
	for _, validType := range types {
		t.Run(validType.String(), func(t *testing.T) {
			if err := validator.WithType(validType).validateBooleanSchema(); err != nil {
				t.Errorf("expect untyped schema to be compatible with %s type", validType)
			}
		})
	}
}

func TestBooleanSchemaValidatorFailOnWrongType(t *testing.T) {
	booleanSchema := openapi3.NewBoolSchema()
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *booleanSchema, Options{})
	errTemplate := "expect boolean schema to be incompatible with %s type"
	// omit the bool type from all defined test types
	var nonBoolTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != boolType
	})
	for _, nonBoolType := range nonBoolTypes {
		t.Run(nonBoolType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonBoolType, errTemplate, nonBoolType)
		})
	}
}
