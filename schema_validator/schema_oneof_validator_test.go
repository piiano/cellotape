package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"testing"
)

func TestSchemaOneOfValidatorPass(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	notBooleanSchema.OneOf = openapi3.SchemaRefs{
		openapi3.NewBoolSchema().NewRef(),
		openapi3.NewStringSchema().NewRef(),
		openapi3.NewInt64Schema().NewRef(),
	}
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *notBooleanSchema, router.Options{})
	var validTypes = []reflect.Type{boolType, stringType, int64Type}
	errTemplate := "expect schema with oneOf property to be compatible with %s type"
	for _, validType := range validTypes {
		t.Run(validType.String(), func(t *testing.T) {
			expectTypeToBeCompatible(t, validator, validType, errTemplate, validType)
		})
	}
}

func TestSchemaOneOfValidatorFailOnMoreThanOneMatchedType(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	numberSchema := openapi3.NewSchema()
	numberSchema.Type = string(numberSchemaType)
	notBooleanSchema.OneOf = openapi3.SchemaRefs{
		openapi3.NewBoolSchema().NewRef(),
		openapi3.NewStringSchema().NewRef(),
		openapi3.NewInt64Schema().NewRef(),
		numberSchema.NewRef(),
	}
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *notBooleanSchema, router.Options{})
	errTemplate := "expect schema with oneOf property to be incompatible with %s type"
	expectTypeToBeIncompatible(t, validator, int64Type, errTemplate, int64Type)
}

func TestSchemaOneOfValidatorFailOnNoMatchedType(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	notBooleanSchema.OneOf = openapi3.SchemaRefs{
		openapi3.NewBoolSchema().NewRef(),
		openapi3.NewStringSchema().NewRef(),
		openapi3.NewInt64Schema().NewRef(),
	}
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *notBooleanSchema, router.Options{})
	invalidTypes := utils.Filter(types, func(t reflect.Type) bool {
		return t != boolType && t != stringType && t != int64Type
	})
	errTemplate := "expect schema with oneOf property to be incompatible with %s type"
	for _, invalidType := range invalidTypes {
		t.Run(invalidType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, invalidType, errTemplate, invalidType)
		})
	}
}
