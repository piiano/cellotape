package schema_validator

func (c typeSchemaValidatorContext) validateSchemaAllOf() error {
	if c.schema.AllOf == nil {
		return nil
	}
	l := c.newLogger()
	for _, option := range c.schema.AllOf {
		l.LogIfNotNil(c.level, c.WithSchema(*option.Value).Validate())
	}
	return l.MustHaveNoErrorsf(schemaAllOfPropertyIncompatibleWithType(l.Errors(), len(c.schema.AllOf), c.goType))
}
