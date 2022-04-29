package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"reflect"
	"testing"
)

// according to the spec the object validation properties should apply only when the type is set to object
func TestObjectSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchemaWithProperty := openapi3.NewSchema().WithProperty("name", openapi3.NewStringSchema())
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *untypedSchemaWithProperty, Options{})
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
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *simpleStructSchema, Options{})
	simpleStructType := reflect.TypeOf(SimpleStruct{})
	errTemplate := "expect object schema to be %s with %s type"
	expectTypeToBeCompatible(t, validator, simpleStructType, errTemplate, "compatible", simpleStructType)
	for _, invalidType := range types {
		t.Run(invalidType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, invalidType, errTemplate, "incompatible", invalidType)
		})
	}
}

func TestObjectSchemaValidatorWithEmbeddedStruct(t *testing.T) {
	type SimpleA struct{ Field1 bool }
	type SimpleB struct {
		SimpleA
		Field2 int
	}
	type SimpleC struct {
		SimpleA
		Field1 string
		Field2 int
	}
	structBType := reflect.TypeOf(SimpleB{})
	structCType := reflect.TypeOf(SimpleC{})
	structBSchema := *openapi3.NewObjectSchema().
		WithProperty("Field1", openapi3.NewBoolSchema()).
		WithProperty("Field2", openapi3.NewIntegerSchema())
	structCSchema := *openapi3.NewObjectSchema().
		WithProperty("Field1", openapi3.NewStringSchema()).
		WithProperty("Field2", openapi3.NewIntegerSchema())
	validatorB := NewTypeSchemaValidator(structBType, structBSchema, Options{})
	validatorC := NewTypeSchemaValidator(structCType, structCSchema, Options{})
	errTemplate := "expect object schema %s to be %s with %s type"
	expectTypeToBeCompatible(t, validatorB, structBType, errTemplate, "structBSchema", "compatible", structBType)
	expectTypeToBeCompatible(t, validatorC, structCType, errTemplate, "structCSchema", "compatible", structCType)
	expectTypeToBeIncompatible(t, validatorB.WithSchema(structCSchema), structBType, errTemplate, "structCSchema", "incompatible", structBType)
	expectTypeToBeIncompatible(t, validatorC.WithSchema(structBSchema), structCType, errTemplate, "structBSchema", "incompatible", structCType)
}
