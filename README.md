![Cellotape mascot](./cellotape-gopher.png)

# Cellotape - Beta - OpenAPI Router for Go

![98.1%](https://badgen.net/badge/coverage/98.1%25/green?icon=github)

Cellotape requires Go 1.18 or above.

> **🚧 Cellotape is in Beta 🚧**
> 
> Please note that this is a beta version and the API may change.

A type safe approach to HTTP routing with OpenAPI in Go.
We aim to simplify the way REST APIs are developed with OpenAPI.
This project allow you to do it in a **design-first** approach.

Load an OpenAPI spec and use it as a router to call relevant handlers in your code.
The handler signatures are validated with your OpenAPI spec to verify your code is implementing
your design correctly.

- [Concepts](#concepts)
  - [Included features](#included-features)
  - [What this project isn't doing](#what-this-project-isnt-doing)
- [Get started](#get-started)
- [Loading OpenAPI spec (`router.OpenAPISpec`)](#loading-openapi-spec-routeropenapispec)
  - [Using FS embedding (recommended)](#using-fs-embedding-recommended)
  - [Read file at runtime](#read-file-at-runtime)
  - [Define programmatically in Go](#define-programmatically-in-go)
- [Init new HTTP router (`router.OpenAPIRouter`)](#init-new-http-router-routeropenapirouter)
  - [Create from spec with default options](#create-from-spec-with-default-options)
  - [Create from spec with custom options](#create-from-spec-with-custom-options)
- [Add Operation Implementation - `router.OpenAPIRouter.WithOperation`](#add-operation-implementation---routeropenapirouterwithoperation)
- [Define Operation Handler - `router.NewHandler`](#define-operation-handler---routernewhandler)
  - [Request Context `router.Context`](#request-context-routercontext)
  - [Request - `router.Request[B, P, Q]`](#request---routerrequestb-p-q)
  - [Responses - `router.Response[R]`](#responses---routerresponser)
- [Examples](#examples)
  - [Hello World API Example](#hello-world-api-example)
  - [TODO List API Example](#todo-list-api-example)
- [Roadmap](#roadmap)

## Concepts

The OpenAPI spec is a great way to describe accurately your API for your users.

Developing an API in Go, it's often challenging making sure your OpenAPI spec and your Go implementation are in sync.  

Cellotape allow you to develop your APIs in Go in a way that helps you catch inconsistencies between your spec and your code.  

The Go ecosystem provides many packages and frameworks for building HTTP servers for REST APIs.

Most of them rely on concepts from the builtin [net/http](https://pkg.go.dev/net/http) package for defining HTTP request handlers.

The issue cellotape come to solve is the lack of type information when working with such handlers.

Instead of having the same untyped signature for all handlers, Cellotape handlers use generics to define for each handler the types of its body, path params, query params and possible responses.

This extra type information allow us to define handlers which can be validated at runtime with an OpenAPI specification.

### Included features

- Load, parse and validate an OpenAPI spec (`router.OpenAPISpec`).

- Init an HTTP router driven from an OpenAPI spec (`router.OpenAPIRouter`).

- Provide SDK for defining strongly typed handlers mapped to spec operations.

- Verify handler signature compatability with spec operations.

  This enforces during initialization that the handlers correctly implement the spec.

  Validated components includes:

  - **Request Body Schema**

  - **Path Parameters**

  - **Query Parameters**

  - **Responses**

- Support for middleware chains and group mechanism that allow applying middlewares to specific operations or specific groups.

- Compatability with the `http.Handler` interface for both the router itself and the middlewares to allow easy integration of the router in any popular framework.

- Support for custom content types to align with content types defined in the spec.

  This can be done by implementing the `router.ContentType` interface

- Support for customization of validation behaviour and other configuration using `router.Options`.

  For the full documentation of the available options you can check the documentation on the [`router.Options`](./router/options.go) struct.


### What this project isn't doing

- It's not a code generator.

  Code generators have their issues, with minimal control over the generated code, and the difficulty to edit generated code that can overridden whenever the spec changes.

- It's not using comments on your code. Some approaches might suggest you comment your code in a way that produce an OpenAPI spec that can be compared to the spec you design.

  The issue with this approach is that the spec, code and comments can get out of sync too easily.

  Especially when HTTP handlers in go provide an opaque interface with no type information.

Our approach using strongly typed handlers can help you create an implementation for your API that never gets out of sync and are easily maintained.

## Get started

Add cellotape to your project using `go get`:

```bash
go get github.com/piiano/cellotape@latest
```

Add an `openapi.yml` describing your API to your project.
For example, the following spec defines a single `greet` API:

```yaml
openapi: 3.0.3
info:
  title: "Hello World Example API"
  version: 1.0.0
servers:
  - url: 'https'
paths:
  /greet:
    post:
      operationId: greet
      summary: Greet the caller
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
      parameters:
        - in: query
          name: greetTemplate
          description: An optional template to customize the greeting.
          schema:
            type: string
      responses:
        200:
          description: successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  greeting:
                    type: string
```

Add the following go code to your project to implement the OpenAPI spec:

```go
package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	
	"github.com/piiano/cellotape/router"
)

//go:embed openapi.yaml
var specData []byte

func main() {
	if err := handleMain(); err != nil {
		log.Fatal(err)
	}
}

func handleMain() error {
	handler, err := initHandler()
	if err != nil {
		return err
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		return err
	}
	return nil
}

func initHandler() (router.OpenAPIRouter, error) {
	spec, err := router.NewSpecFromData(specData)
	if err != nil {
		return nil, err
	}
	return router.NewOpenAPIRouter(spec).
		WithOperation("greet", router.NewHandler(greetHandler)).
        AsHandler()
}


type body struct {
  Name string `json:"name"`
}
type queryParams struct {
  GreetTemplate string `form:"greetTemplate"`
}
type responses struct {
  OK ok `status:"200"`
}
type ok struct {
  Greeting string `json:"greeting"`
}

func greetHandler(_ router.Context, request router.Request[body, router.Nil, queryParams]) (router.Response[responses], error) {
  var greeting string
  if request.QueryParams.GreetTemplate == "" {
    greeting = fmt.Sprintf("Hello %s!", request.Body.Name)
  } else {
    greeting = fmt.Sprintf(request.QueryParams.GreetTemplate, request.Body.Name)
  }

  return router.SendOKJSON(responses{OK: ok{Greeting: greeting}}), nil
}
```

In this code, we use `router.NewSpecFromData` to load the OpenAPI spec.

We initialize the Cellotape router with `router.NewOpenAPIRouter`.

Calling then to `WithOperation("greet", r.NewHandler(greetHandler))` we tell Cellotape to use `greetHandler` as the implementation of the `greet` operation.

The `greetHandler` function defines typed parameters and typed response.

When calling `AsHandler()`, cellotape check the request and response types with reflection and checks their compatability with the provided spec.

If we were implementing the spec incorrectly we would receive an error at this point during the server initialization.

Finally, we can add a simple test to verify our server is always compatible with the spec:

```go
package main

import (
  "testing"
  
  "github.com/stretchr/testify/assert"
)

func TestServerCompatabilityWithOpenAPI(t *testing.T) {
  _, err := initHandler()
  require.NoError(t, err)
}
```

> **Tip**
> 
> Try changing the spec in a way that is incompatible with the API and run the server again.
> 
> Cellotape is reporting to you the incompatibility and allow you to make sure your implementation and spec are in sync. 

## Loading OpenAPI spec (`router.OpenAPISpec`)

To init a new Cellotape HTTP router you must first load an OpenAPI spec.

The OpenAPI spec defines URL paths and HTTP methods for different operations.

Cellotape routes HTTP call on these paths and methods for the relevant handler implementation of each operation.

Cellotape also validates the handler is implementing correctly the request and response defined in the spec.

There are few ways you can load the OpenAPI spec to be used with Cellotape.

### Using FS embedding (recommended)

Using [Go embedding](https://pkg.go.dev/embed) you can embed your OpenAPI YAML or JSON file directly with the compiled binary of your app. 

For example embedding an `openapi.yaml` as follows. 
```go
//go:embed openapi.yaml
var specData []byte
```

This allows you to init then the spec from its bytes with `router.NewSpecFromData` as follows:

```go
spec, err := router.NewSpecFromData(specData)
```

`router.NewSpecFromData` returns the loaded spec object and an error if failed to parse and validate it.

### Read file at runtime

Using `router.NewSpecFromFile` you can read and load the OpenAPI spec JSON or YAML file at runtime. 

```go
// path to openapi.yml file
openapiFilePath := "./openapi.yml"
spec, err := router.NewSpecFromFile(openapiFilePath)
```

This option is less recommended as if the spec file is changed or missing your application may break.

Similar to `router.NewSpecFromData`, this method returns an error if failed to parse and validate the spec but can also error if failed to read the file.

### Define programmatically in Go

Sometimes you might want to define the spec programmatically in your Go code rather than using a YAML or JSON format.

Cellotape uses internally [kin-openapi](github.com/getkin/kin-openapi) to define the OpenAPI model.

You can define with `kin-openapi` the OpenAPI model programmatically and then define an `router.OpenAPISpec` from it.

```go
openapiModel := openapi3.T{
    //...
}
if err := openapiModel.Validate(); err != nil {
    // potentially validate the spec before using it.
}

spec := router.OpenAPISpec(openapiModel)
```

## Init new HTTP router (`router.OpenAPIRouter`)

The `router.OpenAPIRouter` is the main builder block of Cellotape. 

With it, you define the handlers and middlewares of yor application.


### Create from spec with default options

To init the router from the spec with the default options you can use the `r.NewOpenAPIRouter` function.

```go
spec, err := router.NewSpecFromData(specData)
if err != nil {
    // handle err
}

openapiRouter := router.NewOpenAPIRouter(spec)
```

### Create from spec with custom options

Sometimes the default validations can be overwhelming during development, and you might want to adjust the behavior from error to warning or completely ignore certain validations. 

For that you can initialize the router with custom options using `r.NewOpenAPIRouterWithOptions`.

```go
spec, err := router.NewSpecFromData(specData)
if err != nil {
    // handle err
}

options := router.Options{
    // ...
}

openapiRouter := router.NewOpenAPIRouterWithOptions(spec, options)
```

## Add Operation Implementation - `router.OpenAPIRouter.WithOperation`

To implement API operations defined in the OpenAPI spec, we use the `router.OpenAPIRouter.WithOperation` method to add operation handlers to our router.

```go
openapiRouter := router.NewOpenAPIRouter(spec)

openapiRouter.WithOperation("operation-id", router.NewHandler(...))
```

The operation ID provided must match the operation ID defined in the spec.

The router will route HTTP requests that matches the path template and HTTP method defined in the spec for that operation.

The router also validates the provided handler request and response types are compatible with those defined in the spec. 

## Define Operation Handler - `router.NewHandler`

To create a new handler we often use the `router.NewHandler` to create a handler from a typed handler function `router.HandlerFunc[B, P, Q, R]`.

```go
type Responses struct {
	OK string `status:"200"`
}

var handler = router.NewHandler(func (
	c router.Context,
	request Request[router.Nil, router.Nil, router.Nil],
) (router.Response[Responses], error) {
	return router.SendOKText(Responses{OK: "hello world!"})
})
```

Notice, the handler function defines the types of the request and possible responses.  

This is how Cellotape is capable of reading the types with reflection and check their compatibility with the spec. 

### Request Context `router.Context`

The first parameter of a `router.HandlerFunc[B, P, Q, R]` function is a `router.Context`.

The context includes the native HTTP `http.ResponseWriter` and `*http.Request`.

When in a middleware you can use `router.Context.Next` to call the next handler in the chain.

It includes a `router.SpecOperation` with the operation definition read from the spec.

And after a response was returned by a handler in the chain it contains the raw response using `*router.RawResponse`.

### Request - `router.Request[B, P, Q]`

The second parameter of a `router.HandlerFunc[B, P, Q, R]` function is a `router.Request[B, P, Q]`.

The `router.Request[B, P, Q]` defines 3 generic arguments.

- `B` - The type of the request body.
- `P` - The struct type of the request path parameters.
- `Q` - The struct type of the request query parameters.

These types are reflected as parameters of the request so you can use them in the handler function.  

#### Request Body - <code>router.Request[<strong>B</strong>, P, Q]</code>

The first generic argument (`B`) of `router.Request[B, P, Q]` represents the request body.

For example, having an API that expect to receive a simple JSON object with a `name` property can be defined in the spec as follows: 

```yaml
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
```

To represent this request body in the Go implementation we can define a struct and use the `json` tag to define its properties.

```go
type GreetBody struct {
	Name string `json:"name"`
}
```

Defining this struct with the `router.Request[B, P, Q]` param in the handler will look as follows:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[GreetBody, router.Nil, router.Nil],
) (router.Response[Responses], error) {
    return router.SendOKText(Responses{OK: fmt.Sprintf("hello %s!", request.Body.Name)})
})
```

Notice how we can access the value of `request.Body.Name` with cellotape binding it for you from the HTTP request. 

This type is validated for compatability with the schema defined in the spec `requestBody` property.

#### Path Params - <code>router.Request[B, <strong>P</strong>, Q]</code>

The second generic argument (`P`) of `router.Request[B, P, Q]` represents the request path parameters.

For example, having an API that expect to receive a path parm `name` in the request URL property can be defined in the spec as follows:

```yaml
paths:
  /greet/{name}:
    post:
      operationId: greet
      summary: Greet the caller
      parameters:
        - in: path
          required: true
          name: name
          schema:
            type: string
```

To represent this path parameter in the Go implementation we can define a struct.

Each struct field is representing a path parm. You can use the `uri` tag to define the param, name as stated in the spec.

```go
type GreetPathParams struct {
	Name string `uri:"name"`
}
```

Defining this struct with the `router.Request[B, P, Q]` param in the handler will look as follows:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[router.Nil, GreetPathParams, router.Nil],
) (router.Response[Responses], error) {
    return router.SendOKText(Responses{OK: fmt.Sprintf("hello %s!", request.PathParams.Name)})
})
```

Notice how we can access the value of `request.PathParams.Name` with cellotape binding it for you from the HTTP request.

This type is validated for compatability with the schema defined in each parameter of the spec `parameters` that have `in: path`.

#### Query Params - <code>router.Request[B, P, <strong>Q</strong>]</code>

The third generic argument (`Q`) of `router.Request[B, P, Q]` represents the request query parameters.

For example, having an API that expect to receive a query parm `name` in the request URL property can be defined in the spec as follows:

```yaml
paths:
  /greet/{name}:
    post:
      operationId: greet
      summary: Greet the caller
      parameters:
        - in: query
          required: true
          name: name
          schema:
            type: string
```

To represent this query parameter in the Go implementation we can define a struct.

Each struct field is representing a query parm. You can use the `form` tag to define the param, name as stated in the spec.

```go
type GreetQueryParams struct {
	Name string `form:"name"`
}
```

Defining this struct with the `router.Request[B, P, Q]` param in the handler will look as follows:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[router.Nil, router.Nil, GreetQueryParams],
) (router.Response[Responses], error) {
    return router.SendOKText(Responses{OK: fmt.Sprintf("hello %s!", request.QueryParams.Name)})
})
```

Notice how we can access the value of `request.QueryParams.Name` with cellotape binding it for you from the HTTP request.

This type is validated for compatability with the schema defined in each parameter of the spec `parameters` that have `in: query`.

### Responses - `router.Response[R]`

The first return type of the handler function is a `router.Response[R]`.

This generic argument of this type should be a struct with which each field represent a possible response and is tagged with the `status` tag. 

For example, having an API that can return a 200 and 400 possible responses in a simple JSON can be defined in the spec as follows:

```yaml
      responses:
        200:
          description: Greeting
          content:
            application/json:
              schema:
                type: object
                properties:
                  greeting:
                    type: string
        400:
          description: Bad Request
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: string
```

To represent these responses in the Go implementation we define a struct and use the `status` tag to define each possible response.

```go
type GreetOKResponse struct {
    Greeting string `json:"greeting"`
}
type GreetBadRequestResponse struct {
    Error string `json:"error"`
}
type GreetResponses struct {
	OK         GreetOKResponse         `status:"200"`
    BadRequest GreetBadRequestResponse `status:"400"`
}
```

Defining this struct with the `router.Response[R]` return type in the handler will look as follows:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[GreetBody, router.Nil, router.Nil],
) (router.Response[GreetResponses], error) {
	if request.Body.Name == "" {
        return router.SendJSON(GreetResponses{
            BadRequest: GreetBadRequestResponse{
                Error: "name property is empty",
            },
        }).Status(400)
    }
	return router.SendOKJSON(GreetResponses{
		OK: GreetOKResponse{
            Greeting: fmt.Sprintf("hello %s!", request.Body.Name),
		},
	})
})
```

Notice how the multiple fields of the responses struct defined allow you to send the response you want based on your implementation logic in a way that cover all the responses defined the spec.

Each response type is validated for compatability with the schema defined in each response of the spec `responses`.

## Examples

You can learn more about how to use this package by reviewing the following examples.

### Hello World API Example

You can check the [Hello World Example](./examples/hello_world_example) to see how it works.
We use the following [openapi.yaml](./examples/hello_world_example/openapi.yaml)
to init the server and map to the relevant handlers.

### TODO List API Example

You can check the [TODO List API Example](./examples/todo_list_app_example) to see how it works with a more realistic usage.
We use the following [openapi.yaml](./examples/todo_list_app_example/openapi.yaml)
to init the server and map to the relevant handlers.

## Roadmap

- [ ] Improve documentation with more usage details.

- [ ] Runtime validation for request body & params based on the OpenAPI spec.

  Multiple options are considered:

  - https://github.com/go-playground/validator

  - https://github.com/xeipuuv/gojsonschema

- [ ] Add support for serialization of OpenAPI parameters `style`, and `explode`, and `allowReserved` properties.

- [ ] Add support for OpenAPI Header Params.

- [ ] Add support for OpenAPI Cookie Params.
