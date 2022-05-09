package schema_validator

import (
	"reflect"
)

func (c typeSchemaValidatorContext) validateIntegerSchema() error {
	l := c.newLogger()
	if c.schema.Type != integerSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		l.Logf(c.level, schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
	switch c.schema.Format {
	case int32Format:
		if c.goType.Kind() != reflect.Int32 {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	case int64Format:
		if c.goType.Kind() != reflect.Int64 {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	}
	// TODO: check type compatability with Max, ExclusiveMax, Min, and ExclusiveMin
	return formatMustHaveNoError(l.MustHaveNoErrors(), c.schema.Type, c.goType)
}
