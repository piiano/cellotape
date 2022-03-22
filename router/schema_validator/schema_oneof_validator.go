package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/cellotape/router/utils"
)

func (c typeSchemaValidatorContext) validateSchemaOneOf() error {
	if c.schema.OneOf == nil {
		return nil
	}
	l := c.newLogger()
	pass, failed := validateMultipleSchemas(utils.Map(c.schema.OneOf, func(t *openapi3.SchemaRef) TypeSchemaValidator {
		return c.WithSchema(*t.Value)
	})...)
	if len(pass) == 0 {
		l.Logf(c.level, "schema with oneOf property has no matches for the type %q", c.goType)
		for _, check := range failed {
			l.Logf(c.level, "oneOf[%d] didn't match type %q", check.originalIndex, c.goType)
			l.Log(c.level, check.logger.Printed())
		}
	}
	if len(pass) > 1 {
		l.Logf(c.level, "schema with oneOf property has more than one match for the type %q", c.goType)
		for _, check := range pass {
			l.Logf(c.level, "oneOf[%d] matched type %q", check.originalIndex, c.goType)
		}
	}
	return l.MustHaveNoErrorsf("schema with oneOf property is incompatible with type %s", c.goType)
}
