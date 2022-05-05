package router

import (
	"github.com/piiano/restcontroller/utils"
)

type Group interface {
	// Use middlewares on the root handler
	Use(...Handler) Group
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) Group
	// WithOperation attaches a Handler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) Group
	// This to be able to access groups, handlers, and operations of groups added with the WithGroup(Group) Group method
	group() group
}

func NewGroup() Group {
	return &group{
		handlers:   []handler{},
		groups:     []group{},
		operations: []operation{},
	}
}

func (g *group) Use(handlers ...Handler) Group {
	g.handlers = append(g.handlers, utils.Map(handlers, func(h Handler) handler {
		return handler{
			handlerFunc:    h,
			responses:      h.responseTypes(),
			sourcePosition: h.sourcePosition(),
		}
	})...)
	return g
}
func (g *group) WithGroup(group Group) Group {
	g.groups = append(g.groups, group.group())
	return g
}

func (g *group) WithOperation(id string, handlerFunc Handler, handlers ...Handler) Group {
	g.operations = append(g.operations, operation{
		id: id,
		handlers: utils.Map(handlers, func(h Handler) handler {
			return handler{
				handlerFunc:    h,
				responses:      h.responseTypes(),
				sourcePosition: h.sourcePosition(),
			}
		}),
		handler: handler{
			handlerFunc:    handlerFunc,
			responses:      handlerFunc.responseTypes(),
			sourcePosition: handlerFunc.sourcePosition(),
		},
	})
	return g
}
func (g *group) group() group { return *g }
