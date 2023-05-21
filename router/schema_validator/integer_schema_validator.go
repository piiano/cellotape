package schema_validator

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

func (c typeSchemaValidatorContext) validateIntegerSchema() {

	if c.schema.Type != openapi3.TypeInteger {
		return
	}
	switch c.goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
	switch c.schema.Format {
	case int32Format:
		if c.goType.Kind() != reflect.Int32 {
			c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	case int64Format:
		if c.goType.Kind() != reflect.Int64 {
			c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	}

	// TODO: check type compatability with Max, ExclusiveMax, Min, and ExclusiveMin
}
