package schema_validator

import (
	"fmt"
	"github.com/piiano/restcontroller/utils"
	"reflect"
)

func (c typeSchemaValidatorContext) validateBooleanSchema() utils.MultiError {
	errs := utils.NewErrorsCollector()
	if c.schema.Type != booleanSchemaType {
		return nil
	}
	if c.goType.Kind() != reflect.Bool {
		errs.AddIfNotNil(fmt.Errorf("boolean schema must be of type bool. type %s is incompatible", c.goType))
	}
	return errs.ErrorOrNil()
}
