package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/piiano/restcontroller/utils"
	"log"
	"net/http"
	"regexp"
)

func asHandler(oa *openapi) (http.Handler, error) {
	flatOperations := flattenOperations(oa.group)
	if err := validateAllOperations(oa, flatOperations); err != nil {
		return nil, err
	}
	router := httprouter.New()
	pathParamsMatcher := regexp.MustCompile(`\{([^/]*)}`)
	for _, flatOp := range flatOperations {
		specOp, _ := oa.spec.findSpecOperationByID(flatOp.id)
		path := pathParamsMatcher.ReplaceAllString(specOp.path, ":$1")
		chainHead := chainHandlers(*oa, append(flatOp.handlers, flatOp.handler)...)
		httpRouterHandler := asHttpRouterHandler(*oa, chainHead)
		router.Handle(specOp.method, path, httpRouterHandler)
	}
	return router, nil
}

func validateAllOperations(oa *openapi, flatOperations []operation) error {
	errs := utils.NewErrorsCollector()
	declaredOperation := make(map[string]bool, len(flatOperations))
	for _, flatOp := range flatOperations {
		if declaredOperation[flatOp.id] {
			errs.AddIfNotNil(fmt.Errorf("multiple handlers found for operation id %q", flatOp.id))
			continue
		}
		declaredOperation[flatOp.id] = true
		_, found := oa.spec.findSpecOperationByID(flatOp.id)
		if !found {
			errs.AddIfNotNil(fmt.Errorf("handler recieved for non exising operation id %q is spec", flatOp.id))
		}
		chainResponses := utils.Map(append(flatOp.handlers, flatOp.handler), func(handler handler) handlerResponses {
			return handler.responses
		})
		errs.AddIfNotNil(flatOp.validateOperationTypes(*oa, chainResponses))
		if errs.ErrorOrNil() != nil {
			continue
		}
	}
	for _, pathItem := range oa.spec.Paths {
		for _, specOp := range pathItem.Operations() {
			if !declaredOperation[specOp.OperationID] {
				errs.AddIfNotNil(fmt.Errorf("missing handler for operation id %q", specOp.OperationID))
			}
		}
	}
	return errs.ErrorOrNil()
}

// flattenOperations takes a group with separate operations, handlers, and nested groups and flatten them into a flat
// Operation slice that include for each Operation its own NewOperationHandler, attached handlers, and attached group handlers.
func flattenOperations(g group) []operation {
	groupsOperations := utils.ConcatSlices[operation](utils.Map(g.groups, flattenOperations)...)
	return utils.Map(append(g.operations, groupsOperations...), func(op operation) operation {
		return operation{
			id:       op.id,
			handler:  op.handler,
			handlers: utils.ConcatSlices(g.handlers, op.handlers),
		}
	})
}

func chainHandlers(oa openapi, handlers ...handler) (head groupHandlerFunc[any]) {
	next := func(HandlerContext) (Response[any], error) { return Response[any]{}, nil }
	for i := len(handlers) - 1; i >= 0; i-- {
		next = handlers[i].handlerFunc.handler(oa, next)
	}
	return next
}

func asHttpRouterHandler(oa openapi, head groupHandlerFunc[any]) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if oa.options.RecoverOnPanic {
			defer defaultRecoverBehaviour(writer)
		}
		response, err := head(HandlerContext{Request: request, Params: params})
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		for header, values := range response.Headers {
			for _, value := range values {
				writer.Header().Add(header, value)
			}
		}
		writer.WriteHeader(response.Status)
		writer.Write(response.Bytes)
	}
}

func defaultRecoverBehaviour(writer http.ResponseWriter) {
	if r := recover(); r != nil {
		writer.WriteHeader(500)
		log.Printf("[ERROR] recovered from panic. %v. respond with status 500\n", r)
	}
}
