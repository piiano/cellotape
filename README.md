![Cellotape mascot](./Cellotape-gopher.png)

# Cellotape - Beta - OpenAPI Router for Go

![98.1%](https://badgen.net/badge/coverage/98.1%25/green?icon=github)

Cellotape requires Go 1.18 or above.

> **🚧 Cellotape is in Beta 🚧**
> 
> Please note that this is a beta version, and the API may change.

A type-safe approach to HTTP routing with OpenAPI in Go. The project aims to 
simplify the development of REST APIs with OpenAPI by enabling a 
**design-first** approach.

Cellotape loads an OpenAPI spec and uses it as a router to call relevant handlers 
in your code. The handler signatures are validated with your OpenAPI spec to verify
 your code is implementing the design correctly.

- [Concepts](#concepts)
  - [Included features](#included-features)
  - [What this project isn't doing](#what-this-project-isnt-doing)
- [Get started](#get-started)
- [Loading an OpenAPI spec (`router.OpenAPISpec`)](#loading-openapi-spec-routeropenapispec)
  - [Using FS embedding (recommended)](#using-fs-embedding-recommended)
  - [Read file at runtime](#read-file-at-runtime)
  - [Define programmatically in Go](#define-programmatically-in-go)
- [Initialize new HTTP router (`router.OpenAPIRouter`)](#initialize-new-http-router-routeropenapirouter)
  - [Create from spec with default options](#create-from-spec-with-default-options)
  - [Create from spec with custom options](#create-from-spec-with-custom-options)
- [Add operation implementation - `router.OpenAPIRouter.WithOperation`](#add-operation-implementation---routeropenapirouterwithoperation)
- [Define operation handler - `router.NewHandler`](#define-operation-handler---routernewhandler)
  - [Request Context `router.Context`](#request-context-routercontext)
  - [Request - `router.Request[B, P, Q]`](#request---routerrequestb-p-q)
  - [Responses - `router.Response[R]`](#responses---routerresponser)
- [Examples](#examples)
  - [Hello world API example](#hello-world-api-example)
  - [TODO list API example](#todo-list-api-example)
- [Roadmap](#roadmap)

## Concepts

The OpenAPI spec is a great way to describe your API accurately. However, when 
developing an API in Go, it's often a challenge to synchronize the OpenAPI spec 
and Go implementation.

Cellotape enables you to develop APIs in Go in a way that helps catch 
inconsistencies between your spec and code.

The Go ecosystem provides many packages and frameworks for building HTTP servers 
for REST APIs. Most of them rely on concepts from the built-in 
[net/http](https://pkg.Go.dev/net/http) package to define HTTP request handlers.

The issue Cellotape solves is the lack of type information when working with such 
handlers.

Instead of having the same untyped signature for all handlers, Cellotape handlers 
use generics to define for each handler the types of its body, path parameters, 
query parameters, and responses.

This extra type information enables Cellotape to define handlers that are 
validated at runtime with an OpenAPI specification.

### Included features

- Load, parse, and validate an OpenAPI spec (`router.OpenAPISpec`).
- Initialize an HTTP router driven from an OpenAPI spec (`router.OpenAPIRouter`).
- An SDK that defines strongly typed handlers mapped to spec operations.
- Verify handler signature compatibility with spec operations.
  During initialization, this enforces that the handlers correctly implement the 
  spec. The validated components include:
  - Request body schema
  - Path parameters
  - Query parameters
  - Responses
- Support for middleware chain and group mechanisms so that middleware can be 
  applied to operations or groups.
- Compatibility with the `HTTP.Handler` interface for the router and middleware 
  to enable easy integration of the router into any popular framework.
- Support for custom content types to align with content types defined in the spec.
  This can be done by implementing the `router.ContentType` interface.
- Support for customization of validation behavior and other configuration using
  `router.Options`.
  See the documentation on the [`router.Options`](./router/options.go) struct 
  for details of the options.

### What this project isn't doing

- It's not a code generator.
  Code generators have their issues, such as minimal control over the generated 
  code and the overwriting of manual edits whenever the spec changes.

- It's not using comments in your code. 
  With this approach, you use comments in the code to generate an OpenAPI spec 
  (a code-first approach). Then, to implement a design-first approach, you compare
  the generated spec to the desired spec.
  The issue with this approach is that the specification, code, and comments can 
  get out of sync too easily. Particularly when HTTP handlers in Go provide an 
  opaque interface with no type information.

Cellotape use of strongly typed handlers can help create an implementation for 
your API that never gets out of sync and is easily maintained.

## Get started

Add Cellotape to your project using `go get`:

```bash
go get github.com/piiano/cellotape@latest
```

Add an `openapi.yml` describing your API to your project. For example, this spec 
defines a `greet` API:

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

Adding this Go code to your project to implement the OpenAPI spec:

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

This code uses `router.NewSpecFromData` to load the OpenAPI spec.

You initialize the Cellotape router with `router.NewOpenAPIRouter`.

Then, calling `WithOperation("greet", r.NewHandler(greetHandler))` tells Cellotape 
to use `greetHandler` to implement the `greet` operation.

The `greetHandler` function defines typed parameters and typed responses.

When calling `AsHandler()`, Cellotape checks the request and response types with 
reflection and checks their compatibility with the spec.

If you have implemented the spec incorrectly, you receive an error at this point 
during the server initialization.

Finally, you add a simple test to verify the server is always compatible with 
the spec:

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
> Try changing the spec in a way that is incompatible with the API and run the 
> server again.
> 
> Cellotape reports any incompatibility, enabling you to check that your 
> implementation and spec are in sync. 

## Loading OpenAPI spec (`router.OpenAPISpec`)

To initialize a new Cellotape HTTP router you must first load an OpenAPI spec. 

The OpenAPI spec defines URL paths and HTTP methods for various operations. 
Cellotape routes HTTP calls on these paths and methods for the relevant handler 
implementation of each operation. Cellotape also validates that the handler 
correctly implements the request and response defined in the spec.

There are several ways you can load the OpenAPI spec for use with Cellotape.

### Using FS embedding (recommended)

Using [Go embedding](https://pkg.go.dev/embed), you embed the OpenAPI YAML or 
JSON file with the compiled binary of your app. For example, embedding an 
`openapi.yaml` like this:

```go
//go:embed openapi.yaml
var specData []byte
```

This enables you to initialize the spec from its bytes with `router.NewSpecFromData` 
like this:

```go
spec, err := router.NewSpecFromData(specData)
```

`router.NewSpecFromData` returns the loaded spec object or an error if the object 
fails to parse and validate.

### Read file at runtime

Using `router.NewSpecFromFile`, you read and load the OpenAPI spec JSON or 
YAML file at runtime like this:

```go
// path to openapi.yml file
openapiFilePath := "./openapi.yml"
spec, err := router.NewSpecFromFile(openapiFilePath)
```

This option is not recommended; if the spec file is changed or missing, your 
application may break.

Similar to `router.NewSpecFromData`, this method returns an error if the object 
fails to parse and validate but can also error if it fails to read the file.

### Define programmatically in Go

Sometimes you may want to define the spec programmatically in your Go code, rather
than using a YAML or JSON format.

Cellotape uses [kin-openapi](GitHub.com/getkin/kin-openapi) to define the OpenAPI 
model. You can also use `kin-openapi` to create the OpenAPI model programmatically 
and then define a `router.OpenAPISpec` from it like this:

```go
openapiModel := openapi3.T{
    //...
}
if err := openapiModel.Validate(); err != nil {
    // potentially validate the spec before using it.
}

spec := router.OpenAPISpec(openapiModel)
```

## Initialize new HTTP router (`router.OpenAPIRouter`)

The `router.OpenAPIRouter` is the main building block of Cellotape. With it, you 
define the handlers and middleware of your application.

### Create from spec with default options

To initialize the router from the spec with the default options, you use the 
`r.NewOpenAPIRouter` function like this:

```go
spec, err := router.NewSpecFromData(specData)
if err != nil {
    // handle err
}

openapiRouter := router.NewOpenAPIRouter(spec)
```

### Create from spec with custom options

Sometimes the default validations can be overwhelming during development, and you 
may want to adjust the behavior from error to warning or completely ignore certain 
validations. 

For that, you initialize the router with custom options using 
`r.NewOpenAPIRouterWithOptions` like this:

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

To implement API operations defined in the OpenAPI spec, Cellotape uses the 
`router.OpenAPIRouter.WithOperation` method to add operation handlers to the 
router like this:

```go
openapiRouter := router.NewOpenAPIRouter(spec)

openapiRouter.WithOperation("operation-id", router.NewHandler(...))
```

The operation ID provided must match the operation ID defined in the spec.

The router routes HTTP requests that match the path template and HTTP method 
defined in the spec for that operation. It also validates that the handler request 
and response types are compatible with those defined in the spec. 

## Define Operation Handler - `router.NewHandler`

To create a new handler, you use the `router.NewHandler` to create a handler 
from a typed handler function `router.HandlerFunc[B, P, Q, R]` like this:

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

Notice that the handler function defines the types of the request and its responses.

This is how Cellotape is capable of reading the types with reflection and checking 
their compatibility with the spec.

### Request context `router.Context`

The first parameter of a `router.HandlerFunc[B, P, Q, R]` function is a 
`router.Context`. The context includes the native HTTP `http.ResponseWriter` 
and `*http.Request`.

When in middleware, you use `router.Context.Next` to call the next handler in 
the chain. It includes a `router.SpecOperation` with the operation definition 
read from the spec. After a handler in the chain returns a response, it 
contains the raw response using `*router.RawResponse`.

### Request - `router.Request[B, P, Q]`

The second parameter of a `router.HandlerFunc[B, P, Q, R]` function is a 
`router.Request[B, P, Q]`, which defines 3 generic arguments.

- `B` - The type of the request body.
- `P` - The struct type of the request path parameters.
- `Q` - The struct type of the request query parameters.

These types are reflected as parameters of the request so that you can use them in 
the handler function.

#### Request body - <code>router.Request[<strong>B</strong>, P, Q]</code>

The first generic argument (`B`) of `router.Request[B, P, Q]` represents the 
request body.

For example, an API that expects to receive a simple JSON object with a `name` 
property can be defined in the spec like this: 

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

To represent this request body in the Go implementation, you create a struct and 
use the `json` tag to define its properties like this:

```go
type GreetBody struct {
	Name string `json:"name"`
}
```

Defining this struct with the `router.Request[B, P, Q]` param in the handler looks
like this:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[GreetBody, router.Nil, router.Nil],
) (router.Response[Responses], error) {
    return router.SendOKText(Responses{OK: fmt.Sprintf("hello %s!", request.Body.Name)})
})
```

Notice how you can access the value of `request.Body.Name`, with Cellotape binding 
it for you from the HTTP request. 

This type is validated for compatibility with the schema defined in the spec 
`requestBody` property.

#### Path parameters - <code>router.Request[B, <strong>P</strong>, Q]</code>

The second generic argument (`P`) of `router.Request[B, P, Q]` represents the 
request path parameters.

For example, an API that expects to receive a path parameter `name` in the request URL 
property can be defined in the spec like this:

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

To represent this path parameter in the Go implementation, define a struct.

Each struct field represents a path parameter. Use the `uri` tag to define the 
parameter name, as stated in the spec, like this:

```go
type GreetPathParams struct {
	Name string `uri:"name"`
}
```

Defining this struct with the `router.Request[B, P, Q]` param in the handler looks
like this:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[router.Nil, GreetPathParams, router.Nil],
) (router.Response[Responses], error) {
    return router.SendOKText(Responses{OK: fmt.Sprintf("hello %s!", request.PathParams.Name)})
})
```

Notice how you access the value of `request.PathParams.Name` with Cellotape binding
it for you from the HTTP request.

This type is validated for compatibility with the schema defined in each parameter 
of the spec `parameters` that have `in: path`.

#### Query parameters - <code>router.Request[B, P, <strong>Q</strong>]</code>

The third generic argument (`Q`) of `router.Request[B, P, Q]` represents the r
request query parameters.

For example, an API that expects to receive a query parameter `name` in the request 
URL property can be defined in the spec like this:

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

To represent this query parameter in the Go implementation, define a struct.

Each struct field represents a query parameter. Use the `form` tag to define the 
parameter name, as stated in the spec.

```go
type GreetQueryParams struct {
	Name string `form:"name"`
}
```

Defining this struct with the `router.Request[B, P, Q]` param in the handler looks
like this:

```go
var handler = router.NewHandler(func (
    c router.Context,
    request router.Request[router.Nil, router.Nil, GreetQueryParams],
) (router.Response[Responses], error) {
    return router.SendOKText(Responses{OK: fmt.Sprintf("hello %s!", request.QueryParams.Name)})
})
```

Notice how you access the value of `request.QueryParams.Name`, with Cellotape 
binding it from the HTTP request.

This type is validated for compatibility with the schema defined in each parameter 
of the spec `parameters` that have `in: query`.

### Responses - `router.Response[R]`

The first return type of the handler function is a `router.Response[R]`.

This type of generic argument should be a struct with each field representing a 
response and is tagged with the `status` tag. 

For example, an API that can return 200 and 400 responses in a simple JSON can 
be defined in the spec like this:

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

To represent these responses in the Go implementation, Cellotape defines a struct 
and use the `status` tag to describe each response like this:

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

Defining this struct with the `router.Response[R]` return type in the handler looks like this:

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

Notice how the multiple fields of the responses struct enable you to send the 
response you want based on your implementation logic in a way that covers all the 
responses defined the spec.

Each response type is validated for compatibility with the schema defined in each
 response of the spec `responses`.

## Examples

Examples that show how to use Cellotape.

### Hello World API example

This example uses an [openapi.yaml](./examples/hello_world_example/openapi.yaml)
to initialize the server and map to the relevant handlers. Visit the 
[Hello World Example](./examples/hello_world_example) to see how it works.

### TODO List API example

This example use an [openapi.yaml](./examples/todo_list_app_example/openapi.yaml)
to initialize the server and map to the relevant handlers. Visit the 
[TODO List API Example](./examples/todo_list_app_example) to see how Cellotape 
for a more realistic use.
 
## Roadmap

- [ ] Improve the documentation with more usage details.
- [ ] Runtime validation for request body & parameters based on the OpenAPI spec.
  Considered multiple options, including:
  - https://github.com/Go-playground/validator
  - https://github.com/xeipuuv/Gojsonschema
- [ ] Add support for serialization of OpenAPI parameters `style` and `explode`, 
  and the `allowReserved` property.
- [ ] Add support for OpenAPI Header parameters.
- [ ] Add support for OpenAPI Cookie parameters.
