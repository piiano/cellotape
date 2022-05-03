package router

type Group interface {
	// Use middlewares on the root handler
	Use(...Handler) Group
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) Group
	// WithOperation attaches an Handler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) Group
	// This to be able to access groups, handlers, and operations of groups added with the WithGroup(Group) Group method
	group() group
}

type group struct {
	groups     []group
	handlers   []Handler
	operations map[string]operation
}

func NewGroup() Group {
	return &group{
		handlers:   []Handler{},
		groups:     []group{},
		operations: map[string]operation{},
	}
}

func (g *group) group() group { return *g }
func (g *group) Use(handlers ...Handler) Group {
	g.handlers = append(g.handlers, handlers...)
	return g
}
func (g *group) WithGroup(group Group) Group {
	g.groups = append(g.groups, group.group())
	return g
}
func (g *group) WithOperation(id string, handlerFunc Handler, handlers ...Handler) Group {
	g.operations[id] = operation{
		id:               id,
		responseTypes:    handlerFunc.responseTypes(),
		operationHandler: handlerFunc,
		handlers:         handlers,
	}
	return g
}
