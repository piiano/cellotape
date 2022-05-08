package schema_validator

import (
	"fmt"
	"github.com/piiano/restcontroller/router/utils"
	"reflect"
)

func (c typeSchemaValidatorContext) validateIntegerSchema() utils.MultiError {
	errs := utils.NewErrorsCollector()
	if c.schema.Type != integerSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		errs.AddErrorsIfNotNil(fmt.Errorf("type %s is not compatible with integer schema", c.goType))
	}
	switch c.schema.Format {
	case int32Format:
		if c.goType.Kind() != reflect.Int32 {
			errs.AddErrorsIfNotNil(fmt.Errorf("type %s is not compatible with integer schema with int32 format", c.goType))
		}
	case int64Format:
		if c.goType.Kind() != reflect.Int64 {
			errs.AddErrorsIfNotNil(fmt.Errorf("type %s is not compatible with integer schema with int64 format", c.goType))
		}
	}
	// TODO: check type compatability with Max, ExclusiveMax, Min, and ExclusiveMin
	return errs.ErrorOrNil()
}
