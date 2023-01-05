package schema_validator

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/utils"
)

func (c typeSchemaValidatorContext) matchAllSchemaValidator(name string, schemas openapi3.SchemaRefs) {
	if schemas == nil {
		return
	}

	types := []reflect.Type{c.goType}

	if utils.IsMultiType(c.goType) {
		if multiTypeTypes, err := utils.ExtractMultiTypeTypes(c.goType); err != nil {
			c.err(err.Error())
		} else {
			types = multiTypeTypes
		}
	}

	usedTypes := utils.NewSet[string]()

schemas:
	for index, schema := range schemas {
		schemaErrors := make([]string, 0)

		for _, multiTypeType := range types {
			typeValidator := typeSchemaValidatorContext{
				errors: new([]string),
				schema: *schema.Value,
				goType: multiTypeType,
			}

			if err := typeValidator.Validate(); err == nil {
				usedTypes.Add(multiTypeType.String())
				continue schemas
			}

			schemaErrors = append(schemaErrors, *typeValidator.errors...)
		}

		c.err("%s schema at index %d didn't match type of %q", name, index, c.goType)
		*c.errors = append(*c.errors, schemaErrors...)
	}

	if utils.IsMultiType(c.goType) {
		for _, multiTypeType := range types {
			if !usedTypes.Has(multiTypeType.String()) {
				c.err("non of %s schemas match type %q of %q", name, multiTypeType, c.goType)
			}
		}
	}
}
