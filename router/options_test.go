package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBehaviourZeroValue(t *testing.T) {
	var zeroValueBehaviour Behaviour
	assert.Equal(t, PropagateError, zeroValueBehaviour)
	assert.NotEqual(t, PrintWarning, zeroValueBehaviour)
	assert.NotEqual(t, Off, zeroValueBehaviour)
}
