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

	group := ctx.Group("/greet")
	group.POST("/", restcontroller.GinHandler(example.GreetController))
	operation, err := restcontroller.Operation("Greet", example.GreetController, operationOptions)

	if err != nil {
		panic(err)
	}
	operationJson, err := yaml.Marshal(operation)
	fmt.Println(string(operationJson))
}
