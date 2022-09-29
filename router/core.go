package router

import (
	"log"
	"net/http"
	"regexp"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"

	"github.com/piiano/cellotape/router/utils"
)

func createMainRouterHandler(oa *openapi) (http.Handler, error) {
	flatOperations := flattenOperations(oa.group)
	if err := validateOpenAPIRouter(oa, flatOperations); err != nil {
		return nil, err
	}
	router := httprouter.New()
	router.HandleMethodNotAllowed = false
	////router.PanicHandler = nil
	//router.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {
	//	log.Println("http-router handler")
	//}
	logger := oa.logger()
	pathParamsMatcher := regexp.MustCompile(`\{([^/}]*)}`)

	specOperations := oa.spec.Operations()
	for _, flatOp := range flatOperations {
		specOp := specOperations[flatOp.id]
		path := pathParamsMatcher.ReplaceAllString(specOp.Path, ":$1")
		chainHead := chainHandlers(*oa, append(flatOp.handlers, flatOp.handler)...)
		httpRouterHandler := asHttpRouterHandler(*oa, specOp, chainHead)
		router.Handle(specOp.Method, path, httpRouterHandler)
		logger.Infof("register handler for operation %q - %s %s", flatOp.id, specOp.Method, specOp.Path)
	}
	return router, nil
}

func (oa *openapi) logger() utils.Logger {
	return utils.NewLoggerWithLevel(oa.options.LogOutput, oa.options.LogLevel)
}

// flattenOperations takes a group with separate operations, handlers, and nested groups and flatten them into a flat
// OperationHandler slice that include for each OperationHandler its own NewOperation, attached handlers, and attached group handlers.
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

func chainHandlers(oa openapi, handlers ...handler) (head BoundHandlerFunc) {
	var next BoundHandlerFunc
	for i := len(handlers) - 1; i >= 0; i-- {
		next = handlers[i].handlerFunc.handlerFactory(oa, next)
	}
	return next
}

func asHttpRouterHandler(oa openapi, specOp SpecOperation, head BoundHandlerFunc) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if oa.options.RecoverOnPanic {
			defer defaultRecoverBehaviour(writer)
		}
		ctx := Context{
			Operation:   specOp,
			Writer:      writer,
			Request:     request,
			Params:      &params,
			RawResponse: &RawResponse{Status: 0},
		}
		_, err := head(ctx)
		if err != nil || ctx.RawResponse.Status == 0 {
			log.Println("unhandled error")
			log.Println(err)
			writer.WriteHeader(500)
			return
		}
	}
}

func defaultRecoverBehaviour(writer http.ResponseWriter) {
	if r := recover(); r != nil {
		writer.WriteHeader(500)
		log.Printf("[Error] recovered from panic. %v. respond with status 500\n", r)
		debug.PrintStack()
	}
}
