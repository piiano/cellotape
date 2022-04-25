package router

import (
	"fmt"
	"net/http"
)

type Group interface {
	// Use middlewares on the root handler
	Use(...http.Handler) Group
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) Group
	// WithOperation attaches an OperationHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, OperationHandler, ...http.Handler) Group

	getErrors() []error
	getOperations() map[string]operationHandler
	getHandlers() []http.Handler
}

func NewGroup() Group {
	return group{
		operations: map[string]operationHandler{},
		handlers:   []http.Handler{},
		errors:     []error{},
	}
}

type group struct {
	errors     []error
	operations map[string]operationHandler
	handlers   []http.Handler
}

func (g group) getErrors() []error {
	return g.errors
}
func (g group) getOperations() map[string]operationHandler {
	return g.operations
}
func (g group) getHandlers() []http.Handler {
	return g.handlers
}

func (g *group) appendErrIfNotNil(err error) {
	if err != nil {
		g.errors = append(g.errors, err)
	}
}
func (g group) Use(handlers ...http.Handler) Group {
	g.handlers = append(g.handlers, handlers...)
	return g
}
func (g group) WithGroup(group Group) Group {
	for _, err := range group.getErrors() {
		g.appendErrIfNotNil(err)
	}
	handlers := group.getHandlers()
	for id, oh := range group.getOperations() {
		if _, ok := g.operations[id]; ok {
			g.appendErrIfNotNil(fmt.Errorf("operation with id %q already declared", id))
			continue
		}
		g.operations[id] = operationHandler{
			handlers:         append(handlers, oh.handlers...),
			operationHandler: oh.operationHandler,
		}
	}
	return g
}

func (g group) WithOperation(id string, handler OperationHandler, handlers ...http.Handler) Group {
	g.operations[id] = operationHandler{
		operationHandler: handler,
		handlers:         handlers,
	}
	return g
}
