package router

import (
	"github.com/piiano/restcontroller/utils"
	"net/http"
)

type Group interface {
	// Use middlewares on the root handler
	Use(...http.HandlerFunc) Group
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) Group
	// WithOperation attaches an OperationHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, OperationHandler, ...http.HandlerFunc) Group
	groupInternals
}

func NewGroup() Group {
	return &group{
		handlers:   []http.HandlerFunc{},
		groups:     []groupInternals{},
		operations: map[string]Operation{},
	}
}

type group struct {
	groups     []groupInternals
	handlers   []http.HandlerFunc
	operations map[string]Operation
}

func (g *group) Use(handlers ...http.HandlerFunc) Group {
	g.handlers = append(g.handlers, handlers...)
	return g
}
func (g *group) WithGroup(group Group) Group {
	g.groups = append(g.groups, group)
	return g
}
func (g *group) WithOperation(id string, handlerFunc OperationHandler, handlers ...http.HandlerFunc) Group {
	g.operations[id] = operation{
		operationHandler: handlerFunc,
		handlers:         handlers,
	}
	return g
}

func (g group) getGroups() []groupInternals         { return g.groups }
func (g group) getHandlers() []http.HandlerFunc     { return g.handlers }
func (g group) getOperations() map[string]Operation { return g.operations }

type groupInternals interface {
	getGroups() []groupInternals
	getHandlers() []http.HandlerFunc
	getOperations() map[string]Operation
}

// flattenOperations takes a group with separate operations, handlers, and nested groups and flatten them into a flat
// Operation slice that include for each Operation its own OperationFunc, attached handlers, and attached group handlers.
func flattenOperations(g groupInternals) []Operation {
	flatOperations := make([]Operation, 0)
	for id, op := range g.getOperations() {
		flatOperations = append(flatOperations, operation{
			id:               id,
			handlers:         utils.ConcatSlices(g.getHandlers(), op.getHandlers()),
			operationHandler: op.getHandlerFunc(),
		})
	}
	for _, nestedGroup := range g.getGroups() {
		for _, flatOp := range flattenOperations(nestedGroup) {
			flatOperations = append(flatOperations, operation{
				id:               flatOp.getId(),
				handlers:         utils.ConcatSlices(g.getHandlers(), flatOp.getHandlers()),
				operationHandler: flatOp.getHandlerFunc(),
			})
		}
	}
	return flatOperations
}
