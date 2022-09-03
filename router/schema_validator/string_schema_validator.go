package schema_validator

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(new(time.Time)).Elem()
var uuidType = reflect.TypeOf(new(uuid.UUID)).Elem()

func (c typeSchemaValidatorContext) validateStringSchema() error {
	l := c.newLogger()
	if c.schema.Type != stringSchemaType {
		return nil
	}
	switch c.schema.Format {
	case "":
		if c.goType.Kind() != reflect.String {
			l.Logf(c.level, schemaTypeIsIncompatibleWithType(c.schema, c.goType))
		}
	case uuidFormat:
		if c.goType.Kind() != reflect.String && !uuidType.ConvertibleTo(c.goType) {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	case dateTimeFormat, timeFormat:
		if c.goType.Kind() != reflect.String && !timeType.ConvertibleTo(c.goType) {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	// TODO: add support for more formats compatible types (dateTimeFormat, dateFormat, durationFormat, etc.)
	case dateFormat, durationFormat, emailFormat, idnEmailFormat, hostnameFormat,
		idnHostnameFormat, ipv4Format, ipv6Format, uriFormat, uriReferenceFormat, iriFormat, iriReferenceFormat,
		uriTemplateFormat, jsonPointerFormat, relativeJsonPointerFormat, regexFormat, passwordFormat:
		if c.goType.Kind() != reflect.String {
			l.Logf(c.level, schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
	}
	return formatMustHaveNoError(l.MustHaveNoErrors(), c.schema.Type, c.goType)
}
