package router

import "reflect"

// openapi describes the internal state of the OpenAPIRouter builder
type openapi struct {
	// spec added when creating a new router with NewOpenAPIRouter
	spec OpenAPISpec
	// options added when creating a new router with NewOpenAPIRouterWithOptions
	options Options
	// contentTypes added with WithContentType
	contentTypes ContentTypes
	// group hold internal resources added by OpenAPIRouter.Use, OpenAPIRouter.WithGroup and OpenAPIRouter.WithOperation
	group
}

// group describes the internal state of the Group builder
type group struct {
	// nested groups added by Group.WithGroup
	groups []group
	// handlers added by Group.Use
	handlers []handler
	// operations added by Group.WithOperation
	operations []operation
}

// operation describes the internal state of group and openapi operations added by Group.WithOperation and OpenAPIRouter.WithOperation
type operation struct {
	// the operation id added with OpenAPIRouter.WithOperation
	id string
	// handler hold the internal representation of the operationFunc added with OpenAPIRouter.WithOperation
	handler
	// handlers hold the additional handlers added with the variadic arguments of OpenAPIRouter.WithOperation
	handlers []handler
}

// handler describes the internal state of operation, group and openapi handlers added by Group.WithOperation,
// OpenAPIRouter.WithOperation, Group.Use and OpenAPIRouter.Use
type handler struct {
	// the handler function implementing the Handler interface (can be operationFunc or handlerFunc)
	handlerFunc Handler
	// hold a representation of the request described in the handler parameters
	request requestTypes
	// hold a representation of the responses described in the handler returned type
	responses handlerResponses
	// sourcePosition hold the location in source of the handler function for showing meaningful helpful message for validation errors
	sourcePosition sourcePosition
}

// requestTypes described the parameter types provided in the Request input of a handler function
type requestTypes struct {
	// requestBody is the type of the body parameter. type is nilType if there is no httpRequest body
	requestBody reflect.Type
	// pathParams is the type of the body parameter. type is nilType if there is no path pathParams
	pathParams reflect.Type
	// queryParams is the type of the body parameter.  type is nilType if there is no query pathParams
	queryParams reflect.Type
}

// handlerResponses hold a representation of the responses described in a handler returned type
type handlerResponses map[int]httpResponse

// httpResponse describes a single possible response described by a handler returned responses struct type
type httpResponse struct {
	// status describes the status extracted from the status tag in the responses' struct type
	status int
	// responseType describes the type of the field corresponding with that response
	responseType reflect.Type
	// fieldIndex is the index used to access the response field with reflection in runtime
	fieldIndex []int
	// isNilType is true if the response field type is Nil to sign that this response has no content
	isNilType bool
}
