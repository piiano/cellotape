package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"reflect"
	"testing"
)

// according to the spec the array validation properties should apply oly when the type is set to array
func TestArraySchemaValidatorWithUntypedSchema(t *testing.T) {
	// create with NewSchema and not with NewArraySchema for an untyped schema
	untypedSchemaWithItemsProperty := openapi3.NewSchema().WithItems(openapi3.NewStringSchema())
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *untypedSchemaWithItemsProperty, Options{})
	for _, validType := range types {
		t.Run(validType.String(), func(t *testing.T) {
			if err := validator.WithType(validType).validateArraySchema(); err != nil {
				t.Errorf("expect untyped schema to be compatible with %s type", validType)
			}
		})
	}
}

func TestArraySchemaValidatorPassForSimpleArray(t *testing.T) {
	stringArraySchema := openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *stringArraySchema, Options{})
	stringArraySchema.Title = "StringArray"
	stringArrayType := reflect.TypeOf([1]string{""})
	if err := validator.WithType(stringArrayType).Validate(); err != nil {
		t.Errorf("expect string array schema to be compatible with %s type", stringArrayType.String())
		t.Error(err)
	}
	stringSliceType := reflect.TypeOf(make([]string, 1))
	if err := validator.WithType(stringSliceType).Validate(); err != nil {
		t.Errorf("expect string array schema to be compatible with %s type", stringSliceType.String())
		t.Error(err)
	}
}

func TestArraySchemaValidatorFailOnWrongType(t *testing.T) {
	stringArraySchema := openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *stringArraySchema, Options{})
	stringArraySchema.Title = "StringArray"
	intArrayType := reflect.TypeOf([1]int{1})
	if err := validator.WithType(intArrayType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", intArrayType.String())
	}
	intSliceType := reflect.TypeOf(make([]any, 1))
	if err := validator.WithType(intSliceType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", intSliceType.String())
	}
	stringType := reflect.TypeOf("")
	if err := validator.WithType(stringType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", stringType.String())
	}
	objectType := reflect.TypeOf(struct{}{})
	if err := validator.WithType(objectType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", objectType.String())
	}
	mapIntToStringType := reflect.TypeOf(make(map[int]string))
	if err := validator.WithType(mapIntToStringType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", mapIntToStringType.String())
	}
}
