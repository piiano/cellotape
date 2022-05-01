package schema_validator

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"time"
)

func (c typeSchemaValidatorContext) validateStringSchema() utils.MultiError {
	errs := utils.NewErrorsCollector()
	if c.schema.Type != stringSchemaType {
		return nil
	}
	switch c.schema.Format {
	case "":
		if c.goType.Kind() != reflect.String {
			errs.AddIfNotNil(fmt.Errorf("received type %s for string schema", c.goType))
		}
	case uuidFormat:
		if c.goType.Kind() != reflect.String && c.goType != reflect.TypeOf(uuid.New()) {
			errs.AddIfNotNil(fmt.Errorf("string schema with uuid format is not compatible with type %s", c.goType))
		}
	case timeFormat:
		if c.goType.Kind() != reflect.String && c.goType != reflect.TypeOf(time.Now()) {
			errs.AddIfNotNil(fmt.Errorf("string schema with time format is not compatible with type %s", c.goType))
		}
	// TODO: add support for more formats compatible types (dateTimeFormat, dateFormat, durationFormat, etc.)
	case dateTimeFormat, dateFormat, durationFormat, emailFormat, idnEmailFormat, hostnameFormat,
		idnHostnameFormat, ipv4Format, ipv6Format, uriFormat, uriReferenceFormat, iriFormat, iriReferenceFormat,
		uriTemplateFormat, jsonPointerFormat, relativeJsonPointerFormat, regexFormat, passwordFormat:
		if c.goType.Kind() != reflect.String {
			errs.AddIfNotNil(fmt.Errorf("string schema with %s format is not compatible with type %s", c.schema.Format, c.goType))
		}
	}
	return errs.ErrorOrNil()
}
