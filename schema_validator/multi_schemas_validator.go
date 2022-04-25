package schema_validator

import (
	"github.com/piiano/restcontroller/utils"
)

type schemaValidation struct {
	context       TypeSchemaValidator
	originalIndex int
	multiError    utils.MultiError
}

func validateMultipleSchemas(cs ...TypeSchemaValidator) ([]schemaValidation, []schemaValidation) {
	pass := make([]schemaValidation, 0, len(cs))
	failed := make([]schemaValidation, 0, len(cs))
	for i, c := range cs {
		err := c.Validate()
		validation := schemaValidation{context: c, multiError: err, originalIndex: i}
		if err == nil {
			pass = append(pass, validation)
		}
		if err != nil {
			failed = append(pass, validation)
		}
	}
	return pass, failed
}
