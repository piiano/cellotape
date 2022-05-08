package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/utils"
	"reflect"
	"testing"
)

/* ===============================================================================================================
# Type Tests Infrastructure

This infrastructure allows testing all sort of types with a given schema.
It generates predefined set of types that can be iterated dynamically during tests to check different cases.
 =============================================================================================================== */

var (
	stringType       = reflect.TypeOf("")
	boolType         = reflect.TypeOf(false)
	intType          = reflect.TypeOf(0)
	int8Type         = reflect.TypeOf(int8(0))
	int16Type        = reflect.TypeOf(int16(0))
	int32Type        = reflect.TypeOf(int32(0))
	int64Type        = reflect.TypeOf(int64(0))
	uintType         = reflect.TypeOf(uint(0))
	uint8Type        = reflect.TypeOf(uint8(0))
	uint16Type       = reflect.TypeOf(uint16(0))
	uint32Type       = reflect.TypeOf(uint32(0))
	uint64Type       = reflect.TypeOf(uint64(0))
	float32Type      = reflect.TypeOf(float32(0))
	float64Type      = reflect.TypeOf(float64(0))
	signedIntTypes   = []reflect.Type{intType, int8Type, int16Type, int32Type, int64Type}
	unsignedIntTypes = []reflect.Type{uintType, uint8Type, uint16Type, uint32Type, uint64Type}
	intTypes         = append(signedIntTypes, unsignedIntTypes...)
	floatTypes       = []reflect.Type{float32Type, float64Type}
	numericTypes     = append(intTypes, floatTypes...)
	primitiveTypes   = append(numericTypes, stringType, boolType)
	emptyStructTypes = reflect.TypeOf(struct{}{})

	// slice of many types that can be used in tests.
	// Can be used with utils.Filter to fine tune the types.
	// keep depth low to prevent exponential growth of generated types and slowdown of tests.
	types = allTypes(2)
)

func allTypes(depth int) []reflect.Type {
	depth--
	if depth < 0 {
		return []reflect.Type{}
	}
	return utils.ConcatSlices[reflect.Type](
		primitiveTypes,
		[]reflect.Type{emptyStructTypes},
		arrayTypes(depth),
		sliceTypes(depth),
		mapTypes(depth),
		structTypes(depth),
	)
}
func structTypes(depth int) []reflect.Type {
	allTypes := allTypes(depth)
	results := make([]reflect.Type, len(allTypes))
	for i, baseType := range allTypes {
		fields := []reflect.StructField{{
			PkgPath: "tmp_test_pkg",
			Name:    "Field",
			Type:    baseType,
		}}
		results[i] = reflect.StructOf(fields)
	}
	return results
}
func sliceTypes(depth int) []reflect.Type {
	allTypes := allTypes(depth)
	results := make([]reflect.Type, len(allTypes))
	for i, baseType := range allTypes {
		results[i] = reflect.SliceOf(baseType)
	}
	return results
}
func arrayTypes(depth int) []reflect.Type {
	allTypes := allTypes(depth)
	results := make([]reflect.Type, len(allTypes))
	for i, baseType := range allTypes {
		results[i] = reflect.ArrayOf(3, baseType)
	}
	return results
}
func mapTypes(depth int) []reflect.Type {
	allTypes := allTypes(depth)
	comparableTypes := utils.Filter(allTypes, func(t reflect.Type) bool {
		return t.Comparable()
	})
	results := make([]reflect.Type, len(allTypes)*len(comparableTypes))
	for i, keyType := range comparableTypes {
		for j, valueType := range allTypes {
			results[i*len(allTypes)+j] = reflect.MapOf(keyType, valueType)
		}
	}
	return results
}

func expectTypeToBeCompatible(t *testing.T, validator TypeSchemaValidator, testType reflect.Type, errTemplate string, args ...any) {
	if err := validator.WithType(testType).Validate(); err != nil {
		t.Errorf(errTemplate, args...)
		t.Error(err)
	}
}
func expectTypeToBeIncompatible(t *testing.T, validator TypeSchemaValidator, testType reflect.Type, errTemplate string, args ...any) {
	if err := validator.WithType(testType).Validate(); err == nil {
		t.Errorf(errTemplate, args...)
	}
}

func TestSchemaValidatorWithOptions(t *testing.T) {
	stringSchema := openapi3.NewStringSchema()
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *stringSchema, Options{}).
		// test call to with options
		WithOptions(Options{})
	errTemplate := "expect string schema with time format to be %s with %s type"
	expectTypeToBeCompatible(t, validator, stringType, errTemplate, "compatible", stringType)
	// omit the string type from all defined test types
	var nonStringTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != stringType
	})
	for _, nonStringType := range nonStringTypes {
		t.Run(nonStringType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonStringType, errTemplate, "incompatible", nonStringType)
		})
	}
}
