package schema_validator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/utils"
)

func TestIntegerSchemaValidatorPassForIntType(t *testing.T) {
	integerSchema := openapi3.NewIntegerSchema()
	validator := schemaValidator(*integerSchema)
	errTemplate := "expect integer schema to be compatible with %s type"
	for _, intType := range intTypes {
		t.Run(fmt.Sprintf("test expected pass with %s type", intType), func(t *testing.T) {
			expectTypeToBeCompatible(t, validator, intType, errTemplate, intType)
		})
	}
}

// according to the spec the integer validation properties should apply only when the type is set to integer
func TestIntegerSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchemaWithInt64Format := openapi3.NewSchema().WithFormat(int64Format)
	validator := schemaValidator(*untypedSchemaWithInt64Format)
	for _, validType := range types {
		t.Run(validType.String(), func(t *testing.T) {
			if err := validator.WithType(validType).validateIntegerSchema(); err != nil {
				t.Errorf("expect untyped schema to be compatible with %s type", validType)
			}
		})
	}
}

func TestIntegerSchemaValidatorFailOnWrongType(t *testing.T) {
	integerSchema := openapi3.NewIntegerSchema()
	validator := schemaValidator(*integerSchema)
	errTemplate := "expect integer schema to be incompatible with %s type"
	// filter all int types from all defined test types
	var nonIntTypes = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		_, found := utils.Find[reflect.Type](intTypes, func(intType reflect.Type) bool {
			return t == intType || t == reflect.PointerTo(intType)
		})
		return !found
	})
	for _, nonIntType := range nonIntTypes {
		t.Run(fmt.Sprintf("test expected fail with %s type", nonIntType), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonIntType, errTemplate, nonIntType)
		})
	}
}

func TestInt64FormatSchemaValidatorPassForInt64Type(t *testing.T) {
	int64Schema := openapi3.NewInt64Schema()
	validator := schemaValidator(*int64Schema)
	errTemplate := "expect integer schema with int64 format to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, int64Type, errTemplate, int64Type)
}

func TestInt64FormatSchemaValidatorFailOnWrongType(t *testing.T) {
	int64Schema := openapi3.NewInt64Schema()
	validator := schemaValidator(*int64Schema)
	errTemplate := "expect integer schema with int64 format schema to be incompatible with %s type"
	// omit the int64 type from all defined test types
	var nonInt64Types = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		return t != int64Type && t != reflect.PointerTo(int64Type)
	})
	for _, nonInt64Type := range nonInt64Types {
		t.Run(nonInt64Type.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonInt64Type, errTemplate, nonInt64Type)
		})
	}
}

func TestInt32FormatSchemaValidatorPassForInt32Type(t *testing.T) {
	int32Schema := openapi3.NewInt32Schema()
	validator := schemaValidator(*int32Schema)
	errTemplate := "expect integer schema with int32 format to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, int32Type, errTemplate, int32Type)
}

func TestInt32FormatSchemaValidatorFailOnWrongType(t *testing.T) {
	int32Schema := openapi3.NewInt32Schema()
	validator := schemaValidator(*int32Schema)
	errTemplate := "expect integer schema with int32 format schema to be incompatible with %s type"
	// omit the int32 type from all defined test types
	var nonInt32Types = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		return t != int32Type && t != reflect.PointerTo(int32Type)
	})
	for _, nonInt32Type := range nonInt32Types {
		t.Run(nonInt32Type.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonInt32Type, errTemplate, nonInt32Type)
		})
	}
}
