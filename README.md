![Cellotape mascot](./cellotape-gopher.png)

# Cellotape - Beta - OpenAPI Router for Go

> ### 🚧 Cellotape is in Beta
> Please note that this is a beta version and the API may change.

A type safe approach to HTTP routing with OpenAPI in Golang.
We aim to simplify the way REST APIs are developed with OpenAPI.
This project allow you to do it in a **design-first** approach.

Load an OpenAPI spec and use it as a router to call relevant handlers in your code.
The handler signatures are validated with your OpenAPI spec to verify your code is implementing
your design correctly.

- [Concepts](#concepts)
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
- [Middlewares - `router.OpenAPIRouter.Use`](#middlewares---routeropenapirouteruse)
- [Error handling - `router.ErrorHandler`](#error-handling---routererrorhandler)
- [Handler groups - `router.OpenAPIRouter.WithGroup`](#handler-groups---routeropenapirouterwithgroup)
- [Content types - `router.OpenAPIRouter.WithContentType`](#content-types---routeropenapirouterwithcontenttype)
- [Options - `router.Options`](#options---routeroptions)

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
go get github.com/piiano/cellotape@v1
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

In your go code, load the spec, init the router and provide an implementation to the `greet` API:

```go
package main

import (
	_ "embed"
	"fmt"
	"github.com/piiano/cellotape/router"
	"log"
	"net/http"
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

OpenAPI router will load the spec and init the router.

Calling then to `WithOperation("greet", r.NewHandler(greetHandler))` we tell Cellotape to use `greetHandler` as the implementation of the `greet` operation.

The `greetHandler` function defines typed parameters and typed response.

When calling `AsHandler()`, cellotape check the request and response types with reflection and check their compatability with the provided spec.

If we were implementing the spec incorrectly we would receive an error at this point during the server initialization.

Finally, we can add a simple test to verify our server is always compatible with the spec:

```go
package main

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestServerCompatabilityWithOpenAPI(t *testing.T) {
  _, err := initHandler()
  require.NoError(t, err)
}
```

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

## Adding Operation Implementation - `router.OpenAPIRouter.WithOperation`

## Define Operation Handler - `router.NewHandler`

### Request Context `router.Context`

### Request - `router.Request[B, P, Q]`

#### Request Body - <code>router.Request[<strong>B</strong>, P, Q]</code>

#### Path Params - <code>router.Request[B, <strong>P</strong>, Q]</code>

#### Query Params - <code>router.Request[B, P, <strong>Q</strong>]</code>

### Responses - `router.Response[R]`

#### Defining possible responses

#### Sending a response `router.Send`

##### Response Status Code - `router.Response.Status`

##### Response Content Type - `router.Response.ContentType`

##### Response Content Headers - `router.Response.SetHeader`

## Middlewares - `router.OpenAPIRouter.Use`

## Error handling - `router.ErrorHandler`

## Handler groups - `router.OpenAPIRouter.WithGroup`

## Content types - `router.OpenAPIRouter.WithContentType`

## Options - `router.Options`


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

- [ ] Runtime validation for request body & params based on the OpenAPI spec.

  We might want to look at few options for that:

  - https://github.com/go-playground/validator

  - https://github.com/xeipuuv/gojsonschema

- [ ] Add support for serialization of OpenAPI parameters `style`, and `explode`, and `allowReserved` properties.

- [ ] Add support for OpenAPI Header Params.

- [ ] Add support for OpenAPI Cookie Params.
