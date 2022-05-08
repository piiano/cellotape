package router

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOpenAPIRouter(t *testing.T) {
	spec := NewSpec()
	oa := NewOpenAPIRouter(spec)
	assert.NotNil(t, oa)
}

func TestNewOpenAPIRouterWithOptions(t *testing.T) {
	spec := NewSpec()
	oa := NewOpenAPIRouterWithOptions(spec, Options{})
	assert.NotNil(t, oa)
}

//func TestNewOpenAPIRouterWithOptions(t *testing.T) {
//	spec := NewSpec()
//	oa := NewOpenAPIRouterWithOptions(spec, Options{})
//	assert.NotNil(t, oa)
//}
