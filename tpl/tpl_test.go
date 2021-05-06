package tpl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tp, err := Parse([]byte("/tmp/demo-${pim.pam.poum}-${ hop }.${txt}?popo"))
	assert.NoError(t, err)
	resp, err := tp.Execute(map[string]interface{}{
		"pim.pam.poum": "one",
		"txt":          "three",
		"hop":          "two",
	})
	assert.NoError(t, err)
	for _, c := range tp.chunks {
		fmt.Println(string(c))
	}
	assert.Equal(t, "/tmp/demo-one-two.three?popo", string(resp))
}
