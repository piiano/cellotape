package schema_validator

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/cellotape/router/utils"
	"reflect"
	"testing"
)

func TestNumberSchemaValidatorPassForIntType(t *testing.T) {
	numberSchema := openapi3.NewSchema()
	numberSchema.Type = numberSchemaType
	validator := schemaValidator(*numberSchema)
	errTemplate := "expect number schema to be compatible with %s type"
	for _, numericType := range numericTypes {
		t.Run(fmt.Sprintf("test expected pass with %s type", numericType), func(t *testing.T) {
			expectTypeToBeCompatible(t, validator, numericType, errTemplate, numericType)
		})
	}
}

// according to the spec the number validation properties should apply only when the type is set to number
func TestNumberSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchemaWithDoubleFormat := openapi3.NewSchema().WithFormat(doubleFormat)
	validator := schemaValidator(*untypedSchemaWithDoubleFormat)
	for _, validType := range types {
		t.Run(validType.String(), func(t *testing.T) {
			if err := validator.WithType(validType).validateNumberSchema(); err != nil {
				t.Errorf("expect untyped schema to be compatible with %s type", validType)
			}
		})
	}
}

func TestNumberSchemaValidatorFailOnWrongType(t *testing.T) {
	numberSchema := openapi3.NewSchema()
	numberSchema.Type = numberSchemaType
	validator := schemaValidator(*numberSchema)
	errTemplate := "expect number schema to be incompatible with %s type"
	// filter all numeric types from all defined test types
	var nonNumericTypes = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		_, found := utils.Find[reflect.Type](numericTypes, func(numericType reflect.Type) bool {
			return t == numericType || t == reflect.PointerTo(numericType)
		})
		return !found
	})
	for _, nonNumericType := range nonNumericTypes {
		t.Run(fmt.Sprintf("test expected fail with %s type", nonNumericType), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonNumericType, errTemplate, nonNumericType)
		})
	}
}

func TestFloatFormatSchemaValidatorPassForFloat32Type(t *testing.T) {
	floatSchema := openapi3.NewSchema()
	floatSchema.Type = numberSchemaType
	floatSchema.Format = floatFormat
	validator := schemaValidator(*floatSchema)
	errTemplate := "expect number schema with float format to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, float32Type, errTemplate, float32Type)
}

func TestFloat32FormatSchemaValidatorFailOnWrongType(t *testing.T) {
	floatSchema := openapi3.NewSchema()
	floatSchema.Type = numberSchemaType
	floatSchema.Format = floatFormat
	validator := schemaValidator(*floatSchema)
	errTemplate := "expect number schema with float format to be incompatible with %s type"
	// omit the float32 type from all defined test types
	var nonFloat32Types = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		return t != float32Type && t != reflect.PointerTo(float32Type)
	})
	for _, nonFloat32Type := range nonFloat32Types {
		t.Run(nonFloat32Type.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonFloat32Type, errTemplate, nonFloat32Type)
		})
	}
}

func TestDoubleFormatSchemaValidatorPassForFloat32Type(t *testing.T) {
	doubleSchema := openapi3.NewSchema()
	doubleSchema.Type = numberSchemaType
	doubleSchema.Format = doubleFormat
	validator := schemaValidator(*doubleSchema)
	errTemplate := "expect number schema with double format to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, float64Type, errTemplate, float64Type)
}

func TestDoubleFormatSchemaValidatorFailOnWrongType(t *testing.T) {
	doubleSchema := openapi3.NewSchema()
	doubleSchema.Type = numberSchemaType
	doubleSchema.Format = doubleFormat
	validator := schemaValidator(*doubleSchema)
	errTemplate := "expect number schema with double format to be incompatible with %s type"
	// omit the float64 type from all defined test types
	var nonFloat64Types = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		return t != float64Type && t != reflect.PointerTo(float64Type)
	})
	for _, nonFloat64Type := range nonFloat64Types {
		t.Run(nonFloat64Type.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonFloat64Type, errTemplate, nonFloat64Type)
		})
	}
}
