package router

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/schema_validator"
	"github.com/piiano/cellotape/router/utils"
)

type ContentType interface {
	Mime() string
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
	ValidateTypeSchema(utils.Logger, utils.LogLevel, reflect.Type, openapi3.Schema) error
}

type ContentTypes map[string]ContentType

type OctetStreamContentType struct{}

func (t OctetStreamContentType) Mime() string { return "application/octet-stream" }
func (t OctetStreamContentType) Encode(value any) ([]byte, error) {
	if value == nil {
		return []byte{}, nil
	}
	if bytesSlice, ok := value.([]byte); ok {
		return bytesSlice, nil
	}
	return nil, fmt.Errorf("type %T is incompatible with content type %q. value must be a []byte", value, t.Mime())
}
func (t OctetStreamContentType) Decode(data []byte, value any) error {
	if len(data) == 0 {
		return nil
	}
	switch typedValue := value.(type) {
	case *any:
		*typedValue = data
	case *[]byte:
		*typedValue = data
	default:
		return fmt.Errorf("type %T is incompatible with content type %q. value must be a *[]byte", value, t.Mime())
	}
	return nil
}
func (t OctetStreamContentType) ValidateTypeSchema(
	logger utils.Logger, level utils.LogLevel, goType reflect.Type, schema openapi3.Schema) error {
	if goType != reflect.TypeOf([]byte{}) {
		logger.Logf(level, "type %s is incompatible with content type %q", goType, t.Mime())
	}
	if schema.Type != "string" || schema.Format != "binary" {
		logger.Logf(level, `schema must have a "string" type with a "binary" format when content type is %q`, t.Mime())
	}
	return logger.MustHaveNoErrors()
}

type PlainTextContentType struct{}

func (t PlainTextContentType) Mime() string { return "text/plain" }
func (t PlainTextContentType) Encode(value any) ([]byte, error) {
	if value == nil {
		return []byte{}, nil
	}
	if str, ok := value.(string); ok {
		return []byte(str), nil
	}
	return nil, fmt.Errorf("type %T is incompatible with content type %q. value must be a string", value, t.Mime())
}
func (t PlainTextContentType) Decode(data []byte, value any) error {
	if len(data) == 0 {
		return nil
	}
	switch typedValue := value.(type) {
	case *any:
		*typedValue = string(data)
	case *string:
		*typedValue = string(data)
	default:
		return fmt.Errorf("type %T is incompatible with content type %q. value must be a *string", value, t.Mime())
	}
	return nil
}
func (t PlainTextContentType) ValidateTypeSchema(
	logger utils.Logger, level utils.LogLevel, goType reflect.Type, schema openapi3.Schema) error {
	switch goType {
	case reflect.TypeOf(""):
	case reflect.PointerTo(reflect.TypeOf("")):
	default:
		logger.Logf(level, "type %s is incompatible with content type %q", goType, t.Mime())
	}
	if schema.Type != "string" {
		logger.Logf(level, "schema type %s is incompatible with content type %q", schema.Type, t.Mime())
	}
	return logger.MustHaveNoErrors()
}

type JSONContentType struct{}

func (t JSONContentType) Mime() string                        { return "application/json" }
func (t JSONContentType) Encode(value any) ([]byte, error)    { return json.Marshal(value) }
func (t JSONContentType) Decode(data []byte, value any) error { return json.Unmarshal(data, value) }
func (t JSONContentType) ValidateTypeSchema(
	logger utils.Logger, level utils.LogLevel, goType reflect.Type, schema openapi3.Schema) error {
	validator := schema_validator.NewTypeSchemaValidator(goType, schema)
	err := validator.Validate()
	for _, errMessage := range validator.Errors() {
		logger.Log(level, errMessage)
	}

	return err
}

func DefaultContentTypes() ContentTypes {
	defaultContentTypes := []ContentType{
		OctetStreamContentType{},
		PlainTextContentType{},
		JSONContentType{},
	}
	contentTypes := make(ContentTypes, len(defaultContentTypes))
	for _, contentType := range defaultContentTypes {
		contentTypes[contentType.Mime()] = contentType
	}
	return contentTypes
}
