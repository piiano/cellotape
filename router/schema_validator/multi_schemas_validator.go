package schema_validator

import "github.com/piiano/restcontroller/router/utils"

type schemaValidation struct {
	context       TypeSchemaValidator
	originalIndex int
	logger        utils.InMemoryLogger
}

func validateMultipleSchemas(cs ...TypeSchemaValidator) ([]schemaValidation, []schemaValidation) {
	pass := make([]schemaValidation, 0)
	failed := make([]schemaValidation, 0)
	for i, c := range cs {
		logger := utils.NewInMemoryLoggerWithLevel(c.logLevel())
		c.WithLogger(logger)
		err := c.Validate()
		validation := schemaValidation{context: c, logger: logger, originalIndex: i}
		if err == nil {
			pass = append(pass, validation)
		}
		if err != nil {
			failed = append(pass, validation)
		}
	}
	return pass, failed
}
