package router

import "encoding/json"

type ContentType interface {
	Mime() string
	Marshal(value any) ([]byte, error)
	Unmarshal([]byte, any) error
}

type ContentTypes map[string]ContentType

type JsonContentType struct{}

func (t JsonContentType) Mime() string                          { return "application/json" }
func (t JsonContentType) Marshal(value any) ([]byte, error)     { return json.Marshal(value) }
func (t JsonContentType) Unmarshal(bytes []byte, dst any) error { return json.Unmarshal(bytes, dst) }

func DefaultContentTypes() ContentTypes {
	defaultContentTypes := []ContentType{JsonContentType{}}
	contentTypes := make(ContentTypes, len(defaultContentTypes))
	for _, contentType := range defaultContentTypes {
		contentTypes[contentType.Mime()] = contentType
	}
	return contentTypes
}
