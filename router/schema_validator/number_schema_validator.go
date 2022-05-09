package schema_validator

import (
	"reflect"
)

func (c typeSchemaValidatorContext) validateNumberSchema() error {
	l := c.newLogger()
	if c.schema.Type != numberSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
	default:
		l.Logf(c.level, schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
	switch c.schema.Format {
	case floatFormat:
		if c.goType.Kind() != reflect.Float32 {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	case doubleFormat:
		if c.goType.Kind() != reflect.Float64 {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	}
	// TODO: check type compatability with Max, ExclusiveMax, Min, and ExclusiveMin
	return formatMustHaveNoError(l.MustHaveNoErrors(), c.schema.Type, c.goType)
}
