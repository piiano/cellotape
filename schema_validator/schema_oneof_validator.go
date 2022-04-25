package schema_validator

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/utils"
	"strings"
)

func (c typeSchemaValidatorContext) validateSchemaOneOf() utils.MultiError {
	if c.schema.OneOf == nil {
		return nil
	}
	errs := utils.NewErrorsCollector()
	pass, failed := validateMultipleSchemas(utils.Map(c.schema.OneOf, func(t *openapi3.SchemaRef) TypeSchemaValidator {
		return c.WithSchema(*t.Value)
	})...)
	if len(pass) == 0 {
		lines := make([]string, len(failed))
		lines = append(lines, fmt.Sprintf("schema with oneOf property has no matches for the type %q", c.goType))
		for _, check := range failed {
			lines = append(lines, fmt.Sprintf("oneOf[%d] didn't match type %q", check.originalIndex, c.goType))
			lines = append(lines, check.multiError.Error())
		}
		errs.AddIfNotNil(errors.New(strings.Join(lines, "\n")))
	}
	if len(pass) > 1 {
		lines := make([]string, len(pass))
		lines = append(lines, fmt.Sprintf("schema with oneOf property has more than one match for the type %q", c.goType))
		for _, check := range pass {
			lines = append(lines, fmt.Sprintf("oneOf[%d] matched type %q", check.originalIndex, c.goType))
		}
		lines = append(lines, "")
		errs.AddIfNotNil(errors.New(strings.Join(lines, "\n")))
	}
	return errs.ErrorOrNil()
}
