package router

import (
	"github.com/piiano/restcontroller/router/utils"
	"net/http"
)

type OpenAPIRouter interface {
	// Use middlewares on the root handler
	Use(...Handler) OpenAPIRouter
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPIRouter
	// WithOperation attaches an MiddlewareHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) OpenAPIRouter
	// WithContentType Add support for encoding additional content type
	WithContentType(ContentType) OpenAPIRouter
	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec
	// returns a http.Handler that implements nil value if all checks passed correctly
	AsHandler() (http.Handler, error)
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

// NewOpenAPIRouter creates new OpenAPIRouter router with the provided OpenAPISpec
func NewOpenAPIRouter(spec OpenAPISpec) OpenAPIRouter {
	return NewOpenAPIRouterWithOptions(spec, DefaultOptions())
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

func NewGroup() Group {
	return &group{
		handlers:   []handler{},
		groups:     []group{},
		operations: []operation{},
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
