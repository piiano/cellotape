package schema_validator

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

// according to the spec the array validation properties should apply oly when the type is set to array
func TestArraySchemaValidatorWithUntypedSchema(t *testing.T) {
	// create with NewSchema and not with NewArraySchema for an untyped schema
	untypedSchemaWithItemsProperty := openapi3.NewSchema().WithItems(openapi3.NewStringSchema())

	for _, goType := range types {
		valid := (kindIs(reflect.Array, reflect.Slice)(goType) && isSerializedFromString(goType.Elem())) ||
			!kindIs(reflect.Array, reflect.Slice)(goType) ||
			(uuidType.ConvertibleTo(goType))

		t.Run(goType.String(), func(t *testing.T) {
			err := schemaValidator(*untypedSchemaWithItemsProperty).WithType(goType).Validate()
			if valid {
				require.NoErrorf(t, err, "expect untyped schema to be compatible with %s type", goType)
			} else {
				require.Errorf(t, err, "expect untyped schema to be incompatible with %s type", goType)
			}
		})
	}
}

func TestArraySchemaValidatorPassForSimpleArray(t *testing.T) {
	stringArraySchema := openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	validator := schemaValidator(*stringArraySchema)
	stringArrayType := reflect.TypeOf([1]string{""})
	if err := validator.WithType(stringArrayType).Validate(); err != nil {
		t.Errorf("expect string array schema to be compatible with %s type", stringArrayType)
		t.Error(err)
	}
	stringSliceType := reflect.TypeOf(make([]string, 1))
	if err := validator.WithType(stringSliceType).Validate(); err != nil {
		t.Errorf("expect string array schema to be compatible with %s type", stringSliceType)
		t.Error(err)
	}
}

func TestArraySchemaValidatorFailOnWrongType(t *testing.T) {
	stringArraySchema := openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	validator := schemaValidator(*stringArraySchema)
	intArrayType := reflect.TypeOf([1]int{1})
	if err := validator.WithType(intArrayType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", intArrayType)
	}
	if err := validator.WithType(stringType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", stringType)
	}
	if err := validator.WithType(emptyStructTypes).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", emptyStructTypes)
	}
	mapIntToStringType := reflect.TypeOf(make(map[int]string))
	if err := validator.WithType(mapIntToStringType).Validate(); err == nil {
		t.Errorf("expect string array schema to be incompatible with %s type", mapIntToStringType)
	}
}
