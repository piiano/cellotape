package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/utils"
)

func (c typeSchemaValidatorContext) validateSchemaAnyOf() error {
	if c.schema.AnyOf == nil {
		return nil
	}
	//l := c.newLogger()
	l := utils.NewInMemoryLoggerWithLevel(c.level)
	pass, failed := validateMultipleSchemas(utils.Map(c.schema.AnyOf, func(t *openapi3.SchemaRef) TypeSchemaValidator {
		return c.WithSchema(*t.Value)
	})...)
	if len(pass) == 0 {
		l.Logf(c.level, "schema with anyOf property is incompatible with type %s", c.goType)
		for i, check := range failed {
			l.Logf(c.level, "anyOf[%d] didn't match type %s", i, c.goType)
			l.Log(c.level, check.logger.Printed())
		}
	}
	return l.MustHaveNoErrorsf(schemaAnyOfPropertyIncompatibleWithType(len(c.schema.AnyOf), c.goType))
}
