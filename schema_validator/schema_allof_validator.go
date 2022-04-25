package schema_validator

import "github.com/piiano/restcontroller/utils"

func (c typeSchemaValidatorContext) validateSchemaAllOf() utils.MultiError {
	if c.schema.AllOf == nil {
		return nil
	}
	errs := utils.NewErrorsCollector()
	for _, option := range c.schema.AllOf {
		errs.AddIfNotNil(c.WithSchema(*option.Value).Validate())
	}
	return errs.ErrorOrNil()
}
