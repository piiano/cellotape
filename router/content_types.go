package router

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type ContentType interface {
	Mime() string
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
}

type ContentTypes map[string]ContentType

type JsonContentType struct{}

func (t JsonContentType) Mime() string                        { return "application/json" }
func (t JsonContentType) Encode(value any) ([]byte, error)    { return json.Marshal(value) }
func (t JsonContentType) Decode(data []byte, value any) error { return json.Unmarshal(data, value) }

type YamlContentType struct{}

func (t YamlContentType) Mime() string                        { return "application/yaml" }
func (t YamlContentType) Encode(value any) ([]byte, error)    { return yaml.Marshal(value) }
func (t YamlContentType) Decode(data []byte, value any) error { return yaml.Unmarshal(data, value) }

type TomlContentType struct{}

func (t TomlContentType) Mime() string { return "application/toml" }
func (t TomlContentType) Encode(value any) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	err := toml.NewEncoder(buffer).Encode(value)
	return buffer.Bytes(), err
}
func (t TomlContentType) Decode(data []byte, value any) error { return toml.Unmarshal(data, value) }

func DefaultContentTypes() ContentTypes {
	defaultContentTypes := []ContentType{
		JsonContentType{},
		YamlContentType{},
		TomlContentType{},
	}
	contentTypes := make(ContentTypes, len(defaultContentTypes))
	for _, contentType := range defaultContentTypes {
		contentTypes[contentType.Mime()] = contentType
	}
	return contentTypes
}
