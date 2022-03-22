# restcontroller

A Gin REST Controller created from an OpenAPI spec.

## Roadmap

Add support for design-first API approaches for building Gin REST Routers using the following pattern:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/piiano/restcontroller/example"
	"github.com/piiano/restcontroller/restcontroller"
)

func main() {
    spec := restcontroller.LoadSpec("...")
	controllersMap := map[string]any{
		"greetOperationId": example.GreetController,
	}
    options := restcontroller.Options{ /**commonErrors, middlewares, configurations, etc.*/ }
	ctx := gin.New()
	ctx.Use(spec.GinRouter(controllersMap, options))
	ctx.Run()
}
```
