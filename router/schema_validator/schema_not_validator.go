package schema_validator

import (
	"fmt"
)

func (c typeSchemaValidatorContext) validateSchemaNot() error {
	if c.schema.Not == nil {
		return nil
	}
	if err := c.WithSchema(*c.schema.Not.Value).Validate(); err == nil {
		return fmt.Errorf("schema with not property is incompatible with type %s", c.goType)
	}
	return nil
}
