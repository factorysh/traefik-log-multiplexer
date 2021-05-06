package tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tp, err := Parse([]byte("/tmp/demo-${pim.pam.poum}-${ hop }.${txt}"))
	assert.NoError(t, err)
	resp, err := tp.Execute(map[string]interface{}{
		"pim.pam.poum": "one",
		"txt":          "three",
		"hop":          "two",
	})
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/demo-one-two.three", string(resp))
}
