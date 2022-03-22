# restcontroller

A Gin REST Controller created from an OpenAPI spec.

## Roadmap

- [ ] Shift to a design-first API approaches for building Gin REST Routers from a spec and operationId to Controller map.
  
  This approach retire the current implementation of code-first spec extraction.
  
  The original goal was to create specs in both design-first and code-first approaches and then check for compatability between them using an OpenAPI diff tool. 
 
  The new approach planed, is taking a design-first spec and an operationId to Controller map and using them to init a Gin router with builtin validations based on the spec.  

  The plan is to get something that follows this pattern:
  
  ```go
  package main
  
  import (
      "github.com/gin-gonic/gin"
      "github.com/piiano/restcontroller/example"
      "github.com/piiano/restcontroller/restcontroller"
  )
  
  func main() {
      options := restcontroller.Options{ /**commonErrors, middlewares, configurations, etc.*/ }
      spec := restcontroller.LoadSpec("...")
      controllersMap := map[string]any{
          "greetOperationId": example.GreetController,
      }
      router := spec.GinRouter(controllersMap, options)
      ctx := gin.New()
      ctx.Use(router)
      ctx.Run()
  }
  ```
- [ ] Runtime validation for request body & params based on the OpenAPI spec.
  
  We currently use `github.com/getkin/kin-openapi` that from a quick look seems to have only basic validations support. 

