package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router"
	"reflect"
	"testing"
)

// according to the spec the object validation properties should apply only when the type is set to object
func TestObjectSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchemaWithProperty := openapi3.NewSchema().WithProperty("name", openapi3.NewStringSchema())
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *untypedSchemaWithProperty, router.Options{})
	for _, validType := range types {
		t.Run(validType.String(), func(t *testing.T) {
			if err := validator.WithType(validType).validateObjectSchema(); err != nil {
				t.Errorf("expect untyped schema to be compatible with %s type", validType)
			}
		})
	}
}

func TestObjectSchemaValidatorWithSimpleStruct(t *testing.T) {
	type SimpleStruct struct {
		Field1 string
		Field2 struct {
			Field2_A string `json:"renamed_field_2_a,omitempty"`
			Field2_B []bool
		}
		Field3 int
	}
	simpleStructSchema := openapi3.NewObjectSchema().
		WithProperty("Field1", openapi3.NewStringSchema()).
		WithProperty("Field2", openapi3.NewObjectSchema().
			WithProperty("renamed_field_2_a", openapi3.NewStringSchema()).
			WithProperty("Field2_B", openapi3.NewArraySchema().WithItems(openapi3.NewBoolSchema()))).
		WithProperty("Field3", openapi3.NewIntegerSchema())
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *simpleStructSchema, router.Options{})
	simpleStructType := reflect.TypeOf(SimpleStruct{})
	errTemplate := "expect object schema with to be %s with %s type"
	expectTypeToBeCompatible(t, validator, simpleStructType, errTemplate, "compatible", simpleStructType)
	for _, invalidType := range types {
		t.Run(invalidType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, invalidType, errTemplate, "incompatible", invalidType)
		})
	}
}
