package schema_validator

func (c typeSchemaValidatorContext) validateSchemaAllOf() {
	if c.schema.AllOf == nil {
		return
	}
	errors := len(*c.errors)
	for _, option := range c.schema.AllOf {
		_ = c.WithSchema(*option.Value).Validate()
	}
	if len(*c.errors) > errors {
		c.err(schemaAllOfPropertyIncompatibleWithType(len(*c.errors)-errors, len(c.schema.AllOf), c.goType))
	}
}
