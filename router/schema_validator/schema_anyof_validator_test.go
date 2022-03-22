package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/cellotape/router/utils"
	"reflect"
	"testing"
)

func TestSchemaAnyOfValidatorPass(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	notBooleanSchema.AnyOf = openapi3.SchemaRefs{
		openapi3.NewBoolSchema().NewRef(),
		openapi3.NewStringSchema().NewRef(),
		openapi3.NewInt64Schema().NewRef(),
	}
	validator := schemaValidator(*notBooleanSchema)
	var validTypes = []reflect.Type{boolType, stringType, int64Type}
	errTemplate := "expect schema with anyOf property to be compatible with %s type"
	for _, validType := range validTypes {
		t.Run(validType.String(), func(t *testing.T) {
			expectTypeToBeCompatible(t, validator, validType, errTemplate, validType)
		})
	}
}

func TestSchemaAnyOfValidatorPassOnMoreThanOneMatchedType(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	numberSchema := openapi3.NewSchema()
	numberSchema.Type = numberSchemaType
	notBooleanSchema.AnyOf = openapi3.SchemaRefs{
		openapi3.NewBoolSchema().NewRef(),
		openapi3.NewStringSchema().NewRef(),
		openapi3.NewInt64Schema().NewRef(),
		numberSchema.NewRef(),
	}
	validator := schemaValidator(*notBooleanSchema)
	errTemplate := "expect schema with anyOf property to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, int64Type, errTemplate, int64Type)
}

func TestSchemaAnyOfValidatorFailOnNoMatchedType(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	notBooleanSchema.AnyOf = openapi3.SchemaRefs{
		openapi3.NewBoolSchema().NewRef(),
		openapi3.NewStringSchema().NewRef(),
		openapi3.NewInt64Schema().NewRef(),
	}
	validator := schemaValidator(*notBooleanSchema)
	invalidTypes := utils.Filter(types, func(t reflect.Type) bool {
		return t != boolType && t != stringType && t != int64Type
	})
	errTemplate := "expect schema with anyOf property to be incompatible with %s type"
	for _, invalidType := range invalidTypes {
		t.Run(invalidType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, invalidType, errTemplate, invalidType)
		})
	}
}
