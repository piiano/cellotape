# restcontroller

A DRY approach to REST Controller with OpenAPI in Golang.
We aim to simplify the way REST APIs are developed with OpenAPI.
This project allow you to do it by either a design-first or code-first approach.

- **Design First** - Load an OpenAPI spec and use it as a router to call relevant controllers in your code.
- **Code First** - Write your code and use reflection on the controllers signatures to produce an OpenAPI spec.

## What this project isn't doing

- It's not using code generation to support design-first approach.
  
  Code generators have their issues, code can get outdated, changes to generated code are hard to maintain because they might be overridden whenever the spec changes.
  
  We load the spec on server initialization and init a router based on it.

- It's not producing all the extra documentation properties in the spec such as descriptions, examples, etc. 
  
  The purpose of the spec produced in the code-first approach is to only describe the API signatures for compatability checks, API Gateways and for client generation.

## Examples

### Hello World Example

You can check the [Hello World Example](./examples/hello_world_example) to see how it works.
We use the following [openapi.yaml](./examples/hello_world_example/openapi.yaml)
([see with UI](https://editor.swagger.io?url=https://raw.githubusercontent.com/piiano/restcontroller/main/example/hello-world-openapi.yaml?token%3DGHSAT0AAAAAABSHBLZSQVEWSF62YUJLYSK6YSDMK5A))
to init the server and map to the relevant handlers. 

### TODO List API Example

You can check the [TODO List API Example](./examples/todo_list_app_example) to see how it works with a more realistic usage.
We use the following [openapi.yaml](./examples/todo_list_app_example/openapi.yaml)
([see with UI](https://editor.swagger.io?url=https://raw.githubusercontent.com/piiano/restcontroller/main/example/hello-world-openapi.yaml?token%3DGHSAT0AAAAAABSHBLZSQVEWSF62YUJLYSK6YSDMK5A))
to init the server and map to the relevant handlers.


## Roadmap

- [x] Shift to a design-first API approaches for building Gin REST Routers from a spec and operationId to Controller map.
    
  The original goal was to create specs in both design-first and code-first approaches and then check for compatability between them using an OpenAPI diff tool. 
  
  The new approach, is taking a design-first spec and an operationId to Controller map and using them to init a Gin router with builtin validations based on the spec.
  
- [x] Reflection validation during initialization for request body & params based on the OpenAPI spec.
  
- [x] Add support for middlewares - Require some designing of how this API will look like.
  
- [x] Add support for custom binders - support for more content-types and ways to bind them to the relevant params.
  
- [ ] Reflection validation during initialization for responses and possible errors based on the OpenAPI spec.
  
- [ ] Runtime validation for request body & params based on the OpenAPI spec.

  We currently use [getkin/kin-openapi](https://github.com/getkin/kin-openapi) that from a quick look seems to have only basic validations support.

  We might want to use [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema) that seems to have a good support for JSON schema validation (which is part of OpenAPI).

- [ ] We might want to add support for [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter).
  
- [ ] We might want to create a version that validates signature at compile time using the AST
  
  The advantages of this approach is it can capture all return statements of a function and extract their exact type.
  
  It can also result with a final cleaner API for the user.

  
