package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/piiano/restcontroller/example"
	"github.com/piiano/restcontroller/restcontroller"
	"os"
	"path/filepath"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	controllers := map[string]restcontroller.Controller{
		"greet": example.GreetController,
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	openapiSpecPath := filepath.Join(cwd, "example/hello-world-openapi.yaml")
	if err := restcontroller.InitRouter(engine, openapiSpecPath, controllers); err != nil {
		panic(err)
	}
	port := 8080
	fmt.Printf("Starting server on port: %d.", port)
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
