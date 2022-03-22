package restcontroller

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func InitRouter(engine *gin.Engine, specPath string, controllers map[string]Controller) error {
	spec, err := openapi3.NewLoader().LoadFromFile(specPath)
	if err != nil {
		return err
	}
	for path, operations := range spec.Paths {
		for method, operation := range operations.Operations() {
			controller := controllers[operation.OperationID]
			if controller == nil {
				return errors.New(fmt.Sprintf("Couldn't find a controler with operation ID %q in controllers map.", operation.OperationID))
			}
			var methodRouterFn func(string, ...gin.HandlerFunc) gin.IRoutes
			switch method {
			case http.MethodGet:
				methodRouterFn = engine.GET
			case http.MethodPost:
				methodRouterFn = engine.POST
			case http.MethodPut:
				methodRouterFn = engine.PUT
			case http.MethodPatch:
				methodRouterFn = engine.PATCH
			case http.MethodDelete:
				methodRouterFn = engine.DELETE
			case http.MethodHead:
				methodRouterFn = engine.HEAD
			case http.MethodOptions:
				methodRouterFn = engine.OPTIONS
			default:
				return errors.New(fmt.Sprintf("Method %q not supported.", method))
			}
			pathParamsMatcher := regexp.MustCompile(`\{([^/]*)}`)
			formattedPath := pathParamsMatcher.ReplaceAllString(path, ":$1")
			fmt.Println(formattedPath)
			methodRouterFn(formattedPath, controller.GinHandler())
		}
	}
	return nil
}

func (fn ControllerFn[B, P, Q, R]) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params Params[B, P, Q]
		if err := c.BindUri(&params.Path); err != nil {
			c.JSON(400, err.Error())
			return
		}
		if err := c.BindQuery(&params.Query); err != nil {
			c.JSON(400, err.Error())
			return
		}
		if err := c.Bind(&params.Body); err != nil {
			c.JSON(400, err.Error())
			return
		}
		params.Headers = c.Request.Header
		resp, err := fn(params)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
