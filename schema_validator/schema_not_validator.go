package schema_validator

import (
	"fmt"
	"github.com/piiano/restcontroller/utils"
)

func (c typeSchemaValidatorContext) validateSchemaNot() utils.MultiError {
	if c.schema.Not == nil {
		return nil
	}
	errs := utils.NewErrorsCollector()
	if err := c.WithSchema(*c.schema.Not.Value).Validate(); err == nil {
		errs.AddIfNotNil(fmt.Errorf("schema %q with a not valudation failed for type %s", c.schema.Title, c.goType))
	}
	return errs.ErrorOrNil()
}
