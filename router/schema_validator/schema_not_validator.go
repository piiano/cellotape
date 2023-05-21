package schema_validator

func (c typeSchemaValidatorContext) validateSchemaNot() {
	if c.schema.Not == nil {
		return
	}
	errors := len(*c.errors)
	if err := c.WithSchema(*c.schema.Not.Value).Validate(); err == nil {
		c.err("schema with not property is incompatible with type %s", c.goType)
		return
	}
	*c.errors = (*c.errors)[:errors]
}
