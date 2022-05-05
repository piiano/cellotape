package router

import (
	"net/http"
)

type OpenAPIRouter interface {
	// Use middlewares on the root handler
	Use(...Handler) OpenAPIRouter
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPIRouter
	// WithOperation attaches an Handler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) OpenAPIRouter
	// WithContentType Add support for encoding additional content type
	WithContentType(ContentType) OpenAPIRouter
	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec
	// returns a http.Handler that implements nil value if all checks passed correctly
	AsHandler() (http.Handler, error)
}

// NewOpenAPIRouter creates new OpenAPIRouter router with the provided OpenAPISpec
func NewOpenAPIRouter(spec OpenAPISpec) OpenAPIRouter {
	return NewOpenAPIRouterWithOptions(spec, Options{
		RecoverOnPanic: true,
	})
}

// NewOpenAPIRouterWithOptions creates new OpenAPIRouter router with the provided OpenAPISpec and custom options
func NewOpenAPIRouterWithOptions(spec OpenAPISpec, options Options) OpenAPIRouter {
	return &openapi{
		spec:         spec,
		options:      options,
		contentTypes: DefaultContentTypes(),
		group: group{
			handlers:   []handler{},
			groups:     []group{},
			operations: []operation{},
		},
	}
}

func (oa *openapi) Use(handlers ...Handler) OpenAPIRouter {
	oa.group.Use(handlers...)
	return oa
}
func (oa *openapi) WithGroup(group Group) OpenAPIRouter {
	oa.group.WithGroup(group)
	return oa
}
func (oa *openapi) WithOperation(id string, handlerFunc Handler, handlers ...Handler) OpenAPIRouter {
	oa.group.WithOperation(id, handlerFunc, handlers...)
	return oa
}
func (oa *openapi) WithContentType(contentType ContentType) OpenAPIRouter {
	oa.contentTypes[contentType.Mime()] = contentType
	return oa
}
func (oa *openapi) AsHandler() (http.Handler, error) {
	return asHandler(oa)
}
