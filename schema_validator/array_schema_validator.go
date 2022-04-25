package schema_validator

import (
	"fmt"
	"github.com/piiano/restcontroller/utils"
	"reflect"
)

func (c typeSchemaValidatorContext) validateArraySchema() utils.MultiError {
	if SchemaType(c.schema.Type) != arraySchemaType {
		return nil
	}
	errs := utils.NewErrorsCollector()
	if c.goType.Kind() != reflect.Array && c.goType.Kind() != reflect.Slice {
		errs.AddIfNotNil(fmt.Errorf("schema %q must be used with slice or array but is used with %s", c.schema.Title, c.goType))
		return errs.ErrorOrNil()
	}
	errs.AddIfNotNil(c.WithSchema(*c.schema.Items.Value).WithType(c.goType.Elem()).Validate())
	return errs.ErrorOrNil()
}
