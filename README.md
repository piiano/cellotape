# OpenAPI Router

A safe DRY approach to HTTP routing with OpenAPI in Golang.
We aim to simplify the way REST APIs are developed with OpenAPI.
This project allow you to do it in a **design-first** approach.

Load an OpenAPI spec and use it as a router to call relevant handlers in your code.
The router handler signatures are validated with your OpenAPI spec and you can rest assured your code is implementing
your design correctly with no gaps between your code and your spec. 

## What this project isn't doing

- It's not using code generation to support design-first approach.
  
  Code generators have their issues, code can get outdated, changes to generated code are hard to maintain because they might be overridden whenever the spec changes.

- It's not using comments on your code. Some approaches might suggest you comment your code in a way that produce an OpenAPI spec that can be compared to the spec you design.
  
  The issue with this approach is that same as code the spec, code and comments can get out of sync too.

  Especially when HTTP handlers in go provide an opaque interface with no type information.

Our approach using strongly typed handlers can help you create an implementation for your API that never gets out of sync and are easily maintained. 

## Features

- [x] Load, parse and validate an OpenAPI spec (`router.OpenAPISpec`).

- [x] Init an HTTP router driven from an OpenAPI spec (`router.OpenAPI`).

- [x] Provide SDK for defining strongly typed handler mapped to spec operations.

- [x] Verify the handler signature types are compatible with those defined in the spec.
  
  This enforces during initialization that the handlers correctly implement the spec. 
 
  Validated components includes:
  
  - **Request Body Schema**
    
  - **Path Parameters**
    
  - **Query Parameters**
    
  - **Responses**
  
- [x] Support for middleware chains and group mechanism that allow applying middlewares to specific operations or specific groups
  
- [x] Compatability with the `http.Handler` interface for both the router itself and the middlewares to allow easy integration of the router in any popular framework

- [x] Support for custom content types to align with content types defined in the spec. 
 
  This can be done by implementing the `router.ContentType` interface

## Examples

You can learn more about how to use this package by reviewing the following examples.  

### Hello World API Example

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
 
- [ ] Runtime validation for request body & params based on the OpenAPI spec.
  
  We might want to look at few options for that:
  
  - https://github.com/go-playground/validator
    
  - https://github.com/xeipuuv/gojsonschema
  
- [ ] Add support for better customization using an `router.Options` parameter 
  
- [ ] Add support for additional OpenAPI features such as Header Params, Cookie Params, etc.

- [ ] We might want to replace the internal implementation with [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter).

- [ ] Improve test coverage for the router package