package tpl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tp, err := Parse("/tmp/demo-${pim.pam.poum}-${ hop }.${txt}")
	assert.NoError(t, err)
	buff := &bytes.Buffer{}
	tp.Execute(buff, map[string]string{
		"pim.pam.poum": "one",
		"txt":          "three",
		"hop":          "two",
	})
	assert.Equal(t, "/tmp/demo-one-two.three", buff.String())
}
