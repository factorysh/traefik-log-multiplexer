package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInput(t *testing.T) {
	cfg := map[string]interface{}{
		"path": "/tmp/test.cfg",
	}
	i, err := Inputs["file"](cfg)
	assert.NoError(t, err)
	ii, ok := i.(*FileInput)
	assert.True(t, ok)
	assert.Equal(t, "/tmp/test.cfg", ii.path)
}
