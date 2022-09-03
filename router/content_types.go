package router

import (
	"encoding/json"
	"fmt"
)

type ContentType interface {
	Mime() string
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
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

type JSONContentType struct{}

func (t JSONContentType) Mime() string                        { return "application/json" }
func (t JSONContentType) Encode(value any) ([]byte, error)    { return json.Marshal(value) }
func (t JSONContentType) Decode(data []byte, value any) error { return json.Unmarshal(data, value) }

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
