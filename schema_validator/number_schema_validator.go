package schema_validator

import (
	"fmt"
	"github.com/piiano/restcontroller/utils"
	"reflect"
)

func (c typeSchemaValidatorContext) validateNumberSchema() utils.MultiError {
	errs := utils.NewErrorsCollector()
	if c.schema.Type != numberSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
	default:
		errs.AddIfNotNil(fmt.Errorf("number schema is incompatible with type %s", c.goType))
	}
	switch c.schema.Format {
	case floatFormat:
		if c.goType.Kind() != reflect.Float32 {
			errs.AddIfNotNil(fmt.Errorf("number schema with float format is not compatible with type %s", c.goType))
		}
	case doubleFormat:
		if c.goType.Kind() != reflect.Float64 {
			errs.AddIfNotNil(fmt.Errorf("number schema with double format is not compatible with type %s", c.goType))
		}
	}
	// TODO: check type compatability with Max, ExclusiveMax, Min, and ExclusiveMin
	return errs.ErrorOrNil()
}
