package router

import (
	"net/http"

	"github.com/piiano/cellotape/router/utils"
)

type OpenAPIRouter interface {
	// Use middlewares on the root handler
	Use(...Handler) OpenAPIRouter
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPIRouter
	// WithOperation attaches an MiddlewareHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) OpenAPIRouter

	// WithContentType add to the router implementation for encoding additional content types.
	//
	// By default, the OpenAPIRouter supports the ContentTypes defined by DefaultContentTypes.
	//
	// An implementation of ContentType add support for additional content types by implementing serialization and
	// deserialization.
	WithContentType(ContentType) OpenAPIRouter

	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec.
	//
	// Returns a http.Handler and nil error if all checks passed correctly
	//
	// The returned handler can be used easily with any go http web framework as it is implementing the builtin
	// http.Handler interface.
	AsHandler() (http.Handler, error)

	// Spec returns the OpenAPI spec used by the router.
	Spec() OpenAPISpec
}

type Group interface {
	// Use middlewares on the root handler
	Use(...Handler) Group
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) Group
	// WithOperation attaches a MiddlewareHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) Group
	// This to be able to access groups, handlers, and operations of groups added with the WithGroup(Group) Group method
	group() group
}

// NewOpenAPIRouter creates new OpenAPIRouter router with the provided OpenAPISpec.
//
// The OpenAPISpec can be obtained by calling any of the OpenAPISpec loading functions:
//
// 		- NewSpecFromData - Init an OpenAPISpec from bytes represent the spec JSON or YAML.
// 		- NewSpecFromFile - Init an OpenAPISpec from a path to a spec file.
// 		- NewSpec - Create an empty spec object that can be initialized programmatically.
//
// This function also initialize OpenAPIRouter with the DefaultOptions, you can also initialize the router with custom
// Options using the NewOpenAPIRouterWithOptions function instead.
func NewOpenAPIRouter(spec OpenAPISpec) OpenAPIRouter {
	return NewOpenAPIRouterWithOptions(spec, DefaultOptions())
}

// NewOpenAPIRouterWithOptions creates new OpenAPIRouter router with the provided OpenAPISpec and custom Options.
//
// OpenAPISpec can be obtained by calling any of the OpenAPISpec loading functions:
//
// 		- NewSpecFromData - Init an OpenAPISpec from bytes represent the spec JSON or YAML.
// 		- NewSpecFromFile - Init an OpenAPISpec from a path to a spec file.
// 		- NewSpec - Create an empty spec object that can be initialized programmatically.
//
// This function also receives an Options object that allow customizing the default behaviour of the router.
// Check the Options documentation for more details.
//
// If you want to use the OpenAPIRouter with the DefaultOptions, you can use the shorter NewOpenAPIRouter function instead.
func NewOpenAPIRouterWithOptions(spec OpenAPISpec, options Options) OpenAPIRouter {
	return &openapi{
		spec:         spec,
		options:      options,
		contentTypes: DefaultContentTypes(),
	}
}

func NewGroup() Group { return new(group) }

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
	return createMainRouterHandler(oa)
}
func (oa *openapi) Spec() OpenAPISpec {
	return oa.spec
}

func (g *group) Use(handlers ...Handler) Group {
	g.handlers = append(g.handlers, utils.Map(handlers, asHandlerModel)...)
	return g
}
func (g *group) WithGroup(group Group) Group {
	g.groups = append(g.groups, group.group())
	return g
}

func (g *group) WithOperation(id string, handlerFunc Handler, handlers ...Handler) Group {
	g.operations = append(g.operations, operation{
		id:       id,
		handlers: utils.Map(handlers, asHandlerModel),
		handler:  asHandlerModel(handlerFunc),
	})
	return g
}

func (g *group) group() group { return *g }

func asHandlerModel(h Handler) handler {
	return handler{
		handlerFunc:    h,
		request:        h.requestTypes(),
		responses:      h.responseTypes(),
		sourcePosition: h.sourcePosition(),
	}
}
