package schema_validator

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/utils"
	"strings"
)

func (c typeSchemaValidatorContext) validateSchemaAnyOf() utils.MultiError {
	if c.schema.AnyOf == nil {
		return nil
	}
	errs := utils.NewErrorsCollector()
	pass, failed := validateMultipleSchemas(utils.Map(c.schema.AnyOf, func(t *openapi3.SchemaRef) TypeSchemaValidator {
		return c.WithSchema(*t.Value)
	})...)
	if len(pass) == 0 {
		lines := make([]string, len(failed))
		lines = append(lines, fmt.Sprintf("schema with anyOf property is not compatible with type %s", c.goType))
		for i, check := range failed {
			lines = append(lines, fmt.Sprintf("anyOf[%d] didn't match type %s", i, c.goType))
			lines = append(lines, check.multiError.Error())
		}
		errs.AddErrorsIfNotNil(errors.New(strings.Join(lines, "\n")))
	}
	return errs.ErrorOrNil()
}
