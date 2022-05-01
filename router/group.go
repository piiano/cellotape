package router

import (
	"net/http"
)

type Group interface {
	// Use middlewares on the root handler
	Use(...http.HandlerFunc) Group
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) Group
	// WithOperation attaches an OperationHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, OperationHandler, ...http.HandlerFunc) Group
	// This to be able to access groups, handlers, and operations of groups added with the WithGroup(Group) Group method
	group() group
}

type group struct {
	groups     []group
	handlers   []http.HandlerFunc
	operations map[string]operation
}

func NewGroup() Group {
	return &group{
		handlers:   []http.HandlerFunc{},
		groups:     []group{},
		operations: map[string]operation{},
	}
}

func (g *group) group() group { return *g }
func (g *group) Use(handlers ...http.HandlerFunc) Group {
	g.handlers = append(g.handlers, handlers...)
	return g
}
func (g *group) WithGroup(group Group) Group {
	g.groups = append(g.groups, group.group())
	return g
}
func (g *group) WithOperation(id string, handlerFunc OperationHandler, handlers ...http.HandlerFunc) Group {
	g.operations[id] = operation{
		operationHandler: handlerFunc,
		handlers:         handlers,
	}
	return g
}
