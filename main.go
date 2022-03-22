package main

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"github.com/piiano/restcontroller/example"
	"github.com/piiano/restcontroller/restcontroller"
)

func main() {
	operationOptions := &restcontroller.OperationOptions{Errors: openapi3.NewResponses()}
	gin.SetMode(gin.ReleaseMode)
	ctx := gin.New()
	controllers := map[string]restcontroller.Controller{
		"greet": example.GreetController,
	}
	for operationID, controller := range controllers {
		group := ctx.Group("/greet")             // TODO extract path from spec
		group.POST("/", controller.GinHandler()) // TODO extract method from spec
		operation, err := controller.OpenAPIOperation(operationID, operationOptions)
		if err != nil {
			panic(err)
		}
		operationJson, err := yaml.Marshal(operation)
		fmt.Println(string(operationJson))
	}
}
