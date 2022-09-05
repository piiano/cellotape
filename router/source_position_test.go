package router

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFileName = "source_position_test.go"

func TestFunctionSourcePosition(t *testing.T) {
	testFuncPos := functionSourcePosition(TestFunctionSourcePosition)
	assert.True(t, testFuncPos.ok)
	assert.Equal(t, testFileName, path.Base(testFuncPos.file))
	assert.Equal(t, 14, testFuncPos.line)

	anonymousFuncPos := functionSourcePosition(func() {})
	assert.True(t, anonymousFuncPos.ok)
	assert.Equal(t, testFileName, path.Base(anonymousFuncPos.file))
	assert.Equal(t, 20, anonymousFuncPos.line)

	funcPosStr := anonymousFuncPos.String()
	assert.True(t, strings.HasSuffix(funcPosStr, fmt.Sprintf(":%d", anonymousFuncPos.line)))
	assert.True(t, strings.HasPrefix(funcPosStr, anonymousFuncPos.file))

	nonFuncPos := functionSourcePosition(42)
	assert.False(t, nonFuncPos.ok)

	nilTypePos := functionSourcePosition(nil)
	assert.False(t, nilTypePos.ok)
}
