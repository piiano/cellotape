package router

import (
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSchema(t *testing.T) {

	schema := jsonschema.Reflect(&Options{})
	bytes, _ := schema.MarshalJSON()
	log.Println(string(bytes))

}
func TestBehaviourZeroValue(t *testing.T) {
	var zeroValueBehaviour Behaviour
	assert.Equal(t, PropagateError, zeroValueBehaviour)
	assert.NotEqual(t, PrintWarning, zeroValueBehaviour)
	assert.NotEqual(t, Ignore, zeroValueBehaviour)
}
